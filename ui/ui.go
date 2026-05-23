// Package ui contains all Fyne UI code for Cornealius Eyeworth.
package ui

import (
	"fmt"
	"net/url"
	"time"

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

const (
	mainWinW float32 = 420
	mainWinH float32 = 180
)

// Run is the application entry point for the UI layer.
func Run(app fyne.App, cfg *config.Config, repo config.Store) {
	S := i18n.Active

	app.SetIcon(assets.AppIcon)

	scheduleBinding := binding.NewString()
	nextTriggerBinding := binding.NewString()
	standUpBinding := binding.NewString()
	standUpNextTriggerBinding := binding.NewString()

	updateStatus := func(c *config.Config) {
		_ = scheduleBinding.Set(scheduling.Describe(c.CronExpression))
		_ = nextTriggerBinding.Set(
			fmt.Sprintf(S.NextTrigger, scheduling.NextTrigger(c.CronExpression).Format("15:04")),
		)
		_ = standUpBinding.Set(scheduling.Describe(c.StandUpCronExpression))
		_ = standUpNextTriggerBinding.Set(
			fmt.Sprintf(S.NextTriggerStandUp, scheduling.NextTrigger(c.StandUpCronExpression).Format("15:04")),
		)
	}
	updateStatus(cfg)

	// Anchor window: created first so it becomes Fyne's "master" window.
	// It is never shown, so Windows never registers its HWND as a visible
	// window and will not send DPI-change / screen-capture events to it.
	// Its sole purpose is to keep Fyne's event loop alive after the
	// on-demand status window is closed.
	_ = app.NewWindow("")

	// Create scheduler instances and reminder buffer
	eyeScheduler := &scheduling.Scheduler{}
	standUpScheduler := &scheduling.Scheduler{}
	reminderBuffer := scheduling.NewBuffer(2*time.Second, notifications.NewAppFlusher(app))

	trayMgr := systray.NewManager(app)
	trayMgr.Setup(cfg)

	// windowFactory builds a fresh status window each time the user shows it.
	// Closing the window calls win.Close() which destroys the HWND and frees
	// the Fyne/GLFW OpenGL context, preventing GL resource accumulation that
	// is triggered by screen-share and DPI-change OS events.
	trayMgr.SetWindowFactory(func() fyne.Window {
		win := buildMainWindow(app, scheduleBinding, nextTriggerBinding, standUpBinding, standUpNextTriggerBinding)
		win.SetMainMenu(buildMenu(app, repo, eyeScheduler, standUpScheduler, updateStatus, trayMgr, reminderBuffer))
		win.SetOnClosed(func() {
			trayMgr.NotifyHidden()
		})
		win.SetCloseIntercept(func() {
			notifications.SendMinimizedToTray(app)
			win.Close() // triggers SetOnClosed → NotifyHidden
		})
		return win
	})

	// Start both schedulers using scheduling package helper.
	scheduling.ApplySchedule(
		eyeScheduler,
		standUpScheduler,
		scheduling.ScheduleSpec{
			EyeCron:     cfg.CronExpression,
			StandUpCron: cfg.StandUpCronExpression,
		},
		func(notificationCategory string) {
			reminderBuffer.Add(notifications.NewReminder(notificationCategory))
			if latest, err := repo.Load(); err == nil {
				updateStatus(latest)
				trayMgr.UpdateLabels(latest)
			}
		},
	)

	go trayMgr.Run(func() { notifications.SendStartup(app) })
	app.Run()

	// Cleanup
	eyeScheduler.Stop()
	standUpScheduler.Stop()
	reminderBuffer.Close()
}

// buildMenu constructs the window menu bar.
// All dialogs are opened from menu items so the parent window is always
// visible when they appear — no resize juggling needed.
func buildMenu(
	app fyne.App,
	repo config.Store,
	eyeSched scheduling.Runner,
	standUpSched scheduling.Runner,
	updateStatus func(*config.Config),
	trayMgr *systray.Manager,
	buffer scheduling.ReminderAggregator,
) *fyne.MainMenu {
	S := i18n.Active

	quitItem := fyne.NewMenuItem(S.MenuQuit, func() { app.Quit() })
	quitItem.IsQuit = true

	return fyne.NewMainMenu(
		fyne.NewMenu(S.MenuOptions,
			fyne.NewMenuItem(S.MenuTriggerTimes, func() {
				ShowScheduleDialog(app, repo, eyeSched, standUpSched, updateStatus, buffer)
			}),
			fyne.NewMenuItem(S.MenuLanguage, func() {
				ShowLanguageDialog(app, repo)
			}),
			fyne.NewMenuItemSeparator(),
			quitItem,
		),
		fyne.NewMenu(S.MenuHelp,
			fyne.NewMenuItem(S.MenuHowToUse, func() { ShowHelpDialog(app) }),
			fyne.NewMenuItem(S.MenuAbout, func() { ShowAboutDialog(app) }),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem(S.MenuGitHub, func() {
				u, _ := url.Parse(S.GitHubUrl)
				_ = app.OpenURL(u)
			}),
		),
	)
}

// buildMainWindow creates and configures the main status window.
func buildMainWindow(app fyne.App, scheduleBinding, nextTriggerBinding, standUpBinding, standUpNextTriggerBinding binding.String) fyne.Window {
	S := i18n.Active

	win := app.NewWindow(S.AppName)
	win.Resize(fyne.NewSize(mainWinW, mainWinH))

	logo := canvas.NewImageFromResource(assets.Logo)
	logo.SetMinSize(fyne.NewSize(36, 36))
	logo.FillMode = canvas.ImageFillContain
	logo.ScaleMode = canvas.ImageScaleSmooth

	titleLabel := widget.NewLabelWithStyle(S.AppTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	statusLabel := widget.NewLabel(S.StatusServing)

	scheduleLabel := widget.NewLabelWithData(scheduleBinding)
	scheduleLabel.Wrapping = fyne.TextWrapWord

	nextTriggerLabel := widget.NewLabelWithData(nextTriggerBinding)

	standUpLabel := widget.NewLabelWithData(standUpBinding)
	standUpLabel.Wrapping = fyne.TextWrapWord

	standUpNextTriggerLabel := widget.NewLabelWithData(standUpNextTriggerBinding)

	topRow := container.NewHBox(logo, container.NewVBox(titleLabel, statusLabel))
	body := container.New(layout.NewVBoxLayout(),
		topRow,
		widget.NewSeparator(),
		scheduleLabel,
		nextTriggerLabel,
		widget.NewSeparator(),
		standUpLabel,
		standUpNextTriggerLabel,
	)

	win.SetContent(container.NewPadded(body))
	return win
}
