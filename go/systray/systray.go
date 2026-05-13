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
	window          fyne.Window
	scheduleMenu    *systray.MenuItem
	nextTriggerMenu *systray.MenuItem
	windowVisible   bool
	setupCfg        *config.Config
}

// NewManager creates a new systray Manager.
func NewManager(app fyne.App, window fyne.Window) *Manager {
	return &Manager{
		app:    app,
		window: window,
	}
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

	showHideItem := systray.AddMenuItem(S.AppName, "Show/Hide window")
	systray.AddSeparator()

	m.nextTriggerMenu = systray.AddMenuItem(
		fmt.Sprintf(S.NextTrigger, scheduling.NextTrigger(cfg.CronExpression).Format("15:04")),
		"Next trigger time",
	)
	m.nextTriggerMenu.Disable()

	systray.AddSeparator()
	quitItem := systray.AddMenuItem(S.MenuQuit, "Quit application")

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

// toggleWindowVisibility shows or hides the main window.
func (m *Manager) toggleWindowVisibility() {
	if m.windowVisible {
		m.window.Hide()
		m.windowVisible = false
	} else {
		m.window.Show()
		m.windowVisible = true
	}
}

// Run starts the system tray event loop (blocking call).
func (m *Manager) Run() {
	systray.Run(func() {
		m.doSetup()
	}, func() {
		m.app.Quit()
	})
}
