// Package systray handles system tray integration for Cornealius Eyeworth.
package systray

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/systray"

	"github.com/mpyziak/cornealius-eyeworth/assets"
	"github.com/mpyziak/cornealius-eyeworth/config"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// Manager handles systray setup and lifecycle.
type Manager struct {
	app             fyne.App
	windowFactory   func() fyne.Window // builds a fresh window each time it is shown
	currentWindow   fyne.Window        // non-nil while the status window is visible
	scheduleMenu    *systray.MenuItem
	nextTriggerMenu *systray.MenuItem
	setupCfg        *config.Config
}

// NewManager creates a new systray Manager.
func NewManager(app fyne.App) *Manager {
	return &Manager{app: app}
}

// SetWindowFactory sets the function used to build the status window on demand.
// Must be called before the first tray interaction.
func (m *Manager) SetWindowFactory(f func() fyne.Window) {
	m.windowFactory = f
}

// Setup stores the config for use when the tray is ready.
func (m *Manager) Setup(cfg *config.Config) {
	m.setupCfg = cfg
}

// doSetup performs the real systray initialisation. Must be called from inside
// the onReady callback passed to systray.Run so it works correctly on Windows.
func (m *Manager) doSetup() {
	cfg := m.setupCfg
	S := i18n.Active

	systray.SetIcon(assets.IconBytes())
	systray.SetTooltip(S.AppName)

	showHideItem := systray.AddMenuItem(S.TrayTooltipShowHide, S.AppName)
	systray.AddSeparator()

	m.nextTriggerMenu = systray.AddMenuItem(
		fmt.Sprintf(S.NextTrigger, scheduling.NextTrigger(cfg.CronExpression).Format("15:04")),
		S.TrayTooltipNextTrigger,
	)
	m.nextTriggerMenu.Disable()

	systray.AddSeparator()
	quitItem := systray.AddMenuItem(S.MenuQuit, S.TrayTooltipQuit)

	go func() {
		for {
			select {
			case <-showHideItem.ClickedCh:
				m.toggleWindowVisibility()
			case <-quitItem.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

// UpdateLabels updates the schedule and next trigger labels in the tray menu.
func (m *Manager) UpdateLabels(cfg *config.Config) {
	S := i18n.Active
	if m.scheduleMenu != nil {
		m.scheduleMenu.SetTitle(scheduling.Describe(cfg.CronExpression))
	}
	if m.nextTriggerMenu != nil {
		m.nextTriggerMenu.SetTitle(
			fmt.Sprintf(S.NextTrigger, scheduling.NextTrigger(cfg.CronExpression).Format("15:04")),
		)
	}
}

// NotifyHidden records that the status window has been closed by means other
// than the tray toggle (e.g. the window's own X button or its close intercept),
// so the next tray click correctly opens a fresh window.
func (m *Manager) NotifyHidden() {
	m.currentWindow = nil
}

// toggleWindowVisibility opens a fresh status window when none is shown, or
// closes the current one when it is already visible.
func (m *Manager) toggleWindowVisibility() {
	if m.currentWindow != nil {
		w := m.currentWindow
		m.currentWindow = nil
		w.Close() // destroys HWND + GL context; no hidden window left over
	} else if m.windowFactory != nil {
		m.currentWindow = m.windowFactory()
		m.currentWindow.Show()
	}
}

// Run starts the system tray event loop (blocking call).
// onReady is called after the tray icon has been fully initialised;
// pass nil if no post-init work is needed.
func (m *Manager) Run(onReady func()) {
	systray.Run(func() {
		m.doSetup()
		if onReady != nil {
			onReady()
		}
	}, func() {
		m.app.Quit()
	})
}
