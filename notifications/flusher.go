package notifications

import (
	"fyne.io/fyne/v2"

	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// Lives here: scheduling cannot import this package.
type AppFlusher struct {
	app fyne.App
}

func NewAppFlusher(app fyne.App) *AppFlusher {
	return &AppFlusher{app: app}
}

func (af *AppFlusher) Flush(reminders []scheduling.Reminder) {
	SendReminders(af.app, reminders)
}
