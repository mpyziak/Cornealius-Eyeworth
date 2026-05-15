//go:build !windows

// Package notifications sends cross-platform desktop notifications via Fyne.
package notifications

import (
	"math/rand"

	"fyne.io/fyne/v2"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func pickRandom(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[rand.Intn(len(items))]
}

// SendStartup delivers the "on duty" notification shown when the app starts.
func SendStartup(app fyne.App) {
	app.SendNotification(&fyne.Notification{
		Title:   i18n.Active.NotificationOnDuty,
		Content: pickRandom(i18n.Active.NotificationQuips),
	})
}

// SendMinimizedToTray delivers a "still running" hint when the window is hidden.
func SendMinimizedToTray(app fyne.App) {
	app.SendNotification(&fyne.Notification{
		Title:   i18n.Active.AppName,
		Content: i18n.Active.NotificationMinimizedToTray,
	})
}

// SendReminder delivers an eye-rest reminder notification.
func SendReminder(app fyne.App) {
	app.SendNotification(&fyne.Notification{
		Title:   pickRandom(i18n.Active.NotificationReminders),
		Content: pickRandom(i18n.Active.NotificationQuips),
	})
}

// SendReminders delivers multiple aggregated reminders in a single notification.
// Implements scheduling.Flusher interface (Dependency Inversion).
func SendReminders(app fyne.App, reminders []scheduling.Reminder) {
	if len(reminders) == 0 {
		return
	}

	// Build combined title and content from reminder types
	var title string
	var content string

	for _, r := range reminders {
		if content != "" {
			content += "\n"
		}
		content += "• " + r.Message
	}

	// Determine title based on reminder types present
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

	app.SendNotification(&fyne.Notification{
		Title:   title,
		Content: content,
	})
}
