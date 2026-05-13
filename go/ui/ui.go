// Package ui contains all Fyne UI code for Cornealius Eyeworth.
package ui

import (
	"fmt"

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

	// Create and setup system tray manager
	trayMgr := systray.NewManager(app, win)
	trayMgr.Setup(cfg)
	trayMgr.SetCallbacks(
		func() {
			ShowScheduleDialog(win, repo, func(updated *config.Config) {
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
		},
		func() {
			ShowLanguageDialog(win, repo)
		},
		func() {
			ShowAboutDialog(win)
		},
		func() {
			ShowHelpDialog(win)
		},
	)

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
	// win.ShowAndRun()
	win.Hide()

	// Run the system tray event loop (blocking call)
	trayMgr.Run()

	// Cleanup when tray exits
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
	win.SetFixedSize(true)
	win.Resize(fyne.NewSize(420, 180))

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
	win.SetOnClosed(func() {
		// Hide instead of closing to keep the app running in systray
		win.Hide()
	})

	return win
}
