package notifications

import (
	"math/rand/v2"

	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func pickRandom(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[rand.IntN(len(items))]
}

func hasReminderType(reminders []scheduling.Reminder, notificationCategory string) bool {
	for _, r := range reminders {
		if r.NotificationCategory == notificationCategory {
			return true
		}
	}
	return false
}
