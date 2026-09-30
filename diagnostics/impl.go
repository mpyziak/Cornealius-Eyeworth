//go:build diagnostics

package diagnostics

import (
	"fyne.io/fyne/v2"

	devdiag "github.com/mpyziak/cornealius-eyeworth/dev-diagnostics"
)

// Init starts the rotating-file logger and memory sampler under dir. Call
// Close on shutdown, else the sampler keeps running.
func Init(dir string) error { return devdiag.Init(dir) }

// Close stops the logger and memory sampler. Safe to call more than once.
func Close() { devdiag.Close() }

// Info logs an informational line.
func Info(format string, args ...any) { devdiag.Info(format, args...) }

// Event logs a notable but non-error occurrence.
func Event(format string, args ...any) { devdiag.Event(format, args...) }

// Warn logs a recoverable problem.
func Warn(format string, args ...any) { devdiag.Warn(format, args...) }

// Err logs an error. Does not exit.
func Err(format string, args ...any) { devdiag.Err(format, args...) }

// EnableUIDevDiagnosticsSettingsListener wires a Fyne settings-change
// listener into the log, for correlating UI-visible theme changes with the
// registry writes and system events logged elsewhere.
func EnableUIDevDiagnosticsSettingsListener(app fyne.App) {
	devdiag.EnableUIDevDiagnosticsSettingsListener(app)
}
