package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"stonk-trumpet/pkg/truthsocial"
	"time"
)

// Fetcher defines the interface for fetching statuses.
type Fetcher interface {
	FetchStatuses() ([]truthsocial.Status, error)
}

// httpFetcher implements the Fetcher interface using HTTP.
type httpFetcher struct {
	client      *http.Client
	apiEndpoint string
	accountID   string
	mockMode    bool // Flag to enable mock mode
}

// NewHTTPFetcher creates a new httpFetcher instance.
func NewHTTPFetcher(apiEndpoint, accountID string) Fetcher {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &httpFetcher{
		client:      client,
		apiEndpoint: apiEndpoint,
		accountID:   accountID,
		mockMode:    false,
	}
}

// NewMockFetcher creates a fetcher in mock mode for testing.
func NewMockFetcher() Fetcher {
	return &httpFetcher{
		mockMode: true,
	}
}

// FetchStatuses retrieves the latest statuses from the Truth Social API.
func (f *httpFetcher) FetchStatuses() ([]truthsocial.Status, error) {
	// If in mock mode, return mock data
	if f.mockMode {
		return f.generateMockStatuses(), nil
	}

	url := fmt.Sprintf(f.apiEndpoint, f.accountID)

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add headers to mimic a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	// Execute request
	log.Printf("Fetching statuses from: %s", url)
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var statuses []truthsocial.Status
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&statuses); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	log.Printf("Successfully fetched %d statuses", len(statuses))
	return statuses, nil
}

// generateMockStatuses creates fake statuses for testing
func (f *httpFetcher) generateMockStatuses() []truthsocial.Status {
	log.Println("Using mock mode: Generating fake statuses")

	now := time.Now()

	// Create some sample statuses with varying content and timestamps
	statuses := []truthsocial.Status{
		{
			ID:        "mock1",
			CreatedAt: now.Add(-10 * time.Minute),
			Content:   "The stock market is looking very strong today! Economy is booming and I expect it to continue.",
			URL:       "https://truthsocial.com/mock1",
			Account: truthsocial.Account{
				ID:          "107780257626128497",
				Username:    "realDonaldTrump",
				DisplayName: "Donald J. Trump",
			},
		},
		{
			ID:        "mock2",
			CreatedAt: now.Add(-30 * time.Minute),
			Content:   "Just had a wonderful meeting with business leaders. Great things coming for our economy!",
			URL:       "https://truthsocial.com/mock2",
			Account: truthsocial.Account{
				ID:          "107780257626128497",
				Username:    "realDonaldTrump",
				DisplayName: "Donald J. Trump",
			},
		},
		{
			ID:        "mock3",
			CreatedAt: now.Add(-1 * time.Hour),
			Content:   "Beautiful day in New York. Going to play golf with some friends later.",
			URL:       "https://truthsocial.com/mock3",
			Account: truthsocial.Account{
				ID:          "107780257626128497",
				Username:    "realDonaldTrump",
				DisplayName: "Donald J. Trump",
			},
		},
	}

	log.Printf("Generated %d mock statuses", len(statuses))
	return statuses
}
