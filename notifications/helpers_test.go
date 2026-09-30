package notifications

import (
	"testing"

	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func TestPickRandomEmpty(t *testing.T) {
	result := pickRandom([]string{})
	if result != "" {
		t.Errorf("pickRandom([]) = %q, want empty string", result)
	}
}

func TestPickRandomSingle(t *testing.T) {
	items := []string{"only"}
	result := pickRandom(items)
	if result != "only" {
		t.Errorf("pickRandom([\"only\"]) = %q, want \"only\"", result)
	}
}

func TestPickRandomValid(t *testing.T) {
	items := []string{"a", "b", "c"}
	for range 100 {
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

func TestHasReminderTypeEmpty(t *testing.T) {
	result := hasReminderType([]scheduling.Reminder{}, scheduling.CategoryEye)
	if result {
		t.Errorf("hasReminderType([], CategoryEye) = true, want false")
	}
}

func TestHasReminderTypePresent(t *testing.T) {
	reminders := []scheduling.Reminder{
		{Category: scheduling.CategoryEye, Message: "test"},
	}
	result := hasReminderType(reminders, scheduling.CategoryEye)
	if !result {
		t.Errorf("hasReminderType with CategoryEye present = false, want true")
	}
}

func TestHasReminderTypeAbsent(t *testing.T) {
	reminders := []scheduling.Reminder{
		{Category: scheduling.CategoryEye, Message: "test"},
	}
	result := hasReminderType(reminders, scheduling.CategoryStandUp)
	if result {
		t.Errorf("hasReminderType without CategoryStandUp = true, want false")
	}
}
