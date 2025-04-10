package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
// TODO: Define actual configuration fields (API keys, endpoint, polling interval, etc.)
type Config struct {
	APIEndpoint     string `mapstructure:"api_endpoint"`
	OpenAIKey       string `mapstructure:"openai_key"`
	AccountID       string `mapstructure:"account_id"`
	PollIntervalSec int    `mapstructure:"poll_interval_sec"`
	NotifyMethod    string `mapstructure:"notify_method"` // e.g., "sms", "log"
	NotifyTarget    string `mapstructure:"notify_target"` // e.g., phone number
	PersistenceFile string `mapstructure:"persistence_file"`
	Timezone        string `mapstructure:"timezone"` // e.g., "America/New_York"
}

// LoadConfig loads configuration from a given file path using Viper.
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("api_endpoint", "https://truthsocial.com/api/v1/accounts/%s/statuses?with_muted=true")
	v.SetDefault("poll_interval_sec", 60)
	v.SetDefault("notify_method", "log")
	v.SetDefault("persistence_file", "processed_tweets.json")
	v.SetDefault("timezone", "America/New_York")

	// Get the file name and directory
	filename := filepath.Base(path)
	dir := filepath.Dir(path)

	// Configure Viper
	v.SetConfigName(strings.TrimSuffix(filename, filepath.Ext(filename)))
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config into struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Check for environment variables and override config values
	if openAIKey := os.Getenv("OPENAI_API_KEY"); openAIKey != "" {
		cfg.OpenAIKey = openAIKey
	}

	// Validate required fields
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validateConfig ensures all required configuration values are set.
func validateConfig(cfg *Config) error {
	if cfg.OpenAIKey == "" || cfg.OpenAIKey == "YOUR_OPENAI_API_KEY_HERE" {
		return fmt.Errorf("OpenAI API key is required. Set it in config.yaml or using the OPENAI_API_KEY environment variable")
	}

	if cfg.AccountID == "" {
		return fmt.Errorf("account_id is required in config file")
	}

	if cfg.PollIntervalSec <= 0 {
		return fmt.Errorf("poll_interval_sec must be greater than 0")
	}

	if cfg.Timezone == "" {
		return fmt.Errorf("timezone is required in config file")
	}

	return nil
}
