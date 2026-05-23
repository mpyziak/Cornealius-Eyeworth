package notifications

import (
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func aggregatedNotification(reminders []scheduling.Reminder) (string, string) {
	if len(reminders) == 0 {
		return "", ""
	}

	return aggregatedTitle(reminders), aggregatedContent(reminders)
}
