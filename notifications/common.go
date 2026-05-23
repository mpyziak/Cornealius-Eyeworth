package notifications

import (
	"math/rand"
	"strings"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

const (
	reminderTypeEye     = "eye"
	reminderTypeStandup = "standup"
)

func pickRandom(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[rand.Intn(len(items))]
}

func hasReminderType(reminders []scheduling.Reminder, kind string) bool {
	for _, r := range reminders {
		if r.Type == kind {
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
	for _, targetType := range []string{reminderTypeStandup, reminderTypeEye} {
		for _, r := range reminders {
			if r.Type == targetType {
				lines = append(lines, "• "+r.Message)
			}
		}
	}

	for _, r := range reminders {
		if r.Type != reminderTypeStandup && r.Type != reminderTypeEye {
			lines = append(lines, "• "+r.Message)
		}
	}

	content := strings.Join(lines, "\n")
	return content
}

func aggregatedTitle(reminders []scheduling.Reminder) string {
	hasStandUp := hasReminderType(reminders, reminderTypeStandup)
	hasDistanceGlance := hasReminderType(reminders, reminderTypeEye)
	if hasStandUp && hasDistanceGlance {
		return pickRandom(i18n.Active.NotificationCombinedHeaders)
	} else if hasStandUp {
		return pickRandom(i18n.Active.NotificationMovementHeaders)
	}
	return pickRandom(i18n.Active.NotificationDistanceGlanceHeaders)
}

func aggregatedNotification(reminders []scheduling.Reminder) (string, string) {
	if len(reminders) == 0 {
		return "", ""
	}

	return aggregatedTitle(reminders), aggregatedContent(reminders)
}
