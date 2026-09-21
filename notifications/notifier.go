//go:build !windows

package notifications

import (
	"fyne.io/fyne/v2"
	log "github.com/mpyziak/cornealius-eyeworth/diagnostics"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

func SendStartup(app fyne.App) {
	log.Event("notification: startup")
	app.SendNotification(&fyne.Notification{
		Title:   i18n.Active.NotificationOnDuty,
		Content: pickRandom(i18n.Active.NotificationDistanceGlanceQuips),
	})
}

func SendMinimizedToTray(app fyne.App) {
	log.Event("notification: minimized-to-tray")
	app.SendNotification(&fyne.Notification{
		Title:   i18n.Active.AppName,
		Content: i18n.Active.NotificationMinimizedToTray,
	})
}

func SendReminder(app fyne.App) {
	log.Event("notification: reminder (eye)")
	app.SendNotification(&fyne.Notification{
		Title:   pickRandom(i18n.Active.NotificationDistanceGlanceHeaders),
		Content: pickRandom(i18n.Active.NotificationDistanceGlanceQuips),
	})
}

func SendReminders(app fyne.App, reminders []scheduling.Reminder) {
	title, content := aggregatedNotification(reminders)
	if title == "" && content == "" {
		return
	}
	log.Event("notification: reminder batch - count=%d", len(reminders))
	app.SendNotification(&fyne.Notification{
		Title:   title,
		Content: content,
	})
}
