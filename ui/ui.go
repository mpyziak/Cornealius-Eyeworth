// Package ui contains all Fyne UI code for Cornealius Eyeworth.
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
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/notifications"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
	"github.com/mpyziak/cornealius-eyeworth/systray"
)

const (
	mainWinW float32 = 420
	mainWinH float32 = 180
)

// statusState holds the current displayable status strings and, while a status
// window is open, direct label pointers for live updates. Using direct label
// references instead of fyne data bindings avoids accumulating binding
// listeners from each window-open cycle — binding listeners are never cleaned
// up automatically when a window is closed, causing a listener-per-cycle leak
// that prevents old Label widgets from being garbage-collected.
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

	// Update live labels if the window is currently open.
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

// Run is the application entry point for the UI layer.
func Run(app fyne.App, cfg *config.Config, repo config.Store) {
	app.SetIcon(assets.AppIcon)

	status := &statusState{}
	status.update(cfg)

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
		schedule, nextTrigger, standUp, standUpNextTrigger := status.snapshot()
		win, sl, ntl, sul, suntl := buildMainWindow(app, schedule, nextTrigger, standUp, standUpNextTrigger)
		status.attach(sl, ntl, sul, suntl)
		win.SetMainMenu(buildMenu(app, repo, eyeScheduler, standUpScheduler, status.update, trayMgr, reminderBuffer))
		win.SetOnClosed(func() {
			status.detach()
			trayMgr.NotifyHidden()
		})
		win.SetCloseIntercept(func() {
			notifications.SendMinimizedToTray(app)
			win.Close() // triggers SetOnClosed → detach + NotifyHidden
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
				status.update(latest)
				trayMgr.UpdateLabels(latest)
			}
		},
	)

	trayMgr.Run(func() { notifications.SendStartup(app) })
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
// It returns the window and the four status labels so the caller can register
// them for live updates without using data bindings.
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
