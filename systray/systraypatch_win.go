//go:build windows

package systray

import (
	"syscall"
	"time"
	"unsafe"
)

// hwndMessage is the HWND_MESSAGE pseudo-parent (-3 as uintptr). A window
// whose parent is HWND_MESSAGE is a "message-only window": it never appears
// on screen or in the taskbar, but still receives its messages.
const (
	hwndMessage uintptr = ^uintptr(2) // HWND_MESSAGE, (HWND)(LONG_PTR)(-3)
)

var (
	modUser32p   = syscall.NewLazyDLL("user32.dll")
	procFindWinP = modUser32p.NewProc("FindWindowW")
	procSetPar   = modUser32p.NewProc("SetParent")
)

// Hides SystrayClass from EnumWindows, which is how endpoint-security shell
// extensions inject their own "Quit" into the tray menu. Message-only
// windows are not enumerable; Shell_NotifyIcon holds the HWND directly, so the
// tray still works. Style flags (WS_EX_TOOLWINDOW) are too late by then.
//
// Do not add NIM_SETVERSION(NOTIFYICON_VERSION_4)
// fyne.io/systray switches on lParam == WM_RBUTTONUP
// v4 repacks lParam, so the menu stops opening.
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
	// Never reparent SystrayMonitor. It is a GLFW-managed window with an OpenGL context.
	// calling SetParent cross-thread triggers WM_WINDOWPOSCHANGING/CHANGED/SIZE/MOVE back on
	// the main thread, which can leave GLFW in an inconsistent geometry state
	// and cause a continuous-polling loop (100% CPU) and unbounded GL allocation.
	// It has no Shell_NotifyIcon entry, so nothing can inject into it anyway.
}
