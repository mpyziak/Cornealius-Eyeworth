//go:build !diagnostics

// Release builds compile only this file - every call is a no-op.
package diagnostics

import "fyne.io/fyne/v2"

// Init is a no-op in release builds.
func Init(_ string) error { return nil }

// Close is a no-op in release builds.
func Close() {}

// Info is a no-op in release builds.
func Info(_ string, _ ...any) {}

// Event is a no-op in release builds.
func Event(_ string, _ ...any) {}

// Warn is a no-op in release builds.
func Warn(_ string, _ ...any) {}

// Err is a no-op in release builds.
func Err(_ string, _ ...any) {}

// EnableUIDevDiagnosticsSettingsListener is a no-op in release builds.
func EnableUIDevDiagnosticsSettingsListener(_ fyne.App) {}
