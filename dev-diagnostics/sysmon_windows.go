//go:build diagnostics && windows && sysmon

package diagnostics

// Watches the two event sources behind the leaks in LEAKS.md
// * writes to the Themes\Personalize registry key
// * HID arrivals on WM_DEVICECHANGE.
// Observational only - nothing here changes application behaviour.

import (
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"fyne.io/fyne/v2"
	"golang.org/x/sys/windows"
)

// Lines setupTheme drains up against the registry writes logged below.
func EnableUIDevDiagnosticsSettingsListener(app fyne.App) {
	app.Settings().AddListener(func(_ fyne.Settings) {
		Event("fyne: Settings.AddListener fired (setupTheme drained from funcQueue)")
	})
}

// -- Win32 constants ----------------------------------------------------------

const (
	regNotifyChangeLastSet = uintptr(0x00000004) // REG_NOTIFY_CHANGE_LAST_SET
	keyNotify              = uint32(0x0010)      // KEY_NOTIFY

	wmDeviceChange          = uint32(0x0219)
	dbtDeviceArrival        = uintptr(0x8000)
	dbtDeviceRemoveComplete = uintptr(0x8004)
	dbtDevNodesChanged      = uintptr(0x0007)
	dbtDevtypDeviceIface    = uint32(5) // DBT_DEVTYP_DEVICEINTERFACE

	deviceNotifyWindowHandle = uintptr(0)
	wmPowerBroadcast         = uint32(0x0218)
	pbtApmSuspend            = uintptr(0x0004)
	pbtApmResumeAuto         = uintptr(0x0012) // system-initiated resume
	pbtApmResumeSuspend      = uintptr(0x0007) // user-triggered resume

	wmDisplayChange  = uint32(0x007E)
	wmDpiChanged     = uint32(0x02E0)
	wmCompacting     = uint32(0x0041)
	wmFontChange     = uint32(0x001D)
	wmSettingChange  = uint32(0x001A)
	wmThemeChanged   = uint32(0x031A)
	wmSysColorChange = uint32(0x0015)

	fileNotifyChangeSize      = uintptr(0x00000008) // FILE_NOTIFY_CHANGE_SIZE
	fileNotifyChangeLastWrite = uintptr(0x00000010) // FILE_NOTIFY_CHANGE_LAST_WRITE
	invalidHandleValue        = ^uintptr(0)         // INVALID_HANDLE_VALUE

	dbtDevtypNet = uint32(4) // DBT_DEVTYP_NET: legacy VPN / dial-up adapters

	wmWtsSessionChange   = uint32(0x02B1)
	wtsSessionLock       = uintptr(0x7)
	wtsSessionUnlock     = uintptr(0x8)
	wtsRemoteConnect     = uintptr(0x3)
	wtsRemoteDisconnect  = uintptr(0x4)
	notifyForThisSession = uint32(0)
	waitObject0          = uintptr(0x00000000)
	waitInfinite         = uint32(0xFFFFFFFF)

	hwndMessageSM = ^uintptr(2) // HWND_MESSAGE = (HWND)-3
	wmQuit        = uint32(0x0012)
)

// GUID_DEVINTERFACE_HID = {4D1E55B2-F16F-11CF-88CB-001111000030}
// Stored in the mixed-endian layout Windows uses for GUIDs.
var guidDevIfaceHID = [16]byte{
	0xB2, 0x55, 0x1E, 0x4D, // Data1 LE: 4D1E55B2
	0x6F, 0xF1, // Data2 LE: F16F
	0xCF, 0x11, // Data3 LE: 11CF
	0x88, 0xCB, // Data4[0..1]
	0x00, 0x11, 0x11, 0x00, 0x00, 0x30, // Data4[2..7]
}

// GUID_DEVINTERFACE_NET = {CAC88484-7515-4C03-82E6-71A87ABAC361}
// Fired for VPN adapter arrival/removal and ethernet plug/unplug.
var guidDevIfaceNet = [16]byte{
	0x84, 0x84, 0xC8, 0xCA, // Data1 LE: CAC88484
	0x15, 0x75, // Data2 LE: 7515
	0x03, 0x4C, // Data3 LE: 4C03
	0x82, 0xE6, // Data4[0..1]
	0x71, 0xA8, 0x7A, 0xBA, 0xC3, 0x61, // Data4[2..7]
}

// -- Win32 structs ------------------------------------------------------------

type wndClassExW struct {
	Size, Style uint32
	WndProc     uintptr
	ClsExtra    int32
	WndExtra    int32
	Instance    uintptr
	Icon        uintptr
	Cursor      uintptr
	Background  uintptr
	MenuName    *uint16
	ClassName   *uint16
	IconSm      uintptr
}

type devBroadcastHdr struct {
	Size, DeviceType, Reserved uint32
}

type devBroadcastDeviceIface struct {
	Size, DeviceType, Reserved uint32
	ClassGuid                  [16]byte
	Name                       [2]uint16
}

// msg mirrors MSG (variable tail on Win10+ but the first 6 fields are fixed).
type msgW struct {
	Hwnd    uintptr
	Message uint32
	_       uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	PtX     int32
	PtY     int32
}

// -- Lazy Win32 procs ---------------------------------------------------------

var (
	modAdvapi32 = syscall.NewLazyDLL("advapi32.dll")
	modKernel32 = syscall.NewLazyDLL("kernel32.dll")
	modUser32   = syscall.NewLazyDLL("user32.dll")

	procRegNotify = modAdvapi32.NewProc("RegNotifyChangeKeyValue")

	procCreateEvent  = modKernel32.NewProc("CreateEventW")
	procCloseHandle  = modKernel32.NewProc("CloseHandle")
	procWaitForMulti = modKernel32.NewProc("WaitForMultipleObjects")
	procSetEvent     = modKernel32.NewProc("SetEvent")
	procGetModule    = modKernel32.NewProc("GetModuleHandleW")
	procGetThreadID  = modKernel32.NewProc("GetCurrentThreadId")
	procPostThread   = modUser32.NewProc("PostThreadMessageW")

	procRegisterClassEx = modUser32.NewProc("RegisterClassExW")
	procCreateWindowEx  = modUser32.NewProc("CreateWindowExW")
	procGetMsg          = modUser32.NewProc("GetMessageW")
	procTranslateMsg    = modUser32.NewProc("TranslateMessage")
	procDispatchMsg     = modUser32.NewProc("DispatchMessageW")
	procDefWndProc      = modUser32.NewProc("DefWindowProcW")
	procRegDevNotif     = modUser32.NewProc("RegisterDeviceNotificationW")

	procFindFirstChangeNotif = modKernel32.NewProc("FindFirstChangeNotificationW")
	procFindNextChangeNotif  = modKernel32.NewProc("FindNextChangeNotification")

	modWtsapi32            = syscall.NewLazyDLL("wtsapi32.dll")
	procWTSRegisterSession = modWtsapi32.NewProc("WTSRegisterSessionNotification")
)

// -- Stop handle --------------------------------------------------------------

var sysmonStop syscall.Handle // manual-reset event, set by stopSysmon()

// startSysmon creates the shared stop event and launches all four monitor goroutines.
// Called from Init() after the log file is open.
func startSysmon() {
	h, _, _ := procCreateEvent.Call(0, 1 /*manual-reset*/, 0, 0)
	if h == 0 {
		Warn("sysmon: CreateEvent failed; OS event monitoring disabled")
		return
	}
	sysmonStop = syscall.Handle(h)

	go monitorThemeRegistry()
	go monitorHIDDevices()
	go monitorGPRegistry()
	go monitorNotifications()
}

// -- 3. Group Policy registry watcher -----------------------------------------

// Written at the end of every policy refresh
// A write with no WM_SETTINGCHANGE "Policy" behind it is a scheduled task
// rather than an interactive gpupdate.
var gpWatchKeys = []string{
	`Software\Microsoft\Windows\CurrentVersion\Group Policy\State`,
	`Software\Microsoft\Windows\CurrentVersion\Group Policy\History`,
}

func monitorGPRegistry() {
	type keyEntry struct {
		path string
		h    syscall.Handle
	}

	var keys []keyEntry
	for _, path := range gpWatchKeys {
		keyPtr, _ := syscall.UTF16PtrFromString(path)
		var hKey syscall.Handle
		if err := syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_USER, keyPtr, 0, keyNotify, &hKey); err != nil {
			continue // key may not exist on all machines
		}
		keys = append(keys, keyEntry{path, hKey})
		Info("sysmon: watching HKCU\\%s for GPUPDATE activity", path)
	}
	if len(keys) == 0 {
		Info("sysmon: GP registry keys absent; GPUPDATE monitoring disabled")
		return
	}
	defer func() {
		for _, k := range keys {
			syscall.RegCloseKey(k.h)
		}
	}()

	for {
		// One auto-reset event per key, plus the shared stop event.
		notifyEvts := make([]uintptr, len(keys))
		handleSlice := make([]uintptr, 0, len(keys)+1)
		for i, k := range keys {
			evt, _, _ := procCreateEvent.Call(0, 0 /*auto-reset*/, 0, 0)
			if evt == 0 {
				Warn("sysmon: CreateEvent failed in GP registry watcher; stopping")
				return
			}
			// bWatchSubtree=true so child-key writes are also caught.
			procRegNotify.Call(uintptr(k.h), 1, regNotifyChangeLastSet, evt, 1)
			notifyEvts[i] = evt
			handleSlice = append(handleSlice, evt)
		}
		handleSlice = append(handleSlice, uintptr(sysmonStop))

		stopIdx := uintptr(len(handleSlice) - 1)
		result, _, _ := procWaitForMulti.Call(
			uintptr(len(handleSlice)),
			uintptr(unsafe.Pointer(&handleSlice[0])),
			0,
			uintptr(waitInfinite),
		)
		for _, evt := range notifyEvts {
			procCloseHandle.Call(evt)
		}

		if result == waitObject0+stopIdx {
			return
		}
		if idx := int(result - waitObject0); idx >= 0 && idx < len(keys) {
			Event("system: GPUPDATE-related registry write to HKCU\\%s; likely GPUPDATE cycle (~90 min) or logon policy refresh; check for WM_SETTINGCHANGE \"Policy\" in next few seconds", keys[idx].path)
		}
	}
}

// stopSysmon signals all four monitor goroutines to exit.
// Called from Close() before the log file is closed.
func stopSysmon() {
	if sysmonStop != 0 {
		procSetEvent.Call(uintptr(sysmonStop))
	}
}

// -- 1. Theme-registry watcher ------------------------------------------------

const themeRegKey = `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`

// ERROR_KEY_DELETED. Windows deletes and recreates Themes\Personalize
// during a full theme switch, so reopen the key rather than exiting
// This is the same condition as handled by Fyne patch - see LEAKS.md.
const errKeyDeleted = uintptr(1018)

func openThemeKey() (syscall.Handle, error) {
	keyPtr, _ := syscall.UTF16PtrFromString(themeRegKey)
	var hKey syscall.Handle
	err := syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_USER, keyPtr, 0, keyNotify, &hKey)
	return hKey, err
}

func monitorThemeRegistry() {
	hKey, err := openThemeKey()
	if err != nil {
		Warn("sysmon: cannot open theme key: %v; registry monitoring disabled", err)
		return
	}

	Info("sysmon: watching HKCU\\%s for WatchTheme triggers", themeRegKey)

	var burstCount int
	var lastWrite time.Time

	for {
		// Per-iteration auto-reset event for the async notification.
		notifyEvt, _, _ := procCreateEvent.Call(0, 0 /*auto-reset*/, 0, 0)
		if notifyEvt == 0 {
			syscall.RegCloseKey(hKey)
			Warn("sysmon: CreateEvent failed in registry watcher; stopping")
			return
		}

		ret, _, _ := procRegNotify.Call(
			uintptr(hKey),
			0,
			regNotifyChangeLastSet,
			notifyEvt,
			1, // fAsynchronous
		)
		if ret != 0 {
			procCloseHandle.Call(notifyEvt)
			syscall.RegCloseKey(hKey)
			if ret == errKeyDeleted {
				// Deleted during a theme switch; wait for it to come back.
				Event("system: theme-registry key deleted; full theme switch in progress (dark/light, Teams call restore); Fyne WatchTheme handle also invalidated; watching for key recreation")
				for {
					// Check stop event without blocking.
					stopHandles := [1]uintptr{uintptr(sysmonStop)}
					r, _, _ := procWaitForMulti.Call(1, uintptr(unsafe.Pointer(&stopHandles[0])), 0, 100 /*ms*/)
					if r == waitObject0 {
						return // stop signalled
					}
					hKey, err = openThemeKey()
					if err == nil {
						Event("system: theme-registry key recreated; resuming watch")
						break
					}
				}
				continue
			}
			Warn("sysmon: RegNotifyChangeKeyValue failed (ret=%d); registry monitoring stopped", ret)
			return
		}

		handles := [2]uintptr{notifyEvt, uintptr(sysmonStop)}
		result, _, _ := procWaitForMulti.Call(
			2,
			uintptr(unsafe.Pointer(&handles[0])),
			0, // bWaitAll=false
			uintptr(waitInfinite),
		)
		procCloseHandle.Call(notifyEvt)

		if result == waitObject0+1 { // stop event signalled
			syscall.RegCloseKey(hKey)
			return
		}
		if result != waitObject0 {
			syscall.RegCloseKey(hKey)
			Warn("sysmon: WaitForMultipleObjects unexpected result=%d", result)
			return
		}

		// Registry key was written.
		now := time.Now()
		if now.Sub(lastWrite) < 600*time.Millisecond {
			burstCount++
			Event("system: theme-registry write; burst #%d within 600 ms; Fyne WatchTheme queues another setupTheme onto funcQueue (cache reset + layout pass per call)", burstCount)
		} else {
			burstCount = 1
			Event("system: theme-registry write; Fyne WatchTheme will queue setupTheme onto funcQueue")
		}
		lastWrite = now
	}
}

// -- 2. HID device-notification watcher ---------------------------------------

// hidWndProc is the window procedure for our HID-monitoring window.
// Package-level so its address is stable for the lifetime of the process.
var hidWndProc = syscall.NewCallback(func(hwnd, msg, wp, lp uintptr) uintptr {
	switch uint32(msg) {
	case wmDeviceChange:
		// DBT_DEVNODES_CHANGED has no lParam.
		if wp == dbtDevNodesChanged {
			Event("system: WM_DEVICECHANGE DBT_DEVNODES_CHANGED; device tree changed")
		} else if lp != 0 {
			// lParam is a WndProc-supplied pointer to a DEV_BROADCAST_HDR (or a
			// longer struct sharing its layout); converting it is inherent to
			// Win32 callbacks like this one, not an accidental unsafe use. The
			// go vet warning below is expected.
			hdr := (*devBroadcastHdr)(unsafe.Pointer(lp))
			switch hdr.DeviceType {
			case dbtDevtypDeviceIface:
				dbi := (*devBroadcastDeviceIface)(unsafe.Pointer(lp))
				isHID := dbi.ClassGuid == guidDevIfaceHID
				isNet := dbi.ClassGuid == guidDevIfaceNet
				switch wp {
				case dbtDeviceArrival:
					if isHID {
						Event("system: WM_DEVICECHANGE DBT_DEVICEARRIVAL (HID); GLFW will call IDirectInput8_EnumDevices; corporate AV agents may allocate scan context here")
					} else if isNet {
						Event("system: WM_DEVICECHANGE DBT_DEVICEARRIVAL (network adapter); VPN connect or ethernet plug-in; IP stack rebuild may follow")
					} else {
						Event("system: WM_DEVICECHANGE DBT_DEVICEARRIVAL (other device interface)")
					}
				case dbtDeviceRemoveComplete:
					if isHID {
						Event("system: WM_DEVICECHANGE DBT_DEVICEREMOVECOMPLETE (HID)")
					} else if isNet {
						Event("system: WM_DEVICECHANGE DBT_DEVICEREMOVECOMPLETE (network adapter); VPN disconnect or ethernet unplug")
					} else {
						Event("system: WM_DEVICECHANGE DBT_DEVICEREMOVECOMPLETE (other device interface)")
					}
				}
			case dbtDevtypNet:
				// Legacy mechanism used by older VPN clients and dial-up adapters.
				switch wp {
				case dbtDeviceArrival:
					Event("system: WM_DEVICECHANGE DBT_DEVICEARRIVAL (DBT_DEVTYP_NET); legacy VPN/dial-up adapter connected")
				case dbtDeviceRemoveComplete:
					Event("system: WM_DEVICECHANGE DBT_DEVICEREMOVECOMPLETE (DBT_DEVTYP_NET); legacy VPN/dial-up adapter disconnected")
				}
			}
		}

	case wmPowerBroadcast:
		switch wp {
		case pbtApmSuspend:
			Event("system: WM_POWERBROADCAST PBT_APMSUSPEND; system going to sleep/hibernate")
		case pbtApmResumeAuto:
			Event("system: WM_POWERBROADCAST PBT_APMRESUMEAUTOMATIC; system resumed (auto/timer)")
		case pbtApmResumeSuspend:
			Event("system: WM_POWERBROADCAST PBT_APMRESUMESUSPEND; system resumed (user action)")
		}

	case wmDisplayChange:
		// Resolution/depth change, screen-share start/stop, RDP attach, docking.
		depth := wp & 0xFFFF
		width := lp & 0xFFFF
		height := (lp >> 16) & 0xFFFF
		Event("system: WM_DISPLAYCHANGE depth=%d resolution=%dx%d; Fyne reloadScale+SetDirty on all canvases; GL framebuffers reallocated", depth, width, height)

	case wmDpiChanged:
		// Docking, RDP client resize, VDI resize, display-scaling GPO.
		dpiX := wp & 0xFFFF
		dpiY := (wp >> 16) & 0xFFFF
		Event("system: WM_DPICHANGED new DPI x=%d y=%d; Fyne RescaleContext+GL framebuffer realloc on all visible windows", dpiX, dpiY)

	case wmCompacting:
		// If this shows up, the leak is already system-wide.
		Event("system: WM_COMPACTING; system physical memory critically low; Windows compacting working sets")

	case wmFontChange:
		Event("system: WM_FONTCHANGE; installed font set changed (GPO font deployment?); Fyne glyph atlas rebuilt on next render")

	case wmSettingChange:
		// ImmersiveColorSet = Teams/Focus Assist toggling dark mode.
		// Policy            = GPUPDATE finished. Environment = some VPN clients.
		if lp != 0 {
			// lParam is a WndProc-supplied pointer to a NUL-terminated string;
			// converting it is inherent to Win32 callbacks (the go vet warning
			// below is expected). UTF16PtrToString walks to the real NUL
			// instead of assuming a fixed-size window that may over- or
			// under-read the sender's buffer.
			param := windows.UTF16PtrToString((*uint16)(unsafe.Pointer(lp)))
			if param == "Policy" {
				Event("system: WM_SETTINGCHANGE param=%q; GROUP POLICY REFRESH; correlate with GP-registry writes to confirm GPUPDATE cycle", param)
			} else {
				Event("system: WM_SETTINGCHANGE param=%q; may trigger Fyne settings listener cascade", param)
			}
		} else {
			Event("system: WM_SETTINGCHANGE (no param); may trigger Fyne settings listener cascade")
		}

	case wmThemeChanged:
		Event("system: WM_THEMECHANGED; Fyne theme cache invalidated on all windows")

	case wmSysColorChange:
		// High-contrast and accessibility toggles.
		Event("system: WM_SYSCOLORCHANGE; may trigger Fyne theme/settings listener cascade")

	case wmWtsSessionChange:
		switch wp {
		case wtsSessionLock:
			Event("system: WM_WTSSESSION_CHANGE WTS_SESSION_LOCK; workstation locked")
		case wtsSessionUnlock:
			Event("system: WM_WTSSESSION_CHANGE WTS_SESSION_UNLOCK; workstation unlocked")
		case wtsRemoteConnect:
			Event("system: WM_WTSSESSION_CHANGE WTS_REMOTE_CONNECT; RDP session attached; expect WM_DISPLAYCHANGE + DPI cascade")
		case wtsRemoteDisconnect:
			Event("system: WM_WTSSESSION_CHANGE WTS_REMOTE_DISCONNECT; RDP session detached")
		}
	}

	ret, _, _ := procDefWndProc.Call(hwnd, msg, wp, lp)
	return ret
})

func monitorHIDDevices() {
	// This goroutine owns its Win32 window; it must stay on the same OS thread.
	runtime.LockOSThread()
	// Intentionally not unlocking: the goroutine is dedicated to this loop.

	className, _ := syscall.UTF16PtrFromString("CornealiusHIDWatcher")
	hInst, _, _ := procGetModule.Call(0)

	wc := wndClassExW{
		WndProc:   hidWndProc,
		Instance:  hInst,
		ClassName: className,
	}
	wc.Size = uint32(unsafe.Sizeof(wc))

	if atom, _, _ := procRegisterClassEx.Call(uintptr(unsafe.Pointer(&wc))); atom == 0 {
		Warn("sysmon: RegisterClassEx failed; HID monitoring disabled")
		return
	}

	// Message-only window (HWND_MESSAGE parent), invisible to EnumWindows.
	hwnd, _, _ := procCreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		0, 0,
		0, 0, 0, 0,
		hwndMessageSM, 0, hInst, 0,
	)
	if hwnd == 0 {
		Warn("sysmon: CreateWindowEx failed; HID monitoring disabled")
		return
	}

	dbi := devBroadcastDeviceIface{
		DeviceType: dbtDevtypDeviceIface,
		ClassGuid:  guidDevIfaceHID,
	}
	dbi.Size = uint32(unsafe.Sizeof(dbi))

	hNotif, _, _ := procRegDevNotif.Call(hwnd, uintptr(unsafe.Pointer(&dbi)), deviceNotifyWindowHandle)
	if hNotif == 0 {
		Warn("sysmon: RegisterDeviceNotification(HID) failed; HID events may not be captured")
	} else {
		Info("sysmon: watching GUID_DEVINTERFACE_HID for WM_DEVICECHANGE triggers")
	}

	// Also register for network adapter events (VPN connect/disconnect, ethernet).
	dbiNet := devBroadcastDeviceIface{
		DeviceType: dbtDevtypDeviceIface,
		ClassGuid:  guidDevIfaceNet,
	}
	dbiNet.Size = uint32(unsafe.Sizeof(dbiNet))
	if hNotifNet, _, _ := procRegDevNotif.Call(hwnd, uintptr(unsafe.Pointer(&dbiNet)), deviceNotifyWindowHandle); hNotifNet == 0 {
		Warn("sysmon: RegisterDeviceNotification(NET) failed; VPN/ethernet events may not be captured")
	} else {
		Info("sysmon: watching GUID_DEVINTERFACE_NET for VPN/ethernet WM_DEVICECHANGE triggers")
	}

	// Register for session-change notifications (lock/unlock).
	// WTSRegisterSessionNotification requires a real HWND, not HWND_MESSAGE on
	// some older Windows versions; if it fails we still get power events.
	if ret, _, _ := procWTSRegisterSession.Call(hwnd, uintptr(notifyForThisSession)); ret == 0 {
		Warn("sysmon: WTSRegisterSessionNotification failed; session lock/unlock events will not be logged")
	} else {
		Info("sysmon: watching WM_WTSSESSION_CHANGE for session lock/unlock")
	}

	Info("sysmon: watching WM_POWERBROADCAST for sleep/resume (delivered automatically by Windows)")

	// Post WM_QUIT to this thread when the stop event fires.
	threadID, _, _ := procGetThreadID.Call()
	go func() {
		handles := [1]uintptr{uintptr(sysmonStop)}
		procWaitForMulti.Call(1, uintptr(unsafe.Pointer(&handles[0])), 0, uintptr(waitInfinite))
		procPostThread.Call(threadID, uintptr(wmQuit), 0, 0)
	}()

	Info("sysmon: HID/power/session monitor running")

	var msg msgW
	for {
		ret, _, _ := procGetMsg.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		switch int32(ret) {
		case -1:
			Warn("sysmon: GetMessage error in HID watcher")
			return
		case 0:
			return // WM_QUIT
		default:
			procTranslateMsg.Call(uintptr(unsafe.Pointer(&msg)))
			procDispatchMsg.Call(uintptr(unsafe.Pointer(&msg)))
		}
	}
}

// -- 4. Windows notification database watcher ---------------------------------

// wpndatabase.db is touched on every toast from any app
// this dates notification bursts against memory growth
// Ending a Teams call writes both this and the theme key
func monitorNotifications() {
	localAppData, _ := syscall.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		Warn("sysmon: LOCALAPPDATA not set; notification monitoring disabled")
		return
	}

	notifDir := localAppData + `\Microsoft\Windows\Notifications`
	dirPtr, _ := syscall.UTF16PtrFromString(notifDir)

	hChange, _, _ := procFindFirstChangeNotif.Call(
		uintptr(unsafe.Pointer(dirPtr)),
		0, // bWatchSubtree = false
		fileNotifyChangeLastWrite|fileNotifyChangeSize,
	)
	if hChange == 0 || hChange == invalidHandleValue {
		Warn("sysmon: FindFirstChangeNotification for WPN directory failed; notification monitoring disabled")
		return
	}
	defer procCloseHandle.Call(hChange)

	Info("sysmon: watching %s for WPN notification database writes", notifDir)

	for {
		handles := [2]uintptr{hChange, uintptr(sysmonStop)}
		result, _, _ := procWaitForMulti.Call(
			2,
			uintptr(unsafe.Pointer(&handles[0])),
			0,
			uintptr(waitInfinite),
		)

		if result == waitObject0+1 { // stop event
			return
		}
		if result != waitObject0 {
			Warn("sysmon: WaitForMultipleObjects unexpected result=%d in notification watcher", result)
			return
		}

		Event("system: WPN notification delivered; Windows notification platform wrote to notification database (toast from Teams, Outlook, Defender, calendar, etc.)")

		// Rearm the change handle for the next write.
		if ret, _, _ := procFindNextChangeNotif.Call(hChange); ret == 0 {
			Warn("sysmon: FindNextChangeNotification failed; notification monitoring stopped")
			return
		}
	}
}
