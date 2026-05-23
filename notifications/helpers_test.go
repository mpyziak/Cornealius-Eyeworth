package notifications

import (
	"testing"

	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// TestPickRandomEmpty verifies pickRandom returns empty string for empty list.
func TestPickRandomEmpty(t *testing.T) {
	result := pickRandom([]string{})
	if result != "" {
		t.Errorf("pickRandom([]) = %q, want empty string", result)
	}
}

// TestPickRandomSingle verifies pickRandom returns the only item.
func TestPickRandomSingle(t *testing.T) {
	items := []string{"only"}
	result := pickRandom(items)
	if result != "only" {
		t.Errorf("pickRandom([\"only\"]) = %q, want \"only\"", result)
	}
}

// TestPickRandomValid verifies pickRandom returns an item from the list.
func TestPickRandomValid(t *testing.T) {
	items := []string{"a", "b", "c"}
	for i := 0; i < 100; i++ {
		result := pickRandom(items)
		found := false
		for _, item := range items {
			if result == item {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("pickRandom returned %q, which is not in the input list", result)
		}
	}
}

// TestHasReminderTypeEmpty verifies hasReminderType returns false for empty list.
func TestHasReminderTypeEmpty(t *testing.T) {
	result := hasReminderType([]scheduling.Reminder{}, scheduling.ReminderTypeEye)
	if result {
		t.Errorf("hasReminderType([], ReminderTypeEye) = true, want false")
	}
}

// TestHasReminderTypePresent verifies hasReminderType finds a matching type.
func TestHasReminderTypePresent(t *testing.T) {
	reminders := []scheduling.Reminder{
		{NotificationCategory: scheduling.ReminderTypeEye, Message: "test"},
	}
	result := hasReminderType(reminders, scheduling.ReminderTypeEye)
	if !result {
		t.Errorf("hasReminderType with ReminderTypeEye present = false, want true")
	}
}

// TestHasReminderTypeAbsent verifies hasReminderType returns false when type is missing.
func TestHasReminderTypeAbsent(t *testing.T) {
	reminders := []scheduling.Reminder{
		{NotificationCategory: scheduling.ReminderTypeEye, Message: "test"},
	}
	result := hasReminderType(reminders, scheduling.ReminderTypeStandup)
	if result {
		t.Errorf("hasReminderType without ReminderTypeStandup = true, want false")
	}
}
