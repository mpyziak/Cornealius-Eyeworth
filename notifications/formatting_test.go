package notifications

import (
	"strings"
	"testing"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func TestAggregatedContentEmpty(t *testing.T) {
	result := aggregatedContent([]scheduling.Reminder{})
	if result != "" {
		t.Errorf("aggregatedContent([]) = %q, want empty string", result)
	}
}

func TestAggregatedContentSingle(t *testing.T) {
	reminders := []scheduling.Reminder{
		{Category: scheduling.CategoryEye, Message: "Look away"},
	}
	result := aggregatedContent(reminders)
	expected := "• Look away"
	if result != expected {
		t.Errorf("aggregatedContent single reminder = %q, want %q", result, expected)
	}
}

func TestAggregatedContentMultiple(t *testing.T) {
	reminders := []scheduling.Reminder{
		{Category: scheduling.CategoryEye, Message: "Eye care"},
		{Category: scheduling.CategoryStandUp, Message: "Stand up"},
		{Category: scheduling.CategoryEye, Message: "Blink"},
	}
	result := aggregatedContent(reminders)
	lines := strings.Split(result, "\n")

	if len(lines) != 3 {
		t.Errorf("aggregatedContent produced %d lines, want 3", len(lines))
	}

	if !strings.Contains(lines[0], "Stand up") {
		t.Errorf("First line = %q, want it to contain standup reminder", lines[0])
	}

	if !strings.Contains(lines[1], "Eye care") {
		t.Errorf("Second line = %q, want \"Eye care\"", lines[1])
	}
	if !strings.Contains(lines[2], "Blink") {
		t.Errorf("Third line = %q, want \"Blink\"", lines[2])
	}
}

func TestAggregatedTitleEyeOnly(t *testing.T) {
	reminders := []scheduling.Reminder{
		{Category: scheduling.CategoryEye, Message: "test"},
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

func TestAggregatedTitleStandupOnly(t *testing.T) {
	reminders := []scheduling.Reminder{
		{Category: scheduling.CategoryStandUp, Message: "test"},
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

func TestAggregatedTitleBoth(t *testing.T) {
	reminders := []scheduling.Reminder{
		{Category: scheduling.CategoryEye, Message: "test"},
		{Category: scheduling.CategoryStandUp, Message: "test"},
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
