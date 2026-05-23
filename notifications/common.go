package notifications

import (
	"math/rand"
	"strings"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func pickRandom(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[rand.Intn(len(items))]
}

func hasReminderType(reminders []scheduling.Reminder, notificationCategory string) bool {
	for _, r := range reminders {
		if r.NotificationCategory == notificationCategory {
			return true
		}
	}
	return false
}

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

func NewReminder(notificationCategory string) scheduling.Reminder {
	switch notificationCategory {
	case scheduling.ReminderTypeStandup:
		return scheduling.Reminder{NotificationCategory: notificationCategory, Message: pickRandom(i18n.Active.NotificationMovementQuips)}
	case scheduling.ReminderTypeEye:
		return scheduling.Reminder{NotificationCategory: notificationCategory, Message: pickRandom(i18n.Active.NotificationDistanceGlanceQuips)}
	default:
		return scheduling.Reminder{NotificationCategory: notificationCategory, Message: ""}
	}
}

func aggregatedNotification(reminders []scheduling.Reminder) (string, string) {
	if len(reminders) == 0 {
		return "", ""
	}

	return aggregatedTitle(reminders), aggregatedContent(reminders)
}
