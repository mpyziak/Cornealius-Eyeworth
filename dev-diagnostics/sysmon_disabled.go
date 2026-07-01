//go:build diagnostics && !sysmon

package diagnostics

// No-op stubs used when the binary is built without -tags sysmon

import "fyne.io/fyne/v2"

func startSysmon()                                      {}
func stopSysmon()                                       {}
func EnableUiDevDiagnosticsSettingsListener(_ fyne.App) {}
