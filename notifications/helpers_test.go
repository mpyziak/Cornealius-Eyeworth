package notifications

import (
	"slices"
	"testing"

	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func TestPickRandom(t *testing.T) {
	tests := []struct {
		name  string
		items []string
	}{
		{"empty", []string{}},
		{"single", []string{"only"}},
		{"multiple", []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for range 100 {
				got := pickRandom(tt.items)
				if len(tt.items) == 0 {
					if got != "" {
						t.Fatalf("pickRandom(%v) = %q, want empty string", tt.items, got)
					}
					continue
				}
				if !slices.Contains(tt.items, got) {
					t.Fatalf("pickRandom(%v) = %q, want one of %v", tt.items, got, tt.items)
				}
			}
		})
	}
}

func TestHasReminderType(t *testing.T) {
	reminders := []scheduling.Reminder{
		{Category: scheduling.CategoryEye, Message: "test"},
	}

	tests := []struct {
		name      string
		reminders []scheduling.Reminder
		category  scheduling.Category
		want      bool
	}{
		{"empty", nil, scheduling.CategoryEye, false},
		{"present", reminders, scheduling.CategoryEye, true},
		{"absent", reminders, scheduling.CategoryStandUp, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasReminderType(tt.reminders, tt.category); got != tt.want {
				t.Errorf("hasReminderType(%v, %s) = %v, want %v", tt.reminders, tt.category, got, tt.want)
			}
		})
	}
}
