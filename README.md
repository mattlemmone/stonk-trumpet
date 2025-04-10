# Stonk Trumpet 🎺

A Go application that monitors Truth Social statuses (tweets) from a specific account, analyzes their sentiment regarding potential stock market impact using OpenAI, and sends notifications for significant findings.

## Features (Planned)

- Polls a Truth Social account's status endpoint for new tweets.
- Operates only during specified hours (e.g., 7am - midnight Eastern Time).
- Uses OpenAI's API to analyze the sentiment of fetched tweets.
- Classifies tweets as positive, negative, or neutral regarding stock market relevance.
- Sends notifications (e.g., Log, SMS - planned) for relevant, positive tweets.
- Tracks processed tweets using their unique IDs to prevent duplicate processing.
- Persists the list of processed tweets across application restarts.
- Runs as a continuous server application.

## Getting Started

### Prerequisites

- Go (version 1.21 or later recommended)
- Access to OpenAI API (requires an API key)
- (Optional) Access to a notification service like Twilio for SMS.

### Configuration

1.  Copy `config.example.yaml` to `config.yaml`.
2.  Edit `config.yaml` and fill in the required values:
    - `account_id`: The Truth Social account ID to monitor (e.g., `107780257626128497` for @realDonaldTrump).
    - (Optional) Customize `poll_interval_sec`, `notify_method`, `notify_target`, `persistence_file`, `timezone`.
3.  Set your OpenAI API key as an environment variable:
    ```bash
    export OPENAI_API_KEY=your_openai_api_key
    ```
    Alternatively, you can include it in the config file as `openai_key`, but using an environment variable is more secure.

### Building

```bash
go build ./cmd/server
```

This will create an executable named `server` (or `server.exe` on Windows) in the current directory.

### Running

```bash
./server
# Or run directly without building:
# go run ./cmd/server/main.go
```

The application will start, load the configuration, initialize components, and begin polling according to the schedule and time restrictions defined in the configuration.

#### Testing with Mock Mode

For testing without making real API calls to Truth Social or OpenAI, you can enable mock mode:

```bash
MOCK_MODE=true go run ./cmd/server/main.go
```

In mock mode:

- The fetcher generates fake Truth Social statuses instead of calling the real API
- The analyzer uses a simple keyword-based approach instead of calling OpenAI
- This is useful for testing and development when you don't want to use your API quota

## Development

### Running Tests

```bash
go test ./...
```

### Dependencies

Dependencies are managed using Go Modules. They will be downloaded automatically when building or testing.

- (Planned) `github.com/spf13/viper` for configuration management.
- (Planned) `github.com/sashabaranov/go-openai` for OpenAI interaction.
- (Planned) `github.com/twilio/twilio-go` for Twilio SMS notifications.
