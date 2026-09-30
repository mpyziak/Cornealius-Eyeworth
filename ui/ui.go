package ui

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mpyziak/cornealius-eyeworth/assets"
	"github.com/mpyziak/cornealius-eyeworth/config"
	"github.com/mpyziak/cornealius-eyeworth/diagnostics"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/notifications"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
	"github.com/mpyziak/cornealius-eyeworth/systray"
)

const (
	mainWinW float32 = 420
	mainWinH float32 = 180
)

// Direct label pointers, not data bindings: Fyne never releases
// binding listeners on window close, so each open/close cycle leaks one.
type statusState struct {
	mu sync.Mutex

	schedule           string
	nextTrigger        string
	standUp            string
	standUpNextTrigger string

	// non-nil while the status window is open
	scheduleLabel           *widget.Label
	nextTriggerLabel        *widget.Label
	standUpLabel            *widget.Label
	standUpNextTriggerLabel *widget.Label
}

func (s *statusState) update(c *config.Config) {
	str := i18n.Active
	schedule := scheduling.Describe(c.CronExpression)
	nextTrigger := fmt.Sprintf(str.NextTrigger, scheduling.NextTrigger(c.CronExpression).Format("15:04"))
	standUp := scheduling.Describe(c.StandUpCronExpression)
	standUpNextTrigger := fmt.Sprintf(str.NextTriggerStandUp, scheduling.NextTrigger(c.StandUpCronExpression).Format("15:04"))

	s.mu.Lock()
	s.schedule = schedule
	s.nextTrigger = nextTrigger
	s.standUp = standUp
	s.standUpNextTrigger = standUpNextTrigger
	sl := s.scheduleLabel
	ntl := s.nextTriggerLabel
	sul := s.standUpLabel
	suntl := s.standUpNextTriggerLabel
	s.mu.Unlock()

	// Update live labels if window is currently open.
	if sl != nil {
		sl.SetText(schedule)
		ntl.SetText(nextTrigger)
		sul.SetText(standUp)
		suntl.SetText(standUpNextTrigger)
	}
}

func (s *statusState) attach(sl, ntl, sul, suntl *widget.Label) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scheduleLabel = sl
	s.nextTriggerLabel = ntl
	s.standUpLabel = sul
	s.standUpNextTriggerLabel = suntl
}

func (s *statusState) detach() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scheduleLabel = nil
	s.nextTriggerLabel = nil
	s.standUpLabel = nil
	s.standUpNextTriggerLabel = nil
}

func (s *statusState) snapshot() (schedule, nextTrigger, standUp, standUpNextTrigger string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.schedule, s.nextTrigger, s.standUp, s.standUpNextTrigger
}

// Run wires up the schedulers, tray, and status window, and blocks in the
// Fyne event loop until the app quits.
func Run(app fyne.App, cfg *config.Config, repo *config.Repository) {
	status := &statusState{}
	status.update(cfg)

	// Anchor window. Never shown - it only keeps the event loop alive once the
	// status window closes. Showing it would register an HWND that receives
	// DPI-change and screen-capture events.
	_ = app.NewWindow("")

	eyeScheduler := &scheduling.Scheduler{}
	standUpScheduler := &scheduling.Scheduler{}
	reminderBuffer := scheduling.NewBuffer(2*time.Second, notifications.NewAppFlusher(app))

	trayMgr := systray.NewManager(app)
	trayMgr.Setup(cfg)

	// Every ApplySchedule caller must re-arm with this one. A local replacement
	// drops the status and tray refresh.
	onFire := cronCallback(repo, status, trayMgr, reminderBuffer)

	trayMgr.SetWindowFactory(windowFactory(app, repo, status, eyeScheduler, standUpScheduler, trayMgr, onFire))

	if err := scheduling.ApplySchedule(
		eyeScheduler,
		standUpScheduler,
		scheduling.Spec{
			EyeCron:     cfg.CronExpression,
			StandUpCron: cfg.StandUpCronExpression,
		},
		onFire,
	); err != nil {
		// A cron with no entries looks exactly like a working app that never
		// fires, so fall back rather than run with empty.
		diagnostics.Err("invalid schedule in config.json (%v) - falling back to defaults", err)
		defaults := config.DefaultConfig()
		cfg.CronExpression = defaults.CronExpression
		cfg.StandUpCronExpression = defaults.StandUpCronExpression
		status.update(cfg)
		trayMgr.UpdateLabels(cfg)
		if fallbackErr := scheduling.ApplySchedule(
			eyeScheduler,
			standUpScheduler,
			scheduling.Spec{
				EyeCron:     cfg.CronExpression,
				StandUpCron: cfg.StandUpCronExpression,
			},
			onFire,
		); fallbackErr != nil {
			diagnostics.Err("default schedule rejected too: %v", fallbackErr)
		}
	}

	diagnostics.EnableUIDevDiagnosticsSettingsListener(app)

	trayMgr.Run(func() { notifications.SendStartup(app) })
	diagnostics.Info("event loop running - eye=%s standUp=%s", cfg.CronExpression, cfg.StandUpCronExpression)
	app.Run()

	eyeScheduler.Stop()
	standUpScheduler.Stop()
	reminderBuffer.Close()
	diagnostics.Info("schedulers stopped, shutting down")
}

// A fresh window per show, closed rather than hidden. Hiding keeps the HWND
// alive, and Windows keeps sending it DPI and screen-share events that Fyne
// answers by reallocating GL buffers.
func windowFactory(
	app fyne.App,
	repo *config.Repository,
	status *statusState,
	eyeSched scheduling.Runner,
	standUpSched scheduling.Runner,
	trayMgr *systray.Manager,
	onFire func(scheduling.Category),
) func() fyne.Window {
	return func() fyne.Window {
		diagnostics.Event("status window opened")
		schedule, nextTrigger, standUp, standUpNextTrigger := status.snapshot()
		win, sl, ntl, sul, suntl := buildMainWindow(app, schedule, nextTrigger, standUp, standUpNextTrigger)
		status.attach(sl, ntl, sul, suntl)
		win.SetMainMenu(buildMenu(app, repo, eyeSched, standUpSched, status.update, onFire))
		win.SetOnClosed(func() {
			diagnostics.Event("status window closed")
			status.detach()
			trayMgr.NotifyHidden()
		})
		win.SetCloseIntercept(func() {
			notifications.SendMinimizedToTray(app)
			win.Close() // triggers SetOnClosed -> detach + NotifyHidden
		})
		return win
	}
}

func cronCallback(
	repo *config.Repository,
	status *statusState,
	trayMgr *systray.Manager,
	buffer *scheduling.Buffer,
) func(scheduling.Category) {
	return func(category scheduling.Category) {
		diagnostics.Event("cron fired - category=%s", category)
		buffer.Add(notifications.NewReminder(category))
		if latest, err := repo.Load(); err != nil {
			diagnostics.Err("cron fired: cannot load config.json (%v) - labels not refreshed", err)
		} else {
			// Cron's goroutine. Both calls touch Fyne widgets, which are
			// main thread only.
			fyne.Do(func() {
				status.update(latest)
				trayMgr.UpdateLabels(latest)
			})
		}
	}
}

func buildMenu(
	app fyne.App,
	repo *config.Repository,
	eyeSched scheduling.Runner,
	standUpSched scheduling.Runner,
	updateStatus func(*config.Config),
	onFire func(scheduling.Category),
) *fyne.MainMenu {
	str := i18n.Active

	quitItem := fyne.NewMenuItem(str.MenuQuit, func() { app.Quit() })
	quitItem.IsQuit = true

	return fyne.NewMainMenu(
		fyne.NewMenu(str.MenuOptions,
			fyne.NewMenuItem(str.MenuTriggerTimes, func() {
				ShowScheduleDialog(app, repo, eyeSched, standUpSched, updateStatus, onFire)
			}),
			fyne.NewMenuItem(str.MenuLanguage, func() {
				ShowLanguageDialog(app, repo)
			}),
			fyne.NewMenuItemSeparator(),
			quitItem,
		),
		fyne.NewMenu(str.MenuHelp,
			fyne.NewMenuItem(str.MenuHowToUse, func() { ShowHelpDialog(app) }),
			fyne.NewMenuItem(str.MenuAbout, func() { ShowAboutDialog(app) }),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem(str.MenuGitHub, func() {
				// GitHubURL is a constant well-formed URL; url.Parse cannot fail on it.
				u, _ := url.Parse(str.GitHubURL)
				_ = app.OpenURL(u) // best-effort; no channel to report a failed browser launch to
			}),
		),
	)
}

func buildMainWindow(app fyne.App, schedule, nextTrigger, standUp, standUpNextTrigger string) (fyne.Window, *widget.Label, *widget.Label, *widget.Label, *widget.Label) {
	str := i18n.Active

	win := app.NewWindow(str.AppName)
	win.Resize(fyne.NewSize(mainWinW, mainWinH))

	logo := canvas.NewImageFromResource(assets.Logo)
	logo.SetMinSize(fyne.NewSize(36, 36))
	logo.FillMode = canvas.ImageFillContain
	logo.ScaleMode = canvas.ImageScaleSmooth

	titleLabel := widget.NewLabelWithStyle(str.AppTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	statusLabel := widget.NewLabel(str.StatusServing)

	scheduleLabel := widget.NewLabel(schedule)
	scheduleLabel.Wrapping = fyne.TextWrapWord

	nextTriggerLabel := widget.NewLabel(nextTrigger)

	standUpLabel := widget.NewLabel(standUp)
	standUpLabel.Wrapping = fyne.TextWrapWord

	standUpNextTriggerLabel := widget.NewLabel(standUpNextTrigger)

	topRow := container.NewHBox(logo, container.NewVBox(titleLabel, statusLabel))
	body := container.New(layout.NewVBoxLayout(),
		topRow,
		widget.NewSeparator(),
		nextTriggerLabel,
		scheduleLabel,
		widget.NewSeparator(),
		standUpNextTriggerLabel,
		standUpLabel,
	)

	win.SetContent(container.NewPadded(body))
	return win, scheduleLabel, nextTriggerLabel, standUpLabel, standUpNextTriggerLabel
}
