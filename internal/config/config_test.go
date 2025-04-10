package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Valid config
	validConfig := `
api_endpoint: "https://test.com/api/v1/accounts/%s/statuses"
openai_key: "test-openai-key"
account_id: "123456789"
poll_interval_sec: 30
notify_method: "log"
persistence_file: "test_tweets.json"
timezone: "UTC"
`

	// Write valid config to file
	err := os.WriteFile(configPath, []byte(validConfig), 0644)
	require.NoError(t, err)

	// Test loading a valid config
	t.Run("ValidConfig", func(t *testing.T) {
		// Ensure env var is not set
		os.Unsetenv("OPENAI_API_KEY")

		cfg, err := LoadConfig(configPath)
		require.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Equal(t, "https://test.com/api/v1/accounts/%s/statuses", cfg.APIEndpoint)
		assert.Equal(t, "test-openai-key", cfg.OpenAIKey)
		assert.Equal(t, "123456789", cfg.AccountID)
		assert.Equal(t, 30, cfg.PollIntervalSec)
		assert.Equal(t, "log", cfg.NotifyMethod)
		assert.Equal(t, "test_tweets.json", cfg.PersistenceFile)
		assert.Equal(t, "UTC", cfg.Timezone)
	})

	// Invalid config (missing required fields)
	invalidConfig := `
api_endpoint: "https://test.com/api/v1/accounts/%s/statuses"
poll_interval_sec: 30
`

	invalidConfigPath := filepath.Join(tempDir, "invalid_config.yaml")
	err = os.WriteFile(invalidConfigPath, []byte(invalidConfig), 0644)
	require.NoError(t, err)

	// Test loading an invalid config
	t.Run("InvalidConfig", func(t *testing.T) {
		// Ensure env var is not set
		os.Unsetenv("OPENAI_API_KEY")

		cfg, err := LoadConfig(invalidConfigPath)
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "OpenAI API key is required")
	})

	// Test loading a non-existent config file
	t.Run("NonExistentConfig", func(t *testing.T) {
		cfg, err := LoadConfig(filepath.Join(tempDir, "nonexistent.yaml"))
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "failed to read config file")
	})

	// Test default values
	minimalConfig := `
openai_key: "test-openai-key"
account_id: "123456789"
timezone: "UTC"
`

	minimalConfigPath := filepath.Join(tempDir, "minimal_config.yaml")
	err = os.WriteFile(minimalConfigPath, []byte(minimalConfig), 0644)
	require.NoError(t, err)

	// Test loading a minimal config with default values
	t.Run("DefaultValues", func(t *testing.T) {
		// Ensure env var is not set
		os.Unsetenv("OPENAI_API_KEY")

		cfg, err := LoadConfig(minimalConfigPath)
		require.NoError(t, err)
		require.NotNil(t, cfg)

		// Check default values
		assert.Contains(t, cfg.APIEndpoint, "truthsocial.com")
		assert.Equal(t, 60, cfg.PollIntervalSec)
		assert.Equal(t, "log", cfg.NotifyMethod)
		assert.Equal(t, "processed_tweets.json", cfg.PersistenceFile)
	})

	// Test OpenAI key from env var
	t.Run("OpenAIKeyFromEnv", func(t *testing.T) {
		// Set environment variable
		os.Setenv("OPENAI_API_KEY", "env-var-openai-key")
		defer os.Unsetenv("OPENAI_API_KEY")

		// Use a config without OpenAI key
		configWithoutAPIKey := `
api_endpoint: "https://test.com/api/v1/accounts/%s/statuses"
account_id: "123456789"
poll_interval_sec: 30
timezone: "UTC"
`
		configWithoutAPIKeyPath := filepath.Join(tempDir, "config_without_api_key.yaml")
		err = os.WriteFile(configWithoutAPIKeyPath, []byte(configWithoutAPIKey), 0644)
		require.NoError(t, err)

		cfg, err := LoadConfig(configWithoutAPIKeyPath)
		require.NoError(t, err)
		require.NotNil(t, cfg)

		// The key should come from env var
		assert.Equal(t, "env-var-openai-key", cfg.OpenAIKey)
	})

	// Test OpenAI key from env var overriding config file
	t.Run("OpenAIKeyEnvOverridesConfig", func(t *testing.T) {
		// Set environment variable
		os.Setenv("OPENAI_API_KEY", "env-var-overrides-config")
		defer os.Unsetenv("OPENAI_API_KEY")

		cfg, err := LoadConfig(configPath) // Using the valid config which has an API key
		require.NoError(t, err)
		require.NotNil(t, cfg)

		// The key should be overridden by env var
		assert.Equal(t, "env-var-overrides-config", cfg.OpenAIKey)
	})
}
