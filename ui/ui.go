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
	S := i18n.Active
	schedule := scheduling.Describe(c.CronExpression)
	nextTrigger := fmt.Sprintf(S.NextTrigger, scheduling.NextTrigger(c.CronExpression).Format("15:04"))
	standUp := scheduling.Describe(c.StandUpCronExpression)
	standUpNextTrigger := fmt.Sprintf(S.NextTriggerStandUp, scheduling.NextTrigger(c.StandUpCronExpression).Format("15:04"))

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

func Run(app fyne.App, cfg *config.Config, repo config.Store) {
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
		scheduling.ScheduleSpec{
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
			scheduling.ScheduleSpec{
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
	repo config.Store,
	status *statusState,
	eyeSched scheduling.Runner,
	standUpSched scheduling.Runner,
	trayMgr *systray.Manager,
	onFire func(string),
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
	repo config.Store,
	status *statusState,
	trayMgr *systray.Manager,
	buffer scheduling.ReminderAggregator,
) func(string) {
	return func(notificationCategory string) {
		diagnostics.Event("cron fired - category=%s", notificationCategory)
		buffer.Add(notifications.NewReminder(notificationCategory))
		if latest, err := repo.Load(); err == nil {
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
	repo config.Store,
	eyeSched scheduling.Runner,
	standUpSched scheduling.Runner,
	updateStatus func(*config.Config),
	onFire func(string),
) *fyne.MainMenu {
	S := i18n.Active

	quitItem := fyne.NewMenuItem(S.MenuQuit, func() { app.Quit() })
	quitItem.IsQuit = true

	return fyne.NewMainMenu(
		fyne.NewMenu(S.MenuOptions,
			fyne.NewMenuItem(S.MenuTriggerTimes, func() {
				ShowScheduleDialog(app, repo, eyeSched, standUpSched, updateStatus, onFire)
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
				u, _ := url.Parse(S.GitHubURL)
				_ = app.OpenURL(u)
			}),
		),
	)
}

func buildMainWindow(app fyne.App, schedule, nextTrigger, standUp, standUpNextTrigger string) (fyne.Window, *widget.Label, *widget.Label, *widget.Label, *widget.Label) {
	S := i18n.Active

	win := app.NewWindow(S.AppName)
	win.Resize(fyne.NewSize(mainWinW, mainWinH))

	logo := canvas.NewImageFromResource(assets.Logo)
	logo.SetMinSize(fyne.NewSize(36, 36))
	logo.FillMode = canvas.ImageFillContain
	logo.ScaleMode = canvas.ImageScaleSmooth

	titleLabel := widget.NewLabelWithStyle(S.AppTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	statusLabel := widget.NewLabel(S.StatusServing)

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
