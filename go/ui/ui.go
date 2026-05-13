// Package ui contains all Fyne UI code for Cornealius Eyeworth.
package ui

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mpyziak/cornealius-eyeworth/assets"
	"github.com/mpyziak/cornealius-eyeworth/config"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/notifications"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
	"github.com/mpyziak/cornealius-eyeworth/systray"
)

// mainWinW/H is the fixed size of the main status window.
// Every dialog restores the parent to these dimensions on close so Fyne
// doesn't leave it at the dialog's larger size.
const (
	mainWinW float32 = 420
	mainWinH float32 = 180
)

// Run is the application entry point for the UI layer. It creates the main
// window, sets up the system tray, and hands control to the tray event loop.
func Run(app fyne.App, cfg *config.Config, repo *config.Repository) {
	S := i18n.Active

	app.SetIcon(assets.AppIcon)

	scheduleBinding := binding.NewString()
	nextTriggerBinding := binding.NewString()

	updateBindings := func(c *config.Config) {
		_ = scheduleBinding.Set(scheduling.Describe(c.CronExpression))
		_ = nextTriggerBinding.Set(
			fmt.Sprintf(S.NextTrigger, scheduling.NextTrigger(c.CronExpression).Format("15:04")),
		)
	}
	updateBindings(cfg)

	sched := &scheduling.Scheduler{}

	win := buildMainWindow(app, cfg, repo, sched, scheduleBinding, nextTriggerBinding, updateBindings)
	win.SetCloseIntercept(func() {
		win.Hide()
		notifications.SendMinimizedToTray(app)
	})
	trayMgr := systray.NewManager(app, win)
	trayMgr.Setup(cfg)

	// Build the window menu bar — all dialogs open from here so the window
	// is always visible when they appear; no resize gymnastics needed.
	S2 := i18n.Active // alias to avoid shadowing the outer S

	quitItem := fyne.NewMenuItem(S2.MenuQuit, func() { app.Quit() })
	quitItem.IsQuit = true

	win.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu(S2.MenuOptions,
			fyne.NewMenuItem(S2.MenuTriggerTimes, func() {
				ShowScheduleDialog(app, repo, func(updated *config.Config) {
					sched.Start(updated.CronExpression, func() {
						notifications.SendReminder(app)
						if latest, err := repo.Load(); err == nil {
							updateBindings(latest)
							trayMgr.UpdateLabels(latest)
						}
					})
					updateBindings(updated)
					trayMgr.UpdateLabels(updated)
				})
			}),
			fyne.NewMenuItem(S2.MenuLanguage, func() {
				ShowLanguageDialog(app, repo)
			}),
			fyne.NewMenuItemSeparator(),
			quitItem,
		),
		fyne.NewMenu(S2.MenuHelp,
			fyne.NewMenuItem(S2.MenuHowToUse, func() {
				ShowHelpDialog(app)
			}),
			fyne.NewMenuItem(S2.MenuAbout, func() {
				ShowAboutDialog(app)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem(S2.MenuGitHub, func() {
				u, _ := url.Parse(S2.GitHubUrl)
				_ = app.OpenURL(u)
			}),
		),
	))

	notifications.SendStartup(app)
	sched.Start(cfg.CronExpression, func() {
		notifications.SendReminder(app)
		latest, err := repo.Load()
		if err == nil {
			updateBindings(latest)
			trayMgr.UpdateLabels(latest)
		}
	})

	// Hide the window initially (systray is the primary interface)
	win.Hide()

	// Run the system tray in a goroutine; fyne.io/systray manages its own
	// OS thread on Windows, so this is safe.
	go trayMgr.Run()

	// Start Fyne's event loop on the main goroutine (blocks until app.Quit()).
	app.Run()

	// Cleanup when app exits
	sched.Stop()
}

func buildMainWindow(
	app fyne.App,
	cfg *config.Config,
	repo *config.Repository,
	sched *scheduling.Scheduler,
	scheduleBinding binding.String,
	nextTriggerBinding binding.String,
	updateBindings func(*config.Config),
) fyne.Window {
	S := i18n.Active

	win := app.NewWindow(S.AppName)
	win.Resize(fyne.NewSize(mainWinW, mainWinH))

	// ── body ────────────────────────────────────────────────────────────────
	logo := canvas.NewImageFromResource(assets.Logo)
	logo.SetMinSize(fyne.NewSize(36, 36))
	logo.FillMode = canvas.ImageFillContain
	logo.ScaleMode = canvas.ImageScaleSmooth

	titleLabel := widget.NewLabelWithStyle(S.AppTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	statusLabel := widget.NewLabel(S.StatusServing)

	scheduleLabel := widget.NewLabelWithData(scheduleBinding)
	scheduleLabel.Wrapping = fyne.TextWrapWord

	nextTriggerLabel := widget.NewLabelWithData(nextTriggerBinding)

	topRow := container.NewHBox(logo, container.NewVBox(
		titleLabel,
		statusLabel,
	))

	body := container.New(layout.NewVBoxLayout(),
		topRow,
		widget.NewSeparator(),
		scheduleLabel,
		nextTriggerLabel,
	)

	win.SetContent(container.NewPadded(body))

	return win
}
