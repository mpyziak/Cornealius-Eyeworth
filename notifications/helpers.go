package notifications

import (
	"math/rand/v2"
	"slices"

	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func pickRandom(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[rand.IntN(len(items))]
}

func hasReminderType(reminders []scheduling.Reminder, category scheduling.Category) bool {
	return slices.ContainsFunc(reminders, func(r scheduling.Reminder) bool {
		return r.Category == category
	})
}
