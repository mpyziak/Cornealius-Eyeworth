//go:build windows

package notifications

import (
	"fyne.io/fyne/v2"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// AppFlusher adapts a Fyne app to the scheduling.Flusher interface (Adapter pattern).
// This is the glue between scheduling and notifications domains (Dependency Inversion).
type AppFlusher struct {
	app fyne.App
}

// NewAppFlusher creates a flusher bound to the given app.
func NewAppFlusher(app fyne.App) *AppFlusher {
	return &AppFlusher{app: app}
}

// Flush sends aggregated reminders as a single Windows toast notification.
func (af *AppFlusher) Flush(reminders []scheduling.Reminder) {
	if len(reminders) == 0 {
		return
	}

	var title string
	var content string

	for _, r := range reminders {
		if content != "" {
			content += "\n"
		}
		content += "* " + r.Message
	}

	hasEye := false
	hasStandUp := false
	for _, r := range reminders {
		if r.Type == "eye" {
			hasEye = true
		} else if r.Type == "standup" {
			hasStandUp = true
		}
	}

	if hasEye && hasStandUp {
		title = i18n.Active.NotificationCombinedTitle
	} else if hasStandUp {
		title = pickRandom(i18n.Active.NotificationStandUpReminders)
	} else {
		title = pickRandom(i18n.Active.NotificationReminders)
	}

	go sendToast(af.app, title, content)
}
