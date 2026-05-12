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
)

// Run is the application entry point for the UI layer. It creates the main
// window, wires the scheduler, and hands control to the Fyne event loop
// (blocking call — must be invoked from main goroutine).
func Run(app fyne.App, cfg *config.Config, repo *config.Repository) {
	S := i18n.Active

	app.SetIcon(assets.Logo)

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

	notifications.SendStartup(app)
	sched.Start(cfg.CronExpression, func() {
		notifications.SendReminder(app)
		latest, err := repo.Load()
		if err == nil {
			updateBindings(latest)
		}
	})

	win.ShowAndRun()
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

	// ── menu ────────────────────────────────────────────────────────────────
	scheduleItem := fyne.NewMenuItem(S.MenuTriggerTimes, func() {
		ShowScheduleDialog(win, repo, func(updated *config.Config) {
			sched.Start(updated.CronExpression, func() {
				notifications.SendReminder(app)
				if latest, err := repo.Load(); err == nil {
					updateBindings(latest)
				}
			})
			updateBindings(updated)
		})
	})
	languageItem := fyne.NewMenuItem(S.MenuLanguage, func() {
		ShowLanguageDialog(win, repo)
	})
	optionsMenu := fyne.NewMenu(S.MenuOptions, scheduleItem, languageItem)

	aboutItem := fyne.NewMenuItem(S.MenuAbout, func() {
		ShowAboutDialog(win)
	})
	howToItem := fyne.NewMenuItem(S.MenuHowToUse, func() {
		ShowHelpDialog(win)
	})
	githubItem := fyne.NewMenuItem(S.MenuGitHub, func() {
		_ = app.OpenURL(parseURL(S.GitHubUrl))
	})
	helpMenu := fyne.NewMenu(S.MenuHelp, aboutItem, howToItem, githubItem)

	win.SetMainMenu(fyne.NewMainMenu(optionsMenu, helpMenu))

	// ── body ────────────────────────────────────────────────────────────────
	logo := canvas.NewImageFromResource(assets.Logo)
	logo.SetMinSize(fyne.NewSize(36, 36))
	logo.FillMode = canvas.ImageFillContain

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
	win.SetCloseIntercept(func() {
		sched.Stop()
		win.Close()
	})

	return win
}

// parseURL is a small helper to convert a raw URL string; the GitHubUrl
// constant is always well-formed so the error is safely ignored here.
func parseURL(raw string) *url.URL {
	u, _ := url.Parse(raw)
	return u
}
