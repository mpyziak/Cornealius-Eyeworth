package systray

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/mpyziak/cornealius-eyeworth/config"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// Manager handles systray setup and lifecycle.
type Manager struct {
	app             fyne.App
	windowFactory   func() fyne.Window // builds a fresh window each time it is shown
	currentWindow   fyne.Window        // non-nil while the status window is visible
	menu            *fyne.Menu
	nextTriggerMenu *fyne.MenuItem
	setupCfg        *config.Config
}

// NewManager creates a new systray Manager for app.
func NewManager(app fyne.App) *Manager {
	return &Manager{app: app}
}

// SetWindowFactory sets the function used to build a fresh status window
// each time it is shown. Must be called before the first tray interaction.
func (m *Manager) SetWindowFactory(f func() fyne.Window) {
	m.windowFactory = f
}

// Setup builds and installs the tray menu from cfg.
func (m *Manager) Setup(cfg *config.Config) {
	m.setupCfg = cfg
	m.doSetup()
}

func (m *Manager) doSetup() {
	desk, ok := m.app.(desktop.App)
	if !ok {
		return
	}

	cfg := m.setupCfg
	str := i18n.Active

	m.nextTriggerMenu = fyne.NewMenuItem(
		fmt.Sprintf(str.NextTrigger, scheduling.NextTrigger(cfg.CronExpression).Format("15:04")),
		nil,
	)
	m.nextTriggerMenu.Disabled = true

	showHideItem := fyne.NewMenuItem(str.TrayTooltipShowHide, func() {
		m.toggleWindowVisibility()
	})

	quitItem := fyne.NewMenuItem(str.MenuQuit, func() {
		m.app.Quit()
	})
	quitItem.IsQuit = true // prevents Fyne from injecting a second Quit entry

	m.menu = fyne.NewMenu(str.AppName,
		showHideItem,
		fyne.NewMenuItemSeparator(),
		m.nextTriggerMenu,
		fyne.NewMenuItemSeparator(),
		quitItem,
	)

	// The app icon must already be set before this call. Fyne's systray onReady
	// reads fyne.CurrentApp().Icon() to build the tray icon, and it runs on its
	// own goroutine released from inside SetSystemTrayMenu
	desk.SetSystemTrayMenu(m.menu)
}

// UpdateLabels refreshes the tray menu's next-trigger text from cfg. Must
// be called on the main goroutine: fyne.MenuItem.Label is a plain struct
// field with no synchronisation.
func (m *Manager) UpdateLabels(cfg *config.Config) {
	str := i18n.Active
	if m.nextTriggerMenu != nil {
		m.nextTriggerMenu.Label = fmt.Sprintf(str.NextTrigger, scheduling.NextTrigger(cfg.CronExpression).Format("15:04"))
	}
	// A field write alone never reaches the native tray menu - Fyne only reads
	// each item's Label once, when the menu is (re)built. Refresh re-pushes it.
	if m.menu != nil {
		m.menu.Refresh()
	}
}

// NotifyHidden tells the Manager the status window is no longer visible.
func (m *Manager) NotifyHidden() {
	m.currentWindow = nil
}

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
