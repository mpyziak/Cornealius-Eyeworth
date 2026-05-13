//go:build !windows

// Package notifications sends cross-platform desktop notifications via Fyne.
package notifications

import (
	"math/rand"

	"fyne.io/fyne/v2"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
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

// SendReminder delivers an eye-rest reminder notification.
func SendReminder(app fyne.App) {
	app.SendNotification(&fyne.Notification{
		Title:   pickRandom(i18n.Active.NotificationReminders),
		Content: pickRandom(i18n.Active.NotificationQuips),
	})
}
