package scheduler

import (
	"fmt"
	"log"
	"time"

	"stonk-trumpet/internal/analyzer"
	"stonk-trumpet/internal/config"
	"stonk-trumpet/internal/fetcher"
	"stonk-trumpet/internal/notifier"
	"stonk-trumpet/internal/tracker"
	// Added for truthsocial.Status
)

// timeNow is a variable holding the time source, allowing mocking in tests.
var timeNow = time.Now

// Scheduler orchestrates the fetching, analyzing, and notification process.
type Scheduler struct {
	config   *config.Config
	fetcher  fetcher.Fetcher   // Use interface type
	analyzer analyzer.Analyzer // Use interface type
	notifier notifier.Notifier // Already an interface
	tracker  *tracker.Tracker
	loc      *time.Location
	ticker   *time.Ticker
	done     chan bool
}

// NewScheduler creates a new Scheduler instance.
// Now accepts interface types for dependencies.
func NewScheduler(cfg *config.Config, f fetcher.Fetcher, a analyzer.Analyzer, n notifier.Notifier, t *tracker.Tracker) (*Scheduler, error) {
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone %q: %w", cfg.Timezone, err)
	}

	pollInterval := time.Duration(cfg.PollIntervalSec) * time.Second
	if pollInterval <= 0 {
		pollInterval = 60 * time.Second
		log.Printf("Warning: Invalid poll_interval_sec (%d), using default 60s", cfg.PollIntervalSec)
	}

	// Basic validation: Check if dependencies are nil
	if f == nil {
		return nil, fmt.Errorf("fetcher cannot be nil")
	}
	if a == nil {
		return nil, fmt.Errorf("analyzer cannot be nil")
	}
	if n == nil {
		return nil, fmt.Errorf("notifier cannot be nil")
	}
	if t == nil {
		return nil, fmt.Errorf("tracker cannot be nil")
	}

	return &Scheduler{
		config:   cfg,
		fetcher:  f, // Assign interfaces directly
		analyzer: a,
		notifier: n,
		tracker:  t,
		loc:      loc,
		ticker:   time.NewTicker(pollInterval),
		done:     make(chan bool),
	}, nil
}

// Start begins the scheduled polling loop.
func (s *Scheduler) Start() {
	log.Println("Scheduler started. Polling every", s.ticker.C)
	go func() {
		s.runPollingCycle()
		for {
			select {
			case <-s.ticker.C:
				s.runPollingCycle()
			case <-s.done:
				log.Println("Scheduler stopping polling loop.")
				return
			}
		}
	}()
}

// Stop halts the scheduler loop.
func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	s.ticker.Stop()
	close(s.done)
	// No need to explicitly save tracker state here anymore,
	// as it's saved after each successful batch processing.
	log.Println("Scheduler stopped.")
}

// runPollingCycle performs one cycle of fetching, analyzing, and notifying.
func (s *Scheduler) runPollingCycle() {
	if !s.isWithinAllowedTime() {
		log.Println("Current time is outside allowed polling hours (7am-midnight ET). Skipping cycle.")
		return
	}

	log.Println("Running polling cycle...")

	// Get the timestamp of the last successfully processed status
	lastSeenTime := s.tracker.GetLastSeenTime()
	log.Printf("Fetching statuses created after: %s", lastSeenTime.Format(time.RFC3339))

	statuses, err := s.fetcher.FetchStatuses()
	if err != nil {
		log.Printf("Error fetching statuses: %v", err)
		return
	}

	log.Printf("Fetched %d statuses.", len(statuses))
	processedCount := 0
	notifiedCount := 0
	var currentBatchMaxTime time.Time = lastSeenTime // Initialize with the last known time
	newStatusesProcessed := false

	// Process statuses, likely newest first from API
	for i := range statuses {
		status := &statuses[i]

		// Skip statuses that are not newer than the last seen one
		if !status.CreatedAt.After(lastSeenTime) {
			continue
		}

		newStatusesProcessed = true
		log.Printf("Processing new status: %s (%s)", status.ID, status.CreatedAt.Format(time.RFC3339))

		// Update the max time seen in *this* batch
		if status.CreatedAt.After(currentBatchMaxTime) {
			currentBatchMaxTime = status.CreatedAt
		}

		result, err := s.analyzer.AnalyzeSentiment(status)
		if err != nil {
			log.Printf("Error analyzing status %s: %v", status.ID, err)
			// Don't update tracker time if analysis fails for this one?
			// Or update tracker time anyway to avoid reprocessing on next cycle?
			// Let's update anyway to prevent reprocessing loop on persistent analysis error.
			processedCount++ // Count it as processed even if analysis failed
			continue
		}

		processedCount++

		if result.IsRelevant && result.Sentiment == analyzer.Positive {
			log.Printf("Relevant positive status found: %s. Notifying...", status.ID)
			err = s.notifier.Notify(status, result)
			if err != nil {
				log.Printf("Error sending notification for status %s: %v", status.ID, err)
			} else {
				notifiedCount++
			}
		}
	}

	// If we processed any new statuses in this batch, update and save the tracker
	if newStatusesProcessed {
		log.Printf("Updating last seen time to: %s", currentBatchMaxTime.Format(time.RFC3339))
		s.tracker.SetLastSeenTime(currentBatchMaxTime)
		if err := s.tracker.SaveState(); err != nil {
			log.Printf("Error saving tracker state: %v", err)
		}
	}

	log.Printf("Polling cycle finished. Processed: %d, Notified: %d", processedCount, notifiedCount)
}

// isWithinAllowedTime checks if the current time is within allowed hours.
func (s *Scheduler) isWithinAllowedTime() bool {
	now := timeNow().In(s.loc)
	hour := now.Hour()
	return hour >= 7 && hour < 24 // Between 7am and midnight
}
