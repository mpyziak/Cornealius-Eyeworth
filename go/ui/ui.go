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

updateStatus := func(c *config.Config) {
_ = scheduleBinding.Set(scheduling.Describe(c.CronExpression))
_ = nextTriggerBinding.Set(
fmt.Sprintf(S.NextTrigger, scheduling.NextTrigger(c.CronExpression).Format("15:04")),
)
}
updateStatus(cfg)

sched := &scheduling.Scheduler{}

win := buildMainWindow(app, scheduleBinding, nextTriggerBinding)
win.SetCloseIntercept(func() {
win.Hide()
notifications.SendMinimizedToTray(app)
})

trayMgr := systray.NewManager(app, win)
trayMgr.Setup(cfg)
win.SetMainMenu(buildMenu(app, repo, sched, updateStatus, trayMgr))

notifications.SendStartup(app)
sched.Start(cfg.CronExpression, func() {
notifications.SendReminder(app)
if latest, err := repo.Load(); err == nil {
updateStatus(latest)
trayMgr.UpdateLabels(latest)
}
})

win.Hide()
go trayMgr.Run()
app.Run()
sched.Stop()
}

// buildMenu constructs the window menu bar.
// All dialogs are opened from menu items so the parent window is always
// visible when they appear — no resize juggling needed.
func buildMenu(
app fyne.App,
repo config.Store,
sched scheduling.Runner,
updateStatus func(*config.Config),
trayMgr *systray.Manager,
) *fyne.MainMenu {
S := i18n.Active

quitItem := fyne.NewMenuItem(S.MenuQuit, func() { app.Quit() })
quitItem.IsQuit = true

return fyne.NewMainMenu(
fyne.NewMenu(S.MenuOptions,
fyne.NewMenuItem(S.MenuTriggerTimes, func() {
ShowScheduleDialog(app, repo, func(updated *config.Config) {
sched.Start(updated.CronExpression, func() {
notifications.SendReminder(app)
if latest, err := repo.Load(); err == nil {
updateStatus(latest)
trayMgr.UpdateLabels(latest)
}
})
updateStatus(updated)
trayMgr.UpdateLabels(updated)
})
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
func buildMainWindow(app fyne.App, scheduleBinding, nextTriggerBinding binding.String) fyne.Window {
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

topRow := container.NewHBox(logo, container.NewVBox(titleLabel, statusLabel))
body := container.New(layout.NewVBoxLayout(),
topRow,
widget.NewSeparator(),
scheduleLabel,
nextTriggerLabel,
)

win.SetContent(container.NewPadded(body))
return win
}