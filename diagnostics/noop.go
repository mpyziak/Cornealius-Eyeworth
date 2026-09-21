//go:build !diagnostics

// Release builds compile only this file - every call is a no-op.
package diagnostics

import "fyne.io/fyne/v2"

func Init(_ string) error                               { return nil }
func Close()                                            {}
func Info(_ string, _ ...any)                           {}
func Event(_ string, _ ...any)                          {}
func Warn(_ string, _ ...any)                           {}
func Err(_ string, _ ...any)                            {}
func EnableUiDevDiagnosticsSettingsListener(_ fyne.App) {}
