//go:build diagnostics && !windows && sysmon

package diagnostics

import "fyne.io/fyne/v2"

func startSysmon() {}
func stopSysmon()  {}

// noop on non-Windows: listener's purpose is to correlate with
// the theme-registry watcher, which is Windows-only
func EnableUIDevDiagnosticsSettingsListener(_ fyne.App) {}
