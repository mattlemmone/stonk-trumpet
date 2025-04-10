package tracker

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const timeFormat = time.RFC3339Nano // Use a standard format with precision

// Tracker keeps track of the timestamp of the last processed tweet.
type Tracker struct {
	filepath          string
	lastSeenCreatedAt time.Time
	mu                sync.RWMutex
}

// NewTracker creates a new Tracker instance, loading the last seen timestamp.
func NewTracker(filepath string) (*Tracker, error) {
	t := &Tracker{
		filepath: filepath,
		// lastSeenCreatedAt defaults to zero time
	}
	err := t.loadState()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load tracker state from %s: %w", filepath, err)
	}
	// If the file doesn't exist, we start with the zero time.
	return t, nil
}

// GetLastSeenTime returns the timestamp of the last processed tweet.
func (t *Tracker) GetLastSeenTime() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.lastSeenCreatedAt
}

// SetLastSeenTime updates the timestamp of the last processed tweet.
func (t *Tracker) SetLastSeenTime(ts time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	// Only update if the new timestamp is actually later
	if ts.After(t.lastSeenCreatedAt) {
		t.lastSeenCreatedAt = ts
	}
}

// loadState loads the last seen timestamp from the file.
func (t *Tracker) loadState() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	data, err := os.ReadFile(t.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			t.lastSeenCreatedAt = time.Time{} // Not found, start fresh
			return nil                        // Don't treat not found as an error for load
		}
		return err // Other read error
	}

	timestampStr := strings.TrimSpace(string(data))

	// If file is empty or contains only whitespace, start fresh
	if timestampStr == "" {
		t.lastSeenCreatedAt = time.Time{}
		return nil
	}

	parsedTime, err := time.Parse(timeFormat, timestampStr)
	if err != nil {
		// Attempt to parse with plain RFC3339 as fallback
		fallbackTime, fallbackErr := time.Parse(time.RFC3339, timestampStr)
		if fallbackErr != nil {
			return fmt.Errorf("failed to parse timestamp string %q: %w (fallback attempt: %v)", timestampStr, err, fallbackErr)
		}
		parsedTime = fallbackTime
		err = nil // Not needed, just use parsedTime
	}

	t.lastSeenCreatedAt = parsedTime
	return nil
}

// SaveState saves the current last seen timestamp to the file.
func (t *Tracker) SaveState() error {
	t.mu.RLock()
	lastSeen := t.lastSeenCreatedAt
	t.mu.RUnlock()
	if lastSeen.IsZero() {
		return nil
	}
	timestampStr := lastSeen.Format(timeFormat)
	data := []byte(timestampStr)
	tempFilepath := t.filepath + ".tmp"
	err := os.WriteFile(tempFilepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write temporary tracker state file: %w", err)
	}
	err = os.Rename(tempFilepath, t.filepath)
	if err != nil {
		_ = os.Remove(tempFilepath)
		return fmt.Errorf("failed to rename temporary tracker state file to %s: %w", t.filepath, err)
	}
	return nil
}
