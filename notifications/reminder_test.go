package notifications

import (
	"slices"
	"testing"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func TestNewReminder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		category scheduling.Category
		quips    []string
	}{
		{"eye", scheduling.CategoryEye, i18n.Active.NotificationDistanceGlanceQuips},
		{"standup", scheduling.CategoryStandUp, i18n.Active.NotificationMovementQuips},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reminder := NewReminder(tt.category)

			if reminder.Category != tt.category {
				t.Errorf("NewReminder(%s).Category = %q, want %q", tt.category, reminder.Category, tt.category)
			}
			if !slices.Contains(tt.quips, reminder.Message) {
				t.Errorf("NewReminder(%s) message %q not in %v", tt.category, reminder.Message, tt.quips)
			}
		})
	}
}

func TestNewReminderUnknown(t *testing.T) {
	t.Parallel()
	reminder := NewReminder("unknown")

	if reminder.Category != "unknown" {
		t.Errorf("NewReminder(\"unknown\").Category = %q, want \"unknown\"", reminder.Category)
	}

	if reminder.Message != "" {
		t.Errorf("NewReminder(\"unknown\").Message = %q, want empty string", reminder.Message)
	}
}
