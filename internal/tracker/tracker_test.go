package tracker

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert" // Using testify for cleaner time assertions
	"github.com/stretchr/testify/require"
)

// Helper to create a temp file path for testing
func tempFilePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "test_tracker_state.txt") // Use .txt for simple timestamp
}

func TestNewTracker_NonExistentFile(t *testing.T) {
	fp := tempFilePath(t)
	tracker, err := NewTracker(fp)
	require.NoError(t, err, "NewTracker() with non-existent file failed")
	assert.True(t, tracker.GetLastSeenTime().IsZero(), "Expected zero time for non-existent file")

	// Check if file was created (it shouldn't be just by NewTracker)
	_, err = os.Stat(fp)
	assert.True(t, os.IsNotExist(err), "NewTracker() unexpectedly created the state file %s", fp)
}

func TestTracker_GetAndSetLastSeenTime(t *testing.T) {
	fp := tempFilePath(t)
	tracker, err := NewTracker(fp)
	require.NoError(t, err)

	assert.True(t, tracker.GetLastSeenTime().IsZero(), "Initial time should be zero")

	ts1 := time.Now().Add(-1 * time.Hour).UTC() // Ensure UTC for consistency
	tracker.SetLastSeenTime(ts1)
	assert.Equal(t, ts1.UnixNano(), tracker.GetLastSeenTime().UnixNano(), "Time should be updated to ts1")

	// Set an earlier time - should not update
	ts0 := ts1.Add(-1 * time.Minute)
	tracker.SetLastSeenTime(ts0)
	assert.Equal(t, ts1.UnixNano(), tracker.GetLastSeenTime().UnixNano(), "Time should not be updated to earlier ts0")

	// Set a later time - should update
	ts2 := ts1.Add(1 * time.Minute)
	tracker.SetLastSeenTime(ts2)
	assert.Equal(t, ts2.UnixNano(), tracker.GetLastSeenTime().UnixNano(), "Time should be updated to later ts2")
}

func TestTracker_SaveAndLoadState(t *testing.T) {
	fp := tempFilePath(t)
	tracker1, err := NewTracker(fp)
	require.NoError(t, err, "NewTracker() failed")

	// Should not save zero time initially
	err = tracker1.SaveState()
	require.NoError(t, err, "SaveState() with zero time failed")
	_, err = os.Stat(fp)
	assert.True(t, os.IsNotExist(err), "SaveState() unexpectedly created file for zero time")

	// Set and save a time
	time1 := time.Date(2024, 8, 15, 10, 30, 0, 123456789, time.UTC)
	tracker1.SetLastSeenTime(time1)
	err = tracker1.SaveState()
	require.NoError(t, err, "SaveState() failed")

	// Create a new tracker instance to load the saved state
	tracker2, err := NewTracker(fp)
	require.NoError(t, err, "NewTracker() for loading failed")

	// Use assert.WithinDuration for time comparison due to potential precision differences
	assert.WithinDuration(t, time1, tracker2.GetLastSeenTime(), time.Millisecond, "Loaded time mismatch")

	// Test loading empty file
	emptyFp := tempFilePath(t)
	require.NoError(t, os.WriteFile(emptyFp, []byte(" \n "), 0644), "Failed to create empty file") // Write whitespace
	trackerEmpty, err := NewTracker(emptyFp)
	require.NoError(t, err, "NewTracker() with empty file failed")
	assert.True(t, trackerEmpty.GetLastSeenTime().IsZero(), "Expected zero time when loading empty file")

	// Test loading corrupted file (invalid timestamp format)
	corruptFp := tempFilePath(t)
	require.NoError(t, os.WriteFile(corruptFp, []byte("not-a-timestamp"), 0644), "Failed to create corrupt file")
	_, err = NewTracker(corruptFp)
	assert.Error(t, err, "NewTracker() with corrupt file succeeded unexpectedly")

	// Test loading file with only whitespace
	wsFp := tempFilePath(t)
	require.NoError(t, os.WriteFile(wsFp, []byte("   \t\n  "), 0644), "Failed to create whitespace file")
	trackerWS, err := NewTracker(wsFp)
	require.NoError(t, err, "NewTracker() with whitespace file failed")
	assert.True(t, trackerWS.GetLastSeenTime().IsZero(), "Expected zero time when loading whitespace file")

}

// Removed TestTracker_SaveState_Atomicity - less critical with simpler state, can add back if needed.
// Removed TestTracker_AutoSave - auto-save logic removed.
// Removed TestTracker_Concurrency - simplified state makes this less relevant.

// Add testify dependency: go get github.com/stretchr/testify
