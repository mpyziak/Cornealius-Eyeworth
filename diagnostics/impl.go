//go:build diagnostics

package diagnostics

// impl.go wires the dev-diagnostics implementation into this package when the
// diagnostics build tag is set. The alias avoids a name collision: the
// dev-diagnostics package is also named "diagnostics".
import (
	"fyne.io/fyne/v2"

	devdiag "github.com/mpyziak/cornealius-eyeworth/dev-diagnostics"
)

func Init(dir string) error            { return devdiag.Init(dir) }
func Close()                           { devdiag.Close() }
func Info(format string, args ...any)  { devdiag.Info(format, args...) }
func Event(format string, args ...any) { devdiag.Event(format, args...) }
func Warn(format string, args ...any)  { devdiag.Warn(format, args...) }
func Err(format string, args ...any)   { devdiag.Err(format, args...) }
func EnableUiDevDiagnosticsSettingsListener(app fyne.App) {
	devdiag.EnableUiDevDiagnosticsSettingsListener(app)
}
