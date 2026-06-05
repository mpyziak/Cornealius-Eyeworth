//go:build windows

package systray

import (
	"syscall"
	"time"
	"unsafe"
)

const (
	// hwndMessage is the HWND_MESSAGE pseudo-parent (-3 as uintptr).
	// A window whose parent is HWND_MESSAGE is a "message-only window":
	// it is completely invisible to EnumWindows and all shell-extension
	// enumeration, cannot appear in the taskbar or Alt-Tab list, and
	// receives no broadcast/paint messages — but still receives messages
	// sent directly to it, including Shell_NotifyIcon callback messages.
	hwndMessage uintptr = ^uintptr(2) // (HWND)(LONG_PTR)(-3)
)

var (
	modUser32p   = syscall.NewLazyDLL("user32.dll")
	procFindWinP = modUser32p.NewProc("FindWindowW")
	procSetPar   = modUser32p.NewProc("SetParent")
)

// patchTrayWindowsWhenReady polls until both Win32 windows created by Fyne's
// systray subsystem exist, then reparents each to HWND_MESSAGE.
//
// Root cause of the extra "Quit" menu item in the corporate environment:
//
//   - "SystrayClass" (owned by fyne.io/systray) and "SystrayMonitor" (the
//     hidden GLFW window Fyne creates in its driver) are both top-level windows
//     with WS_OVERLAPPEDWINDOW style (includes WS_CAPTION + WS_SYSMENU).
//   - Corporate endpoint-security shell extensions call EnumWindows to build a
//     list of every top-level window in all processes.  When they find a window
//     whose owning process has a notification-area icon, they inject an extra
//     context-menu item ("Quit" / "Close") in the OS UI language.
//   - Style patching (SetWindowLongPtr / WS_EX_TOOLWINDOW) happens after the
//     window is already visible to the shell and does not un-register it from
//     the extension's internal table.
//
// Fix: SetParent(hwnd, HWND_MESSAGE) converts each window into a message-only
// window.  Message-only windows are invisible to EnumWindows by design; no
// shell extension can enumerate or inject items for them.  The tray icon still
// works because Shell_NotifyIcon holds the HWND directly and delivers callback
// messages to it without any enumeration.
//
// NOTE: NIM_SETVERSION(NOTIFYICON_VERSION_4) is intentionally NOT called.
// fyne.io/systray switches on lParam == WM_RBUTTONUP (0x0205) to show the
// context menu.  Version 4 repacks lParam to LOWORD = event / HIWORD = icon ID,
// which breaks that switch and prevents the menu from ever opening.
func patchTrayWindows() { go patchTrayWindowsWhenReady() }

func patchTrayWindowsWhenReady() {
	const (
		maxAttempts  = 50
		pollInterval = 100 * time.Millisecond
	)

	var systrayHwnd uintptr

	for i := 0; i < maxAttempts; i++ {
		if systrayHwnd == 0 {
			cls, _ := syscall.UTF16PtrFromString("SystrayClass")
			systrayHwnd, _, _ = procFindWinP.Call(uintptr(unsafe.Pointer(cls)), 0)
		}
		if systrayHwnd != 0 {
			break
		}
		time.Sleep(pollInterval)
	}

	if systrayHwnd != 0 {
		procSetPar.Call(systrayHwnd, hwndMessage)
	}
	// SystrayMonitor is intentionally NOT reparented.
	// It is a GLFW-managed window with an OpenGL context; calling SetParent on
	// it cross-thread triggers WM_WINDOWPOSCHANGING/CHANGED/SIZE/MOVE back on
	// the main thread, which can leave GLFW in an inconsistent geometry state
	// and cause a continuous-polling loop (100% CPU) and unbounded GL allocation.
	// SystrayMonitor has no Shell_NotifyIcon entry, so no shell extension will
	// enumerate or inject into it regardless of its parent.
}
