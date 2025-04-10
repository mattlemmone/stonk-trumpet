package scheduler

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"stonk-trumpet/internal/analyzer"
	"stonk-trumpet/internal/config"
	"stonk-trumpet/internal/fetcher"
	"stonk-trumpet/internal/notifier"
	"stonk-trumpet/internal/tracker"
	"stonk-trumpet/pkg/truthsocial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---- Test Mocks ----

// Basic mocks now directly implement the required interfaces
type mockFetcher struct {
	statuses []truthsocial.Status
	err      error
}

func (m *mockFetcher) FetchStatuses() ([]truthsocial.Status, error) { return m.statuses, m.err }

type mockAnalyzer struct {
	results map[string]analyzer.AnalysisResult // Use ID as key for simplicity
	err     error
}

func (m *mockAnalyzer) AnalyzeSentiment(status *truthsocial.Status) (analyzer.AnalysisResult, error) {
	if res, ok := m.results[status.ID]; ok {
		return res, m.err
	}
	return analyzer.AnalysisResult{Sentiment: analyzer.Neutral, IsRelevant: false}, m.err
}

type mockNotifier struct {
	notifyCalled map[string]bool
	err          error
}

func (m *mockNotifier) Notify(status *truthsocial.Status, result analyzer.AnalysisResult) error {
	if m.notifyCalled == nil {
		m.notifyCalled = make(map[string]bool)
	}
	m.notifyCalled[status.ID] = true
	return m.err
}

// Embedding mocks no longer needed
// type mockFetcherEmbed struct { ... }
// func (m *mockFetcherEmbed) FetchStatuses() ...
// type mockAnalyzerEmbed struct { ... }
// func (m *mockAnalyzerEmbed) AnalyzeSentiment() ...

// ---- Test Helper ----

func setupTestComponents(t *testing.T) (*config.Config, *tracker.Tracker) {
	t.Helper()
	cfg := &config.Config{
		Timezone:        "UTC",
		PollIntervalSec: 60,
		PersistenceFile: filepath.Join(t.TempDir(), "test_tracker.txt"), // Use .txt
	}
	// Use the simplified NewTracker
	tracker, err := tracker.NewTracker(cfg.PersistenceFile)
	if err != nil {
		t.Fatalf("Failed to create test tracker: %v", err)
	}
	return cfg, tracker
}

// Helper to create scheduler with mocks for runPollingCycle tests
func newTestSchedulerForPolling(t *testing.T, cfg *config.Config, tracker *tracker.Tracker, fetcher fetcher.Fetcher, analyzer analyzer.Analyzer, notifier notifier.Notifier) *Scheduler {
	loc, err := time.LoadLocation(cfg.Timezone)
	require.NoError(t, err) // Check error loading location

	s, err := NewScheduler(cfg, fetcher, analyzer, notifier, tracker)
	require.NoError(t, err)
	s.loc = loc // Ensure location is set (might be redundant if NewScheduler sets it)
	return s
}

// Helper to mock time within allowed range for tests
func mockTimeInAllowedRange(t *testing.T) func() {
	originalNow := timeNow
	// Set time to a fixed point within 7am-midnight ET (e.g., Noon ET)
	mockedTime := time.Date(2024, 1, 1, 17, 0, 0, 0, time.UTC) // Noon ET is 17:00 UTC
	timeNow = func() time.Time { return mockedTime }
	// Return a function to restore the original time
	return func() {
		timeNow = originalNow
	}
}

// ---- Tests ----

func TestScheduler_RunPollingCycle_Basic(t *testing.T) {
	restoreTime := mockTimeInAllowedRange(t) // Mock time
	defer restoreTime()                      // Restore time after test

	cfg, trk := setupTestComponents(t)
	now := timeNow().UTC()
	ts1 := now.Add(-3 * time.Minute)
	ts2 := now.Add(-2 * time.Minute)
	ts3 := now.Add(-1 * time.Minute)
	ts4 := now // Newest

	// Set initial last seen time
	initialLastSeen := now.Add(-5 * time.Minute)
	trk.SetLastSeenTime(initialLastSeen)

	fetcherMock := &mockFetcher{
		statuses: []truthsocial.Status{
			// API might return newest first
			{ID: "4", Content: "<p>Another positive finance tweet.</p>", CreatedAt: ts4, URL: "url4"},
			{ID: "3", Content: "<p>Irrelevant tweet.</p>", CreatedAt: ts3, URL: "url3"},
			{ID: "2", Content: "<p>Negative market outlook.</p>", CreatedAt: ts2, URL: "url2"},
			{ID: "1", Content: "<p>Positive stock news!</p>", CreatedAt: ts1, URL: "url1"},
			{ID: "0", Content: "<p>Old tweet</p>", CreatedAt: initialLastSeen, URL: "url0"}, // Should be skipped
		},
	}
	analyzerMock := &mockAnalyzer{
		results: map[string]analyzer.AnalysisResult{
			"1": {Sentiment: analyzer.Positive, IsRelevant: true},
			"2": {Sentiment: analyzer.Negative, IsRelevant: true},
			"3": {Sentiment: analyzer.Neutral, IsRelevant: false},
			"4": {Sentiment: analyzer.Positive, IsRelevant: true},
		},
	}
	notifierMock := &mockNotifier{}

	scheduler := newTestSchedulerForPolling(t, cfg, trk, fetcherMock, analyzerMock, notifierMock)
	scheduler.runPollingCycle()

	// Check tracker's last seen time - should be the max time from the batch (ts4)
	assert.Equal(t, ts4.UnixNano(), trk.GetLastSeenTime().UnixNano(), "Tracker time not updated to max batch time")

	// Check notifier - only relevant positive tweets newer than initialLastSeen (1 and 4)
	assert.True(t, notifierMock.notifyCalled["1"], "Expected notification for status 1")
	assert.False(t, notifierMock.notifyCalled["2"], "Expected NO notification for status 2 (negative)")
	assert.False(t, notifierMock.notifyCalled["3"], "Expected NO notification for status 3 (irrelevant)")
	assert.True(t, notifierMock.notifyCalled["4"], "Expected notification for status 4")
	assert.False(t, notifierMock.notifyCalled["0"], "Expected NO notification for status 0 (old)")
	assert.Len(t, notifierMock.notifyCalled, 2, "Expected 2 notifications")

	// Check saved state using package function tracker.NewTracker
	trackerSaved, err := tracker.NewTracker(cfg.PersistenceFile)
	require.NoError(t, err)
	assert.Equal(t, ts4.UnixNano(), trackerSaved.GetLastSeenTime().UnixNano(), "Saved tracker time mismatch")
}

func TestScheduler_RunPollingCycle_NoNewTweets(t *testing.T) {
	restoreTime := mockTimeInAllowedRange(t)
	defer restoreTime()

	cfg, trk := setupTestComponents(t)
	now := timeNow().UTC()
	initialLastSeen := now.Add(-5 * time.Minute)
	trk.SetLastSeenTime(initialLastSeen)
	// Explicitly save the initial state *before* running the cycle
	err := trk.SaveState()
	require.NoError(t, err, "Failed to save initial tracker state")

	fetcherMock := &mockFetcher{
		statuses: []truthsocial.Status{
			{ID: "0", Content: "<p>Old tweet</p>", CreatedAt: initialLastSeen, URL: "url0"},
			{ID: "-1", Content: "<p>Even older</p>", CreatedAt: initialLastSeen.Add(-1 * time.Hour), URL: "url-1"},
		},
	}
	analyzerMock := &mockAnalyzer{}
	notifierMock := &mockNotifier{}

	scheduler := newTestSchedulerForPolling(t, cfg, trk, fetcherMock, analyzerMock, notifierMock)
	scheduler.runPollingCycle()

	// Check tracker's last seen time - should NOT have changed
	assert.Equal(t, initialLastSeen.UnixNano(), trk.GetLastSeenTime().UnixNano(), "Tracker time updated unexpectedly")

	// Check notifier
	assert.Len(t, notifierMock.notifyCalled, 0, "Expected 0 notifications")

	// Check that state file still contains the initial time
	trackerSaved, err := tracker.NewTracker(cfg.PersistenceFile)
	require.NoError(t, err)
	assert.Equal(t, initialLastSeen.UnixNano(), trackerSaved.GetLastSeenTime().UnixNano(), "Saved tracker time mismatch (should be initial)")
}

func TestScheduler_RunPollingCycle_AnalysisError(t *testing.T) {
	restoreTime := mockTimeInAllowedRange(t) // Mock time
	defer restoreTime()                      // Restore time after test

	cfg, trk := setupTestComponents(t)
	now := timeNow().UTC()
	ts1 := now.Add(-1 * time.Minute)
	initialLastSeen := now.Add(-5 * time.Minute)
	trk.SetLastSeenTime(initialLastSeen)

	fetcherMock := &mockFetcher{
		statuses: []truthsocial.Status{{ID: "1", Content: "<p>Analyze me</p>", CreatedAt: ts1, URL: "url1"}},
	}
	analyzerMock := &mockAnalyzer{err: fmt.Errorf("analysis failed")}
	notifierMock := &mockNotifier{}

	scheduler := newTestSchedulerForPolling(t, cfg, trk, fetcherMock, analyzerMock, notifierMock)
	scheduler.runPollingCycle()

	// Check tracker's last seen time - should be updated even on analysis error
	assert.Equal(t, ts1.UnixNano(), trk.GetLastSeenTime().UnixNano(), "Tracker time not updated on analysis error")

	// Check notifier
	assert.Len(t, notifierMock.notifyCalled, 0, "Expected 0 notifications on analysis error")

	// Check saved state using package function tracker.NewTracker
	trackerSaved, err := tracker.NewTracker(cfg.PersistenceFile)
	require.NoError(t, err)
	assert.Equal(t, ts1.UnixNano(), trackerSaved.GetLastSeenTime().UnixNano(), "Saved tracker time mismatch on analysis error")
}

func TestScheduler_IsWithinAllowedTime(t *testing.T) {
	cfg := &config.Config{Timezone: "America/New_York"}
	loc, _ := time.LoadLocation(cfg.Timezone)

	tests := []struct {
		name     string
		testTime time.Time
		expected bool
	}{
		{"Midday ET (Allowed)", time.Date(2024, 1, 1, 17, 0, 0, 0, time.UTC), true},
		{"Morning ET (Allowed)", time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), true},
		{"Late Night ET (Allowed)", time.Date(2024, 1, 2, 4, 59, 59, 0, time.UTC), true},
		{"Early Morning ET (Not Allowed)", time.Date(2024, 1, 1, 11, 59, 59, 0, time.UTC), false},
		{"Midnight ET (Not Allowed)", time.Date(2024, 1, 2, 5, 0, 0, 0, time.UTC), false},
		{"Exactly Midnight (Not Allowed)", time.Date(2024, 1, 1, 5, 0, 0, 0, time.UTC), false},
		{"Exactly 7 AM (Allowed)", time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), true},
	}

	originalNow := timeNow
	defer func() { timeNow = originalNow }()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			timeNow = func() time.Time { return tt.testTime }
			schedulerForTest := &Scheduler{loc: loc}
			actual := schedulerForTest.isWithinAllowedTime()
			assert.Equal(t, tt.expected, actual, fmt.Sprintf("isWithinAllowedTime() for %s (%s ET)", tt.testTime.Format(time.RFC3339), timeNow().In(loc).Format("15:04:05")))
		})
	}
}

// TODO: Add tests for Start/Stop behavior.
// TODO: Refactor scheduler fetcher/analyzer dependencies to use interfaces.
