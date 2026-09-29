//go:build diagnostics && !windows && sysmon

package diagnostics

import "fyne.io/fyne/v2"

func startSysmon() {}
func stopSysmon()  {}

// EnableUIDevDiagnosticsSettingsListener is a no-op on non-Windows: its
// purpose is to correlate with the theme-registry watcher, which is
// Windows-only.
func EnableUIDevDiagnosticsSettingsListener(_ fyne.App) {}
