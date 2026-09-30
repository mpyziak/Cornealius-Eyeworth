package notifications

import (
	"testing"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func TestNewReminderEye(t *testing.T) {
	reminder := NewReminder(scheduling.CategoryEye)

	if reminder.Category != scheduling.CategoryEye {
		t.Errorf("NewReminder(CategoryEye).Category = %q, want %q",
			reminder.Category, scheduling.CategoryEye)
	}

	found := false
	for _, quip := range i18n.Active.NotificationDistanceGlanceQuips {
		if reminder.Message == quip {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("NewReminder(CategoryEye) message %q not in NotificationDistanceGlanceQuips", reminder.Message)
	}
}

func TestNewReminderStandup(t *testing.T) {
	reminder := NewReminder(scheduling.CategoryStandUp)

	if reminder.Category != scheduling.CategoryStandUp {
		t.Errorf("NewReminder(CategoryStandUp).Category = %q, want %q",
			reminder.Category, scheduling.CategoryStandUp)
	}

	found := false
	for _, quip := range i18n.Active.NotificationMovementQuips {
		if reminder.Message == quip {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("NewReminder(CategoryStandUp) message %q not in NotificationMovementQuips", reminder.Message)
	}
}

func TestNewReminderUnknown(t *testing.T) {
	reminder := NewReminder("unknown")

	if reminder.Category != "unknown" {
		t.Errorf("NewReminder(\"unknown\").Category = %q, want \"unknown\"", reminder.Category)
	}

	if reminder.Message != "" {
		t.Errorf("NewReminder(\"unknown\").Message = %q, want empty string", reminder.Message)
	}
}
