// Package systray handles system tray integration for Cornealius Eyeworth.
package systray

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"

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
	scheduleMenu    *fyne.MenuItem
	nextTriggerMenu *fyne.MenuItem
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

// Setup stores the config and performs the real systray initialisation.
func (m *Manager) Setup(cfg *config.Config) {
	m.setupCfg = cfg
	m.doSetup()
}

// doSetup performs the real systray initialisation.
func (m *Manager) doSetup() {
	desk, ok := m.app.(desktop.App)
	if !ok {
		return
	}

	cfg := m.setupCfg
	S := i18n.Active

	m.nextTriggerMenu = fyne.NewMenuItem(
		fmt.Sprintf(S.NextTrigger, scheduling.NextTrigger(cfg.CronExpression).Format("15:04")),
		nil,
	)
	m.nextTriggerMenu.Disabled = true

	showHideItem := fyne.NewMenuItem(S.TrayTooltipShowHide, func() {
		m.toggleWindowVisibility()
	})

	quitItem := fyne.NewMenuItem(S.MenuQuit, func() {
		m.app.Quit()
	})
	quitItem.IsQuit = true // prevents Fyne from injecting a second Quit entry

	menu := fyne.NewMenu(S.AppName,
		showHideItem,
		fyne.NewMenuItemSeparator(),
		m.nextTriggerMenu,
		fyne.NewMenuItemSeparator(),
		quitItem,
	)

	desk.SetSystemTrayMenu(menu)
	m.app.SetIcon(assets.Logo)
}

// UpdateLabels updates the schedule and next trigger labels in the tray menu.
func (m *Manager) UpdateLabels(cfg *config.Config) {
	S := i18n.Active
	if m.scheduleMenu != nil {
		m.scheduleMenu.Label = scheduling.Describe(cfg.CronExpression)
	}
	if m.nextTriggerMenu != nil {
		m.nextTriggerMenu.Label = fmt.Sprintf(S.NextTrigger, scheduling.NextTrigger(cfg.CronExpression).Format("15:04"))
	}
	// Fyne's MenuItem automatically reflects changes to its Label if the menu is active.
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

// Run is kept for API compatibility but no longer starts a separate loop.
func (m *Manager) Run(onReady func()) {
	patchTrayWindows()
	if onReady != nil {
		onReady()
	}
}
