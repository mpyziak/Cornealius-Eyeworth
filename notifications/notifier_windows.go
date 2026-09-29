//go:build windows

// Shell_NotifyIcon rather than fyne.App.SendNotification: that one shells out to
// PowerShell, which AV blocks on locked-down machines.
package notifications

import (
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"

	"github.com/mpyziak/cornealius-eyeworth/diagnostics"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// NOTIFYICONDATAW on amd64, cbSize 976.
type notifyIconData struct {
	Size            uint32
	_               [4]byte
	Wnd             uintptr
	ID              uint32
	Flags           uint32
	CallbackMessage uint32
	_               [4]byte
	Icon            uintptr
	Tip             [128]uint16
	State           uint32
	StateMask       uint32
	Info            [256]uint16
	Version         uint32
	InfoTitle       [64]uint16
	InfoFlags       uint32
	GuidItem        [16]byte
	BalloonIcon     uintptr
}

const (
	nimModify = 1
	nifInfo   = 0x00000010

	systrayClass  = "SystrayClass"
	systrayIconID = 100
)

var (
	modShell32      = syscall.NewLazyDLL("shell32.dll")
	modUser32       = syscall.NewLazyDLL("user32.dll")
	procShellNotify = modShell32.NewProc("Shell_NotifyIconW")
	procFindWindow  = modUser32.NewProc("FindWindowW")
)

func msgOnlyWindow() uintptr {
	cls, _ := syscall.UTF16PtrFromString(systrayClass)
	hwnd, _, _ := procFindWindow.Call(uintptr(unsafe.Pointer(cls)), 0)
	return hwnd
}

func showBalloon(title, message string) {
	hwnd := msgOnlyWindow()
	if hwnd == 0 {
		fyne.LogError("notification: systray window not found", nil)
		return
	}

	var nid notifyIconData
	nid.Size = uint32(unsafe.Sizeof(nid))
	nid.Wnd = hwnd
	nid.ID = systrayIconID
	nid.Flags = nifInfo
	nid.InfoFlags = 0
	t16, _ := syscall.UTF16FromString(title)
	copy(nid.InfoTitle[:], t16)
	m16, _ := syscall.UTF16FromString(message)
	copy(nid.Info[:], m16)
	procShellNotify.Call(nimModify, uintptr(unsafe.Pointer(&nid)))
}

func SendStartup(_ fyne.App) {
	diagnostics.Event("notification: startup")
	showBalloon(i18n.Active.NotificationOnDuty, pickRandom(i18n.Active.NotificationDistanceGlanceQuips))
}

func SendMinimizedToTray(_ fyne.App) {
	diagnostics.Event("notification: minimized-to-tray")
	showBalloon(i18n.Active.AppName, i18n.Active.NotificationMinimizedToTray)
}

func SendReminder(_ fyne.App) {
	diagnostics.Event("notification: reminder (eye)")
	showBalloon(pickRandom(i18n.Active.NotificationDistanceGlanceHeaders), pickRandom(i18n.Active.NotificationDistanceGlanceQuips))
}

func SendReminders(_ fyne.App, reminders []scheduling.Reminder) {
	title, content := aggregatedNotification(reminders)
	if title == "" && content == "" {
		return
	}
	diagnostics.Event("notification: reminder batch - count=%d", len(reminders))
	showBalloon(title, content)
}
