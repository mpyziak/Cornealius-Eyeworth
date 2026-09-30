package notifications

import (
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// NewReminder builds the Reminder for category, picking its message text
// from the matching i18n quip pool.
func NewReminder(category scheduling.Category) scheduling.Reminder {
	switch category {
	case scheduling.CategoryStandUp:
		return scheduling.Reminder{Category: category, Message: pickRandom(i18n.Active.NotificationMovementQuips)}
	case scheduling.CategoryEye:
		return scheduling.Reminder{Category: category, Message: pickRandom(i18n.Active.NotificationDistanceGlanceQuips)}
	default:
		return scheduling.Reminder{Category: category, Message: ""}
	}
}
