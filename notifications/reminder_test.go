package notifications

import (
	"testing"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func TestNewReminderEye(t *testing.T) {
	reminder := NewReminder(scheduling.ReminderTypeEye)

	if reminder.NotificationCategory != scheduling.ReminderTypeEye {
		t.Errorf("NewReminder(ReminderTypeEye).NotificationCategory = %q, want %q",
			reminder.NotificationCategory, scheduling.ReminderTypeEye)
	}

	found := false
	for _, quip := range i18n.Active.NotificationDistanceGlanceQuips {
		if reminder.Message == quip {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("NewReminder(ReminderTypeEye) message %q not in NotificationDistanceGlanceQuips", reminder.Message)
	}
}

func TestNewReminderStandup(t *testing.T) {
	reminder := NewReminder(scheduling.ReminderTypeStandup)

	if reminder.NotificationCategory != scheduling.ReminderTypeStandup {
		t.Errorf("NewReminder(ReminderTypeStandup).NotificationCategory = %q, want %q",
			reminder.NotificationCategory, scheduling.ReminderTypeStandup)
	}

	found := false
	for _, quip := range i18n.Active.NotificationMovementQuips {
		if reminder.Message == quip {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("NewReminder(ReminderTypeStandup) message %q not in NotificationMovementQuips", reminder.Message)
	}
}

func TestNewReminderUnknown(t *testing.T) {
	reminder := NewReminder("unknown")

	if reminder.NotificationCategory != "unknown" {
		t.Errorf("NewReminder(\"unknown\").NotificationCategory = %q, want \"unknown\"", reminder.NotificationCategory)
	}

	if reminder.Message != "" {
		t.Errorf("NewReminder(\"unknown\").Message = %q, want empty string", reminder.Message)
	}
}
