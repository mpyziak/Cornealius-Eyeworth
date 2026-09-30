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
	for _, category := range []scheduling.Category{scheduling.CategoryStandUp, scheduling.CategoryEye} {
		for _, r := range reminders {
			if r.Category == category {
				lines = append(lines, "• "+r.Message)
			}
		}
	}

	for _, r := range reminders {
		if r.Category != scheduling.CategoryStandUp && r.Category != scheduling.CategoryEye {
			lines = append(lines, "• "+r.Message)
		}
	}

	content := strings.Join(lines, "\n")
	return content
}

func aggregatedTitle(reminders []scheduling.Reminder) string {
	hasStandUp := hasReminderType(reminders, scheduling.CategoryStandUp)
	hasDistanceGlance := hasReminderType(reminders, scheduling.CategoryEye)
	if hasStandUp && hasDistanceGlance {
		return pickRandom(i18n.Active.NotificationCombinedHeaders)
	} else if hasStandUp {
		return pickRandom(i18n.Active.NotificationMovementHeaders)
	}
	return pickRandom(i18n.Active.NotificationDistanceGlanceHeaders)
}
