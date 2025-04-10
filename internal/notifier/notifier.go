package notifier

import (
	"fmt"
	"log"
	"stonk-trumpet/internal/analyzer" // For Sentiment type
	"stonk-trumpet/pkg/truthsocial"
)

// Notifier is responsible for sending notifications.
type Notifier interface {
	Notify(status *truthsocial.Status, result analyzer.AnalysisResult) error
}

// LogNotifier simply logs the notification.
type LogNotifier struct{}

func NewLogNotifier() *LogNotifier {
	return &LogNotifier{}
}

func (n *LogNotifier) Notify(status *truthsocial.Status, result analyzer.AnalysisResult) error {
	log.Printf("NOTIFICATION: Status %s (%s) - Sentiment: %s, Relevant: %t, URL: %s\n",
		status.ID,
		status.CreatedAt,
		result.Sentiment,
		result.IsRelevant,
		status.URL,
	)
	// Potentially log more details from the status content
	// log.Printf("Content: %s\n", status.Content)
	return nil
}

// TODO: Implement other notifier types (e.g., SMSNotifier using Twilio)

// NewNotifier creates a notifier based on the configuration.
func NewNotifier(method, target string) (Notifier, error) {
	switch method {
	case "log":
		return NewLogNotifier(), nil
	// case "sms":
	// 	 // TODO: Initialize SMS client (e.g., Twilio)
	// 	 return NewSMSNotifier(target /*, twilioClient */), nil
	default:
		return nil, fmt.Errorf("unsupported notification method: %s", method)
	}
}
