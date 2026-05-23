package notifications

import (
	"strings"
	"testing"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// TestAggregatedContentEmpty verifies aggregatedContent returns empty string for no reminders.
func TestAggregatedContentEmpty(t *testing.T) {
	result := aggregatedContent([]scheduling.Reminder{})
	if result != "" {
		t.Errorf("aggregatedContent([]) = %q, want empty string", result)
	}
}

// TestAggregatedContentSingle verifies aggregatedContent formats a single reminder.
func TestAggregatedContentSingle(t *testing.T) {
	reminders := []scheduling.Reminder{
		{NotificationCategory: scheduling.ReminderTypeEye, Message: "Look away"},
	}
	result := aggregatedContent(reminders)
	expected := "• Look away"
	if result != expected {
		t.Errorf("aggregatedContent single reminder = %q, want %q", result, expected)
	}
}

// TestAggregatedContentMultiple verifies aggregatedContent orders by type, then preserves order.
func TestAggregatedContentMultiple(t *testing.T) {
	reminders := []scheduling.Reminder{
		{NotificationCategory: scheduling.ReminderTypeEye, Message: "Eye care"},
		{NotificationCategory: scheduling.ReminderTypeStandup, Message: "Stand up"},
		{NotificationCategory: scheduling.ReminderTypeEye, Message: "Blink"},
	}
	result := aggregatedContent(reminders)
	lines := strings.Split(result, "\n")

	if len(lines) != 3 {
		t.Errorf("aggregatedContent produced %d lines, want 3", len(lines))
	}

	// Standup should come first in the ordering
	if !strings.Contains(lines[0], "Stand up") {
		t.Errorf("First line = %q, want it to contain standup reminder", lines[0])
	}

	// Eye reminders should follow
	if !strings.Contains(lines[1], "Eye care") {
		t.Errorf("Second line = %q, want \"Eye care\"", lines[1])
	}
	if !strings.Contains(lines[2], "Blink") {
		t.Errorf("Third line = %q, want \"Blink\"", lines[2])
	}
}

// TestAggregatedTitleEyeOnly verifies aggregatedTitle uses DistanceGlanceHeaders for eye-only.
func TestAggregatedTitleEyeOnly(t *testing.T) {
	reminders := []scheduling.Reminder{
		{NotificationCategory: scheduling.ReminderTypeEye, Message: "test"},
	}
	result := aggregatedTitle(reminders)

	found := false
	for _, header := range i18n.Active.NotificationDistanceGlanceHeaders {
		if result == header {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("aggregatedTitle returned %q, which is not in NotificationDistanceGlanceHeaders", result)
	}
}

// TestAggregatedTitleStandupOnly verifies aggregatedTitle uses MovementHeaders for standup-only.
func TestAggregatedTitleStandupOnly(t *testing.T) {
	reminders := []scheduling.Reminder{
		{NotificationCategory: scheduling.ReminderTypeStandup, Message: "test"},
	}
	result := aggregatedTitle(reminders)

	found := false
	for _, header := range i18n.Active.NotificationMovementHeaders {
		if result == header {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("aggregatedTitle returned %q, which is not in NotificationMovementHeaders", result)
	}
}

// TestAggregatedTitleBoth verifies aggregatedTitle uses CombinedHeaders for both types.
func TestAggregatedTitleBoth(t *testing.T) {
	reminders := []scheduling.Reminder{
		{NotificationCategory: scheduling.ReminderTypeEye, Message: "test"},
		{NotificationCategory: scheduling.ReminderTypeStandup, Message: "test"},
	}
	result := aggregatedTitle(reminders)

	found := false
	for _, header := range i18n.Active.NotificationCombinedHeaders {
		if result == header {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("aggregatedTitle returned %q, which is not in NotificationCombinedHeaders", result)
	}
}
