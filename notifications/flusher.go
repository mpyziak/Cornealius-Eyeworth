package notifications

import (
	"fyne.io/fyne/v2"
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

// Flush sends aggregated reminders as a single notification.
func (af *AppFlusher) Flush(reminders []scheduling.Reminder) {
	SendReminders(af.app, reminders)
}
