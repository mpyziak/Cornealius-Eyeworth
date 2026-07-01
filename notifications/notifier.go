//go:build !windows

// Package notifications sends cross-platform desktop notifications via Fyne.
package notifications

import (
	"fyne.io/fyne/v2"
	log "github.com/mpyziak/cornealius-eyeworth/diagnostics"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// SendStartup delivers the "on duty" notification shown when the app starts.
func SendStartup(app fyne.App) {
	log.Event("notification: startup")
	app.SendNotification(&fyne.Notification{
		Title:   i18n.Active.NotificationOnDuty,
		Content: pickRandom(i18n.Active.NotificationDistanceGlanceQuips),
	})
}

// SendMinimizedToTray delivers a "still running" hint when the window is hidden.
func SendMinimizedToTray(app fyne.App) {
	log.Event("notification: minimized-to-tray")
	app.SendNotification(&fyne.Notification{
		Title:   i18n.Active.AppName,
		Content: i18n.Active.NotificationMinimizedToTray,
	})
}

// SendReminder delivers an eye-rest reminder notification.
func SendReminder(app fyne.App) {
	log.Event("notification: reminder (eye)")
	app.SendNotification(&fyne.Notification{
		Title:   pickRandom(i18n.Active.NotificationDistanceGlanceHeaders),
		Content: pickRandom(i18n.Active.NotificationDistanceGlanceQuips),
	})
}

// SendReminders delivers multiple aggregated reminders in a single notification.
// Implements scheduling.Flusher interface (Dependency Inversion).
func SendReminders(app fyne.App, reminders []scheduling.Reminder) {
	title, content := aggregatedNotification(reminders)
	if title == "" && content == "" {
		return
	}
	log.Event("notification: reminder batch — count=%d", len(reminders))
	app.SendNotification(&fyne.Notification{
		Title:   title,
		Content: content,
	})
}
