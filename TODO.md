# TODO - Stonk Trumpet

This file tracks the high-level milestones and tasks for the project.

## Phase: MVP (Minimum Viable Product)

### Feature: Core Polling & Analysis

- [x] **Milestone:** Project Setup
  - [x] Initialize Go module
  - [x] Create initial directory structure (`cmd`, `internal`, `pkg`)
  - [x] Define Truth Social API response structs (`pkg/truthsocial`)
  - [x] Create basic `main.go` entry point
  - [x] Create placeholder internal packages (config, fetcher, analyzer, notifier, tracker, scheduler) with basic structs/functions
  - [x] Create basic unit tests for placeholder packages
  - [x] Create `README.md`, `TODO.md`, `.gitignore`, `config.example.yaml`
- [ ] **Milestone:** Configuration Loading
  - [x] Choose and implement a configuration library (e.g., Viper)
  - [x] Load configuration from `config.yaml` in `main.go`
  - [x] Add tests for config loading
- [ ] **Milestone:** Tweet Fetching
  - [x] Implement HTTP client logic in `internal/fetcher`
  - [x] Handle API request construction (URL, headers - potentially mimic browser headers)
  - [x] Parse JSON response into `truthsocial.Status` structs
  - [x] Handle potential API errors and rate limits (basic handling)
  - [ ] Add tests for fetcher (using mock HTTP server)
- [ ] **Milestone:** Sentiment Analysis (OpenAI)
  - [x] Choose and implement OpenAI Go client library (e.g., `go-openai`)
  - [x] Implement logic in `internal/analyzer` to call OpenAI API (e.g., ChatCompletion)
  - [x] Design prompt for sentiment analysis focused on stock market impact (Positive/Negative/Neutral)
  - [x] Implement relevance check (improve basic keyword check, potentially use OpenAI function calling or prompt refinement)
  - [x] Parse OpenAI response to get sentiment and relevance
  - [x] Handle OpenAI API errors
  - [ ] Add tests for analyzer (using mock OpenAI client)
- [ ] **Milestone:** Processed Tweet Tracking
  - [x] Implement tracker logic with file persistence (`internal/tracker`)
  - [x] Implement load/save functionality
  - [x] Implement `MarkProcessed` and `IsProcessed` methods
  - [x] Ensure thread-safety
  - [x] Add comprehensive tests for tracker (persistence, concurrency)
- [ ] **Milestone:** Notification (Logging)
  - [x] Implement basic LogNotifier (`internal/notifier`)
  - [x] Implement Notifier interface and factory function
  - [x] Add tests for LogNotifier
- [ ] **Milestone:** Scheduling & Time Restriction
  - [x] Implement scheduler logic (`internal/scheduler`)
  - [x] Implement time restriction check based on timezone
  - [x] Implement main polling loop coordinating fetcher, tracker, analyzer, notifier
  - [x] Implement graceful shutdown
  - [ ] Add tests for time restriction logic
  - [ ] Add integration tests for scheduler Start/Stop (TODO: Requires more setup)
- [ ] **Milestone:** Component Wiring
  - [x] Initialize all components (Fetcher, Analyzer, Notifier, Tracker, Scheduler) in `main.go` using loaded config
  - [x] Start the Scheduler
  - [x] Handle graceful shutdown signal to stop Scheduler and save tracker state

## Phase: v1

### Feature: Enhanced Notifications

- [ ] **Milestone:** SMS Notifications
  - [ ] Add Twilio (or alternative) client library
  - [ ] Implement `SMSNotifier` conforming to the `Notifier` interface
  - [ ] Add configuration options for SMS (API keys, phone numbers)
  - [ ] Update notifier factory function
  - [ ] Add tests for `SMSNotifier` (using mock Twilio client)

### Feature: Robustness & Error Handling

- [ ] **Milestone:** Improved Error Handling
  - [ ] Refine error handling in fetcher (retries?)
  - [ ] Refine error handling in analyzer (retries?)
  - [ ] Refine error handling in notifier
- [ ] **Milestone:** Observability
  - [ ] Add structured logging (e.g., `log/slog`)
  - [ ] Consider adding basic metrics (e.g., tweets processed, errors, notifications sent)

### Feature: Deployment

- [ ] **Milestone:** Containerization
  - [ ] Create a `Dockerfile` for building the application image
  - [ ] Add instructions for building/running the Docker image to `README.md`
