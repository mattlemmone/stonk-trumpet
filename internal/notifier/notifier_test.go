package notifier

import (
	"stonk-trumpet/internal/analyzer"
	"stonk-trumpet/pkg/truthsocial"
	"testing"
	"time"
)

func TestLogNotifier(t *testing.T) {
	notifier := NewLogNotifier()
	status := &truthsocial.Status{
		ID:        "123",
		CreatedAt: time.Now(),
		URL:       "http://example.com/status/123",
		Content:   "<p>Test status</p>",
	}
	result := analyzer.AnalysisResult{
		Sentiment:  analyzer.Positive,
		IsRelevant: true,
	}

	// Simply call Notify, we can't easily check log output without redirecting
	// or using a mock logger, which might be overkill for this simple case.
	// We're mainly testing that it doesn't panic.
	err := notifier.Notify(status, result)
	if err != nil {
		t.Errorf("LogNotifier.Notify() returned an unexpected error: %v", err)
	}
}

func TestNewNotifier(t *testing.T) {
	_, err := NewNotifier("log", "")
	if err != nil {
		t.Errorf("NewNotifier(\"log\", \"\") returned error: %v", err)
	}

	_, err = NewNotifier("unsupported", "")
	if err == nil {
		t.Errorf("NewNotifier(\"unsupported\", \"\") expected an error but got nil")
	}
}
