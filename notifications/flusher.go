package notifications

import (
	"fyne.io/fyne/v2"

	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// AppFlusher sends buffered reminders as a notification. Lives here:
// scheduling cannot import this package.
type AppFlusher struct {
	app fyne.App
}

// NewAppFlusher returns an AppFlusher that sends through app.
func NewAppFlusher(app fyne.App) *AppFlusher {
	return &AppFlusher{app: app}
}

// Flush sends reminders as a single aggregated notification.
func (af *AppFlusher) Flush(reminders []scheduling.Reminder) {
	SendReminders(af.app, reminders)
}
