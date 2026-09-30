package notifications

import (
	"strings"
	"testing"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func TestAggregatedNotificationEmpty(t *testing.T) {
	title, content := aggregatedNotification([]scheduling.Reminder{})

	if title != "" || content != "" {
		t.Errorf("aggregatedNotification([]) = (%q, %q), want (\"\", \"\")", title, content)
	}
}

func TestAggregatedNotificationSingle(t *testing.T) {
	reminders := []scheduling.Reminder{
		{Category: scheduling.CategoryEye, Message: "Look away"},
	}
	title, content := aggregatedNotification(reminders)

	if title == "" {
		t.Errorf("aggregatedNotification title is empty")
	}

	if !strings.Contains(content, "Look away") {
		t.Errorf("aggregatedNotification content = %q, should contain \"Look away\"", content)
	}
}

func TestAggregatedNotificationMultiple(t *testing.T) {
	reminders := []scheduling.Reminder{
		{Category: scheduling.CategoryEye, Message: "Eye care"},
		{Category: scheduling.CategoryStandUp, Message: "Stand up"},
	}
	title, content := aggregatedNotification(reminders)

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

	if !strings.Contains(content, "Eye care") || !strings.Contains(content, "Stand up") {
		t.Errorf("aggregatedNotification content missing expected messages: %q", content)
	}
}
