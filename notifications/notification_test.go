package notifications

import (
	"strings"
	"testing"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// TestAggregatedNotificationEmpty verifies aggregatedNotification returns empty for no reminders.
func TestAggregatedNotificationEmpty(t *testing.T) {
	title, content := aggregatedNotification([]scheduling.Reminder{})

	if title != "" || content != "" {
		t.Errorf("aggregatedNotification([]) = (%q, %q), want (\"\", \"\")", title, content)
	}
}

// TestAggregatedNotificationSingle verifies aggregatedNotification combines title and content.
func TestAggregatedNotificationSingle(t *testing.T) {
	reminders := []scheduling.Reminder{
		{NotificationCategory: scheduling.ReminderTypeEye, Message: "Look away"},
	}
	title, content := aggregatedNotification(reminders)

	if title == "" {
		t.Errorf("aggregatedNotification title is empty")
	}

	if !strings.Contains(content, "Look away") {
		t.Errorf("aggregatedNotification content = %q, should contain \"Look away\"", content)
	}
}

// TestAggregatedNotificationMultiple verifies aggregatedNotification handles mixed reminders.
func TestAggregatedNotificationMultiple(t *testing.T) {
	reminders := []scheduling.Reminder{
		{NotificationCategory: scheduling.ReminderTypeEye, Message: "Eye care"},
		{NotificationCategory: scheduling.ReminderTypeStandup, Message: "Stand up"},
	}
	title, content := aggregatedNotification(reminders)

	// Title should be from CombinedHeaders
	found := false
	for _, header := range i18n.Active.NotificationCombinedHeaders {
		if title == header {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("aggregatedNotification title %q not in NotificationCombinedHeaders", title)
	}

	// Content should have both messages
	if !strings.Contains(content, "Eye care") || !strings.Contains(content, "Stand up") {
		t.Errorf("aggregatedNotification content missing expected messages: %q", content)
	}
}
