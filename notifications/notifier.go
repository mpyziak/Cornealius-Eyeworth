//go:build !windows

package notifications

import (
	"fyne.io/fyne/v2"

	"github.com/mpyziak/cornealius-eyeworth/diagnostics"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// SendStartup notifies that the app has started.
func SendStartup(app fyne.App) {
	diagnostics.Event("notification: startup")
	app.SendNotification(&fyne.Notification{
		Title:   i18n.Active.NotificationOnDuty,
		Content: pickRandom(i18n.Active.NotificationDistanceGlanceQuips),
	})
}

// SendMinimizedToTray notifies that the window was minimized to the tray.
func SendMinimizedToTray(app fyne.App) {
	diagnostics.Event("notification: minimized-to-tray")
	app.SendNotification(&fyne.Notification{
		Title:   i18n.Active.AppName,
		Content: i18n.Active.NotificationMinimizedToTray,
	})
}

// SendReminder sends a single eye-care reminder.
func SendReminder(app fyne.App) {
	diagnostics.Event("notification: reminder (eye)")
	app.SendNotification(&fyne.Notification{
		Title:   pickRandom(i18n.Active.NotificationDistanceGlanceHeaders),
		Content: pickRandom(i18n.Active.NotificationDistanceGlanceQuips),
	})
}

// SendReminders sends every buffered reminder as one aggregated
// notification.
func SendReminders(app fyne.App, reminders []scheduling.Reminder) {
	title, content := aggregatedNotification(reminders)
	if title == "" && content == "" {
		return
	}
	diagnostics.Event("notification: reminder batch - count=%d", len(reminders))
	app.SendNotification(&fyne.Notification{
		Title:   title,
		Content: content,
	})
}
