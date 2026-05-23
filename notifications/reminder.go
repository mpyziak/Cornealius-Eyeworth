package notifications

import (
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

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
