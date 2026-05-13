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
	onScheduleClick func()
	onLanguageClick func()
	onAboutClick    func()
	onHelpClick     func()
	scheduleMenu    *systray.MenuItem
	nextTriggerMenu *systray.MenuItem
	windowVisible   bool
}

// NewManager creates a new systray Manager.
func NewManager(app fyne.App, window fyne.Window) *Manager {
	return &Manager{
		app:    app,
		window: window,
	}
}

// SetCallbacks sets the menu item click handlers.
func (m *Manager) SetCallbacks(
	onSchedule, onLanguage, onAbout, onHelp func(),
) {
	m.onScheduleClick = onSchedule
	m.onLanguageClick = onLanguage
	m.onAboutClick = onAbout
	m.onHelpClick = onHelp
}

// Setup initializes the system tray with menu items.
func (m *Manager) Setup(cfg *config.Config) {
	S := i18n.Active

	// Set tray icon
	systray.SetIcon(assets.IconBytes())
	systray.SetTooltip(S.AppName)

	// Create menu items
	showHideItem := systray.AddMenuItem(S.AppName, "Show/Hide window")
	systray.AddSeparator()

	m.scheduleMenu = systray.AddMenuItem(scheduling.Describe(cfg.CronExpression), "Current schedule")
	m.nextTriggerMenu = systray.AddMenuItem(
		fmt.Sprintf(S.NextTrigger, scheduling.NextTrigger(cfg.CronExpression).Format("15:04")),
		"Next trigger time",
	)
	m.scheduleMenu.Disable()
	m.nextTriggerMenu.Disable()

	systray.AddSeparator()

	settingsItem := systray.AddMenuItem(S.MenuTriggerTimes, "Configure schedule")
	languageItem := systray.AddMenuItem(S.MenuLanguage, "Change language")
	helpItem := systray.AddMenuItem(S.MenuHowToUse, "View help")
	aboutItem := systray.AddMenuItem(S.MenuAbout, "About application")

	systray.AddSeparator()
	quitItem := systray.AddMenuItem(S.MenuQuit, "Quit application")

	// Handle menu clicks
	go func() {
		for {
			select {
			case <-showHideItem.ClickedCh:
				m.toggleWindowVisibility()
			case <-settingsItem.ClickedCh:
				if m.onScheduleClick != nil {
					m.onScheduleClick()
				}
			case <-languageItem.ClickedCh:
				if m.onLanguageClick != nil {
					m.onLanguageClick()
				}
			case <-helpItem.ClickedCh:
				if m.onHelpClick != nil {
					m.onHelpClick()
				}
			case <-aboutItem.ClickedCh:
				if m.onAboutClick != nil {
					m.onAboutClick()
				}
			case <-quitItem.ClickedCh:
				m.app.Quit()
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
		// onReady callback - tray icon is ready
	}, func() {
		// onExit callback
	})
}
