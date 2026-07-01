//go:build !diagnostics

// Package diagnostics is the stable, always-present logging interface for
// Cornealius Eyeworth. In release builds this file is the only one compiled:
// every call is a no-op — no files, goroutines, or OS handles are created.
//
// The optional implementation lives in dev-diagnostics/ and is wired in by
// impl.go when the diagnostics build tag is set.
package diagnostics

import "fyne.io/fyne/v2"

func Init(_ string) error                               { return nil }
func Close()                                            {}
func Info(_ string, _ ...any)                           {}
func Event(_ string, _ ...any)                          {}
func Warn(_ string, _ ...any)                           {}
func Err(_ string, _ ...any)                            {}
func EnableUiDevDiagnosticsSettingsListener(_ fyne.App) {}
