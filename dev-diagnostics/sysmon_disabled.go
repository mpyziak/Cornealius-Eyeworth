//go:build diagnostics && !sysmon

package devdiag

// No-op stubs used when the binary is built without -tags sysmon

import "fyne.io/fyne/v2"

func startSysmon() {}
func stopSysmon()  {}

// EnableUIDevDiagnosticsSettingsListener is a no-op without -tags sysmon.
func EnableUIDevDiagnosticsSettingsListener(_ fyne.App) {}
