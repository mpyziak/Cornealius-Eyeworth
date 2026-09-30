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

// scheduleDeps bundles the values every level of the status-window and
// schedule-dialog wiring needs, so they travel as one parameter instead of
// four repeated positional ones.
type scheduleDeps struct {
	repo         *config.Repository
	eyeSched     scheduling.Runner
	standUpSched scheduling.Runner
	onFire       func(scheduling.Category)
}

// statusLabels are the status window's live label widgets, bundled so they
// travel and get nilled out together rather than as four positional params.
type statusLabels struct {
	schedule, nextTrigger, standUp, standUpNextTrigger *widget.Label
}

// Direct label pointers, not data bindings: Fyne never releases
// binding listeners on window close, so each open/close cycle leaks one.
type statusState struct {
	mu sync.Mutex

	schedule           string
	nextTrigger        string
	standUp            string
	standUpNextTrigger string

	labels *statusLabels // non-nil while the status window is open
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
	labels := s.labels
	s.mu.Unlock()

	// Update live labels if window is currently open.
	if labels != nil {
		labels.schedule.SetText(schedule)
		labels.nextTrigger.SetText(nextTrigger)
		labels.standUp.SetText(standUp)
		labels.standUpNextTrigger.SetText(standUpNextTrigger)
	}
}

func (s *statusState) attach(labels *statusLabels) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.labels = labels
}

func (s *statusState) detach() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.labels = nil
}

func (s *statusState) snapshot() (schedule, nextTrigger, standUp, standUpNextTrigger string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.schedule, s.nextTrigger, s.standUp, s.standUpNextTrigger
}

// refreshLabels bundles the status window and tray menu refreshes so every
// caller re-arms both from one place. A caller that only calls status.update
// leaves the tray showing a stale "Next eye care" time until the next fire.
func refreshLabels(status *statusState, trayMgr *systray.Manager) func(*config.Config) {
	return func(cfg *config.Config) {
		status.update(cfg)
		trayMgr.UpdateLabels(cfg)
	}
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

	deps := scheduleDeps{
		repo:         repo,
		eyeSched:     eyeScheduler,
		standUpSched: standUpScheduler,
		// Every ApplySchedule caller must re-arm with this one. A local
		// replacement drops the status and tray refresh.
		onFire: cronCallback(repo, status, trayMgr, reminderBuffer),
	}

	trayMgr.SetWindowFactory(windowFactory(app, status, trayMgr, deps))

	if err := scheduling.ApplySchedule(
		deps.eyeSched,
		deps.standUpSched,
		scheduling.Spec{
			EyeCron:     cfg.CronExpression,
			StandUpCron: cfg.StandUpCronExpression,
		},
		deps.onFire,
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
			deps.eyeSched,
			deps.standUpSched,
			scheduling.Spec{
				EyeCron:     cfg.CronExpression,
				StandUpCron: cfg.StandUpCronExpression,
			},
			deps.onFire,
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
	status *statusState,
	trayMgr *systray.Manager,
	deps scheduleDeps,
) func() fyne.Window {
	return func() fyne.Window {
		diagnostics.Event("status window opened")
		schedule, nextTrigger, standUp, standUpNextTrigger := status.snapshot()
		win, labels := buildMainWindow(app, schedule, nextTrigger, standUp, standUpNextTrigger)
		status.attach(labels)
		win.SetMainMenu(buildMenu(app, refreshLabels(status, trayMgr), deps))
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
	refresh := refreshLabels(status, trayMgr)
	return func(category scheduling.Category) {
		diagnostics.Event("cron fired - category=%s", category)
		buffer.Add(notifications.NewReminder(category))
		if latest, err := repo.Load(); err != nil {
			diagnostics.Err("cron fired: cannot load config.json (%v) - labels not refreshed", err)
		} else {
			// Cron's goroutine. Both calls touch Fyne widgets, which are
			// main thread only.
			fyne.Do(func() { refresh(latest) })
		}
	}
}

func buildMenu(
	app fyne.App,
	updateStatus func(*config.Config),
	deps scheduleDeps,
) *fyne.MainMenu {
	str := i18n.Active

	quitItem := fyne.NewMenuItem(str.MenuQuit, func() { app.Quit() })
	quitItem.IsQuit = true

	return fyne.NewMainMenu(
		fyne.NewMenu(str.MenuOptions,
			fyne.NewMenuItem(str.MenuTriggerTimes, func() {
				ShowScheduleDialog(app, deps, updateStatus)
			}),
			fyne.NewMenuItem(str.MenuLanguage, func() {
				ShowLanguageDialog(app, deps.repo)
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

func buildMainWindow(app fyne.App, schedule, nextTrigger, standUp, standUpNextTrigger string) (fyne.Window, *statusLabels) {
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
	return win, &statusLabels{
		schedule:           scheduleLabel,
		nextTrigger:        nextTriggerLabel,
		standUp:            standUpLabel,
		standUpNextTrigger: standUpNextTriggerLabel,
	}
}
