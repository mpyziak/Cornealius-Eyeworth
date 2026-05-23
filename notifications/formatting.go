package notifications

import (
	"strings"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func aggregatedContent(reminders []scheduling.Reminder) string {
	if len(reminders) == 0 {
		return ""
	}

	var lines []string
	for _, targetNotificationCategory := range []string{scheduling.ReminderTypeStandup, scheduling.ReminderTypeEye} {
		for _, r := range reminders {
			if r.NotificationCategory == targetNotificationCategory {
				lines = append(lines, "• "+r.Message)
			}
		}
	}

	for _, r := range reminders {
		if r.NotificationCategory != scheduling.ReminderTypeStandup && r.NotificationCategory != scheduling.ReminderTypeEye {
			lines = append(lines, "• "+r.Message)
		}
	}

	content := strings.Join(lines, "\n")
	return content
}

func aggregatedTitle(reminders []scheduling.Reminder) string {
	hasStandUp := hasReminderType(reminders, scheduling.ReminderTypeStandup)
	hasDistanceGlance := hasReminderType(reminders, scheduling.ReminderTypeEye)
	if hasStandUp && hasDistanceGlance {
		return pickRandom(i18n.Active.NotificationCombinedHeaders)
	} else if hasStandUp {
		return pickRandom(i18n.Active.NotificationMovementHeaders)
	}
	return pickRandom(i18n.Active.NotificationDistanceGlanceHeaders)
}
