package notifications

import (
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// NewReminder builds the Reminder for notificationCategory, picking its
// message text from the matching i18n quip pool.
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
