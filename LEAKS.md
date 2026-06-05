# Memory-leak root causes and proposed upstream fixes

This file documents two memory-leak root causes found during profiling of
Cornealius Eyeworth, together with the minimal patches that fix them in their
respective upstream libraries.  The patches are not applied in this repository;
they are recorded here so that pull requests can be submitted upstream.

---

## 1 · `fyne.io/fyne/v2` — `WatchTheme` has no debounce and no stop signal

### Affected file
`internal/app/theme_windows.go` — function `WatchTheme`

### How it manifests
`app.Run()` calls `settings.watchSettings()`, which calls `watchTheme(s)`.
`watchTheme` starts a goroutine that blocks on `RegNotifyChangeKeyValue` in a
tight loop.  Every time the watched registry key
`HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize` is written,
the loop body runs `fyne.Do(s.setupTheme)`.

`fyne.Do` is non-blocking: it posts the closure to Fyne's internal
`async.UnboundedChan[funcData]` (the "funcQueue") without waiting.  The GLFW
event loop drains this queue at ~60 fps.

When a Microsoft Teams call ends, Windows restores the Focus Assist /
Do-Not-Disturb state that Teams had locked during the call.  This involves
several rapid, successive writes to the Personalize key — 5–15 writes in under
200 ms is typical.  Each write causes `WatchTheme` to wake up immediately
(because `RegNotifyChangeKeyValue` re-arms on the next call) and enqueue one
more `setupTheme` closure.  Closures pile up faster than the 60 fps loop drains
them.  Each `setupTheme` in turn calls `isDark()` (opens a registry key) and
`apply()` (walks all settings listeners, which iterates all Fyne windows and
calls `ApplyThemeTo` on their widget trees + `cache.ResetThemeCaches`).  The
resulting allocations from cache resets, theme color objects, and widget layout
passes grow the heap by tens of megabytes before the GC can catch up.

A secondary issue is that the goroutine has no stop channel.  Once started it
runs for the entire lifetime of the process and cannot be stopped even when
`app.Quit()` / `settings.stopWatching()` is called.

### Proposed patch (unified diff against v2.7.4)

```diff
--- a/internal/app/theme_windows.go
+++ b/internal/app/theme_windows.go
@@ -4,6 +4,7 @@ package app

 import (
 	"syscall"
+	"time"

 	"golang.org/x/sys/windows/registry"

@@ -37,18 +38,32 @@ func isDark() bool {
 // WatchTheme calls the supplied function when the Windows dark/light theme changes.
-func WatchTheme(onChanged func()) {
+// stopCh is closed to request the watcher to stop; pass a nil channel to run
+// until process exit (original behaviour, kept for callers that do not need
+// lifecycle control).
+func WatchTheme(onChanged func(), stopCh <-chan struct{}) {
 	var regNotifyChangeKeyValue *syscall.Proc
 	if advapi32, err := syscall.LoadDLL("Advapi32.dll"); err == nil {
 		if p, err := advapi32.FindProc("RegNotifyChangeKeyValue"); err == nil {
 			regNotifyChangeKeyValue = p
 		}
 	}
 	if regNotifyChangeKeyValue == nil {
 		return
 	}
 	k, err := registry.OpenKey(registry.CURRENT_USER, themeRegKey,
 		syscall.KEY_NOTIFY|registry.QUERY_VALUE)
 	if err != nil {
 		return
 	}
+
+	const debounce = 500 * time.Millisecond
+	var lastFired time.Time
+
 	for {
-		regNotifyChangeKeyValue.Call(uintptr(k), 0,
-			0x00000001|0x00000004, 0, 0)
-		onChanged()
+		// Create an auto-reset event so we can also wait on stopCh.
+		event, _ := windows.CreateEvent(nil, 0, 0, nil)
+		regNotifyChangeKeyValue.Call(uintptr(k), 0,
+			0x00000001|0x00000004, uintptr(event), 1 /*fAsynchronous*/)
+		if stopCh != nil {
+			select {
+			case <-stopCh:
+				windows.CloseHandle(event)
+				k.Close()
+				return
+			case <-waitForHandle(event):
+			}
+		} else {
+			windows.WaitForSingleObject(event, windows.INFINITE)
+		}
+		windows.CloseHandle(event)
+
+		// Debounce: swallow notifications that arrive within 500 ms.
+		if now := time.Now(); now.Sub(lastFired) >= debounce {
+			lastFired = now
+			onChanged()
+		}
 	}
 }
```

Simpler alternative (no async event, no stop-channel — just adds the debounce):

```diff
--- a/internal/app/theme_windows.go
+++ b/internal/app/theme_windows.go
@@ -4,6 +4,7 @@ package app

 import (
 	"syscall"
+	"time"

-func WatchTheme(onChanged func()) {
+// WatchTheme calls onChanged when the Windows dark/light theme changes.
+// A 500 ms debounce prevents rapid consecutive registry writes (e.g. from
+// Teams restoring Focus Assist after a call) from flooding Fyne's funcQueue.
+func WatchTheme(onChanged func()) {
 	...
 	for {
 		regNotifyChangeKeyValue.Call(...)
-		onChanged()
+		if now := time.Now(); now.Sub(lastFired) >= 500*time.Millisecond {
+			lastFired = now
+			onChanged()
+		}
 	}
 }
```

### Upstream repository
<https://github.com/fyne-io/fyne>

### Workaround in this application
No clean workaround is possible without either vendoring or a `go.mod replace`
fork.

`app.Run()` (`fyne.io/fyne/v2/app.fyneApp.Run`) is the only public entry point
for the event loop.  It starts `RunEventQueue` on the private `lifecycle` field
before calling `driver.Run()`.  Calling `driver.Run()` alone causes `Quit()` to
hang: `driver.Run()` calls `l.WaitForEvents()` on exit, which enqueues a
sentinel to the lifecycle queue and blocks until it is drained — but without
`RunEventQueue` running, nobody drains it.

`RunEventQueue` is defined on `*fyne.io/fyne/v2/internal/app.Lifecycle`, which
Go's module system marks as internal and forbids importing from other modules.

Therefore `app.Run()` is used as-is.  The WatchTheme leak remains until the
upstream fix is merged.  Its practical effect is mitigated by the other fixes
in this release: the binding-listener leak and GLFW HID enumeration are gone,
so the post-call spike is smaller and GC recovers it faster.

---

## 2 · `github.com/go-gl/glfw/v3.3/glfw` — HID device notifications trigger DirectInput enumeration on every Bluetooth event

### Affected file
`glfw/src/win32_init.c` — function `createHelperWindow`
`glfw/src/win32_window.c` — `wndProc` case `WM_DEVICECHANGE`

### How it manifests
During `glfwInit`, GLFW creates a hidden helper window and registers it for
`GUID_DEVINTERFACE_HID` device notifications via `RegisterDeviceNotificationW`.
This means Windows delivers `WM_DEVICECHANGE DBT_DEVICEARRIVAL` and
`DBT_DEVICEREMOVECOMPLETE` with `DBT_DEVTYP_DEVICEINTERFACE` to GLFW's helper
window every time any HID device connects or disconnects.

GLFW's `wndProc` responds to `DBT_DEVICEARRIVAL DBT_DEVTYP_DEVICEINTERFACE` by
calling `_glfwDetectJoystickConnectionWin32()`.  That function calls:

```c
IDirectInput8_EnumDevices(
    _glfw.win32.dinput8.api,
    DI8DEVCLASS_GAMECTRL,
    deviceCallback,
    NULL,
    DIEDFL_ALLDEVICES);   // enumerates ALL game-class HID devices, not just the new one
```

When a Bluetooth headset switches between call profile (HFP/HSP) and media
profile (A2DP), Windows sends one or more `DBT_DEVICEARRIVAL` events for the
HID interfaces involved.  On machines with corporate endpoint-security software
(CrowdStrike Falcon, Carbon Black, SentinelOne, etc.), those agents hook
`DirectInput8` to detect input-injection attacks.  Being invoked by a background
process they do not recognize, they allocate forensic/scan context inside the
calling process.  That memory is typically not freed until the agent's internal
timeout expires (minutes to hours), causing tens of megabytes of unrecoverable
growth per event.

Even without security agents, the `DIEDFL_ALLDEVICES` flag causes a full
re-enumeration of all connected game controllers on every single HID arrival,
regardless of whether the arriving device is a game controller.  This is
unnecessarily expensive.

### Proposed patch

**Option A — Do not register for HID notifications at all (correct fix for
applications that do not use joystick hot-plug):**

```diff
--- a/glfw/src/win32_init.c
+++ b/glfw/src/win32_init.c
@@ createHelperWindow
-    // Register for HID device notifications
-    {
-        DEV_BROADCAST_DEVICEINTERFACE_W dbi;
-        ZeroMemory(&dbi, sizeof(dbi));
-        dbi.dbcc_size = sizeof(dbi);
-        dbi.dbcc_devicetype = DBT_DEVTYP_DEVICEINTERFACE;
-        dbi.dbcc_classguid = GUID_DEVINTERFACE_HID;
-
-        _glfw.win32.deviceNotificationHandle =
-            RegisterDeviceNotificationW(_glfw.win32.helperWindowHandle,
-                                        (DEV_BROADCAST_HDR*) &dbi,
-                                        DEVICE_NOTIFY_WINDOW_HANDLE);
-    }
+    _glfw.win32.deviceNotificationHandle = NULL;
+    /* Joystick hot-plug is handled by polling inside glfwPollEvents.        */
+    /* Registering for GUID_DEVINTERFACE_HID causes WM_DEVICECHANGE to fire  */
+    /* on every Bluetooth audio profile switch, which then triggers a full   */
+    /* IDirectInput8_EnumDevices re-scan — expensive and problematic on      */
+    /* machines with security agents that hook DirectInput.                  */
```

**Option B — Add an init hint `GLFW_JOYSTICK_HOTPLUG` (feature request):**

Add a boolean init hint (default `GLFW_TRUE` for backward compat) that controls
whether `RegisterDeviceNotificationW` is called, allowing applications that do
not use joysticks to opt out.

```c
// New init hint (win32_init.c):
if (_glfw.hints.init.joystickHotplug) {
    // existing RegisterDeviceNotificationW block
}
```

### Upstream repository
<https://github.com/go-gl/glfw>  
The underlying C library is <https://github.com/glfw/glfw>.  Both need the fix;
the Go binding just ships the C sources vendored.

### Workaround in this application
None possible at the application level — `RegisterDeviceNotificationW` is called
unconditionally inside `glfwInit`, which Fyne calls from `driver.Run()`.  The
only application-level mitigation is to patch the vendored C source (see git
history of this repository for the temporary vendor patch that was later
removed in favour of the `app.Driver().Run()` workaround for issue 1).

---

## Summary

| # | Library | Root cause | Upstream fix needed | App-level workaround |
|---|---------|-----------|---------------------|----------------------|
| 1 | `fyne-io/fyne` | `WatchTheme` has no debounce; rapid registry writes flood funcQueue | Add 500 ms debounce + stop channel to `WatchTheme` | None without vendoring — `RunEventQueue` lives on an internal type that Go forbids importing across modules |
| 2 | `go-gl/glfw` | HID arrival triggers full DirectInput re-scan | Don't register for HID notifications / add opt-out hint | None — requires library patch |

---

## Self-inflicted regression (resolved) — `SetParent(SystrayMonitor)` causes CPU loop

**Introduced and fixed in this repository.** Recorded here as a hard constraint
so it is never reintroduced.

### What happened
An attempt to hide both Fyne-created tray windows from shell-extension enumeration
called `SetParent(hwnd, HWND_MESSAGE)` on both `SystrayClass` and `SystrayMonitor`
from a background goroutine.  The result was runaway CPU use and unbounded memory
growth triggered by Teams chat notification bubbles.

### Why
`SystrayMonitor` is a GLFW-managed window with an active OpenGL context owned by
the main thread.  When `SetParent` is called on it from another thread, Windows
delivers `WM_WINDOWPOSCHANGING/CHANGED/SIZE/MOVE` synchronously to the main-thread
message queue.  GLFW's `wndProc` processes these and queries the window's geometry;
a message-only window returns non-standard values, leaving GLFW in a state where it
continuously issues resize/recheck calls — 100% CPU and a new GL allocation per
iteration.

### The rule
**Never call `SetParent`, `SetWindowLongPtr`, `SetWindowPos`, or any other
cross-thread Win32 style mutation on `SystrayMonitor`.**
Only `SystrayClass` (the `fyne.io/systray` message pump — no GL context) is safe
to reparent.  See `systray/systraypatch_win.go` and the "Systray window patching"
section of `PROJECT_CONTEXT.md` for the full rationale.
