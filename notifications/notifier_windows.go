//go:build windows

// Package notifications sends Windows balloon tip notifications via Shell_NotifyIcon.
// This bypasses fyne.App.SendNotification(), which shells out to PowerShell and
// can be blocked by AV scanners on some machines.
package notifications

import (
	"math/rand"
	"os"
	"sync"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"
	"github.com/mpyziak/cornealius-eyeworth/assets"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// notifyIconData mirrors NOTIFYICONDATAW on amd64 (cbSize = 976 bytes).
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
	nimModify   = 1
	nimDelete   = 2
	nifInfo     = 0x00000010
	niifUser    = 0x00000004
	niifLarge   = 0x00000020
	hwndMessage = ^uintptr(2)

	// systrayClass is the window class registered by fyne.io/systray.
	// Its tray icon has ID 100. We piggyback on it for balloon delivery so
	// Windows shows the app icon in the notification header.
	systrayClass  = "SystrayClass"
	systrayIconID = 100

	// LoadImage flags
	imageIcon      = 1
	lrLoadFromFile = 0x00000010
	lrDefaultSize  = 0x00000040
)

var (
	modShell32      = syscall.NewLazyDLL("shell32.dll")
	modUser32       = syscall.NewLazyDLL("user32.dll")
	procShellNotify = modShell32.NewProc("Shell_NotifyIconW")
	procFindWindow  = modUser32.NewProc("FindWindowW")
	procLoadImage   = modUser32.NewProc("LoadImageW")

	appIconHwnd  uintptr
	appIconOnce  sync.Once
	iconTempFile string
)

func msgOnlyWindow() uintptr {
	cls, _ := syscall.UTF16PtrFromString(systrayClass)
	hwnd, _, _ := procFindWindow.Call(uintptr(unsafe.Pointer(cls)), 0)
	return hwnd
}

// appIcon writes the embedded .ico to a temp file once and loads it as an HICON.
// Falls back to IDI_APPLICATION if anything fails.
func appIcon() uintptr {
	appIconOnce.Do(func() {
		tmp, err := os.CreateTemp("", "eyeworth-*.ico")
		if err != nil {
			return
		}
		iconTempFile = tmp.Name()
		if _, err := tmp.Write(assets.IconBytes()); err != nil {
			tmp.Close()
			return
		}
		tmp.Close()
		path16, err := syscall.UTF16PtrFromString(iconTempFile)
		if err != nil {
			return
		}
		hicon, _, _ := procLoadImage.Call(
			0,
			uintptr(unsafe.Pointer(path16)),
			imageIcon,
			0, 0,
			lrLoadFromFile|lrDefaultSize,
		)
		appIconHwnd = hicon
	})
	return appIconHwnd
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
	nid.BalloonIcon = appIcon()
	nid.InfoFlags = niifUser | niifLarge
	t16, _ := syscall.UTF16FromString(title)
	copy(nid.InfoTitle[:], t16)
	m16, _ := syscall.UTF16FromString(message)
	copy(nid.Info[:], m16)
	procShellNotify.Call(nimModify, uintptr(unsafe.Pointer(&nid)))
}

func pickRandom(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[rand.Intn(len(items))]
}

func SendStartup(_ fyne.App) {
	showBalloon(i18n.Active.NotificationOnDuty, pickRandom(i18n.Active.NotificationQuips))
}

func SendMinimizedToTray(_ fyne.App) {
	showBalloon(i18n.Active.AppName, i18n.Active.NotificationMinimizedToTray)
}

func SendReminder(_ fyne.App) {
	showBalloon(pickRandom(i18n.Active.NotificationReminders), pickRandom(i18n.Active.NotificationQuips))
}

func SendReminders(_ fyne.App, reminders []scheduling.Reminder) {
	if len(reminders) == 0 {
		return
	}
	var content string
	for _, r := range reminders {
		if content != "" {
			content += "\n"
		}
		content += r.Message
	}
	hasEye, hasStandUp := false, false
	for _, r := range reminders {
		switch r.Type {
		case "eye":
			hasEye = true
		case "standup":
			hasStandUp = true
		}
	}
	var title string
	switch {
	case hasEye && hasStandUp:
		title = i18n.Active.NotificationCombinedTitle
	case hasStandUp:
		title = pickRandom(i18n.Active.NotificationStandUpReminders)
	default:
		title = pickRandom(i18n.Active.NotificationReminders)
	}
	showBalloon(title, content)
}
