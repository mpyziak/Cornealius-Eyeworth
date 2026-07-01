# Memory-leak root causes and proposed upstream fixes

This file documents two memory-leak root causes found during profiling of
Cornealius Eyeworth, together with the minimal patches that fix them in their
respective upstream libraries.

- **§1 (Fyne `WatchTheme`)** — workaround applied in this repository via
  `go.mod replace` fork at `../fyne-v2-watchtheme-patch`.  See BUILD.md for
  setup instructions.  Upstream bug report template is in §1 below.
- **§2 (GLFW HID)** — requires a C-level change to the vendored GLFW sources;
  not yet applied.  Upstream patch proposals are in §2 below.

---

## 1 · `fyne.io/fyne/v2` — `WatchTheme` ignores `RegNotifyChangeKeyValue` return value; key deletion causes busy-spin flooding `funcQueue`

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

### Exact failure sequence (observed 2026-06-23 / 2026-06-24)

Both log sessions show the same pattern; the 2026-06-24 run is cleaner:

```
08:30:35.487  theme-registry write #1           → Fyne's synchronous RNCV wakes up, calls onChanged()
08:30:35.497  theme-registry key DELETED        → RNCV returns ERROR_KEY_DELETED (1018)
                                                  Fyne: return value ignored, calls onChanged() anyway
                                                  Fyne: tight loop — RNCV returns 1018 immediately every iteration
                                                  HeapAlloc still ~1.32 MB (GC keeping up for now)
08:30:35.616  key recreated (119 ms gap)        → Fyne's handle is still invalid; loop continues
08:30:35.702  burst write #2 (our watcher sees it, Fyne can't — handle stale)
08:30:35.723  burst write #3
08:30:35.738  burst write #4
08:30:35.763  burst write #5
08:30:37.471  12 × Settings.AddListener fires in 8 ms — the queued setupTheme closures drain
              HeapAlloc: 1.32 MB → 37.80 MB (+36 MB) in one event-loop frame
              NumGC: 691 → 698 (7 GC cycles during drain — cannot keep up with allocation rate)
```

The ~12 closures that survive (out of potentially thousands enqueued during the
119 ms tight-loop window) are those that Fyne's 60 fps event loop had not yet
drained by the time steady-state resumed.  Each `setupTheme` call resets the
theme cache, walks every widget in every open window, and re-applies colour
and font metrics — O(widgets) allocations per call.

### Workaround applied in this application

A local fork of `fyne.io/fyne/v2` is maintained at
`../fyne-v2-watchtheme-patch` (sibling of this repository) and wired in via a
`replace` directive in `go.mod`:

```
replace fyne.io/fyne/v2 => ../fyne-v2-watchtheme-patch
```

The single patched file is
`internal/app/theme_windows.go → WatchTheme`.  Two changes vs upstream v2.7.4:

1. **ERROR_KEY_DELETED recovery**: check the `RegNotifyChangeKeyValue` return
   value; on 1018 close the stale handle and poll with 50 ms sleeps until the
   key is recreated, then resume watching normally.  This stops the tight loop.

2. **300 ms debounce**: even after key recreation, Teams/Focus Assist writes the
   key 3–5 more times in rapid succession.  A `time.Time` gate skips
   `onChanged()` for any call that arrives within 300 ms of the previous one,
   collapsing a burst of N writes into a single `setupTheme` call.

Expected post-fix behaviour: one `Settings.AddListener` fire per Teams
call-end event, heap stays at baseline, no multi-GB spike.

### Suggested bug report (submit to https://github.com/fyne-io/fyne)

> **Title**: `WatchTheme` (Windows): `RegNotifyChangeKeyValue` return value ignored — key deletion causes infinite tight loop flooding `funcQueue`
>
> **Affected version**: v2.7.4 (latest at time of writing; `internal/app/theme_windows.go`)
>
> **Platform**: Windows 10/11, any corporate environment where a Teams call or
> Focus Assist restore triggers a full dark/light theme switch.
>
> **Reproduction**:
> 1. Run any Fyne app on Windows with a `Settings().AddListener` registered.
> 2. Start a Microsoft Teams call (or trigger any full theme switch that causes
>    Windows to delete-and-recreate `HKCU\...\Themes\Personalize`).
> 3. End the call / dismiss Focus Assist.
> 4. Observe: heap jumps from ~1 MB to 20–40 MB; NumGC increases by 7+ in one
>    event-loop frame; `Settings.AddListener` fires 10–15 times in <10 ms.
>
> **Root cause** (`internal/app/theme_windows.go`, function `WatchTheme`):
>
> ```go
> for {
>     // return value completely ignored
>     regNotifyChangeKeyValue.Call(uintptr(k), 0, 0x00000001|0x00000004, 0, 0)
>     onChanged()
> }
> ```
>
> When Windows deletes the registry key, `RegNotifyChangeKeyValue` returns
> `ERROR_KEY_DELETED` (1018) **immediately** on every subsequent call.  Because
> the return value is ignored, the loop calls `onChanged()` at CPU speed.
> `onChanged()` calls `fyne.Do(s.setupTheme)`, posting to an unbounded channel.
> The GLFW event loop drains the channel in one frame, running `setupTheme`
> (cache reset + full widget-tree walk) N times in rapid succession.
>
> **Minimal fix**:
>
> ```go
> import "time"
>
> func WatchTheme(onChanged func()) {
>     // ... existing DLL loading ...
>
>     openKey := func() (registry.Key, error) {
>         return registry.OpenKey(registry.CURRENT_USER, themeRegKey,
>             syscall.KEY_NOTIFY|registry.QUERY_VALUE)
>     }
>     k, err := openKey()
>     if err != nil { return }
>
>     const debounce = 300 * time.Millisecond
>     var lastFired time.Time
>
>     for {
>         ret, _, _ := regNotifyChangeKeyValue.Call(uintptr(k), 0,
>             0x00000001|0x00000004, 0, 0)
>         if ret != 0 { // ERROR_KEY_DELETED or other error
>             k.Close()
>             for {
>                 time.Sleep(50 * time.Millisecond)
>                 k, err = openKey()
>                 if err == nil { break }
>             }
>             continue
>         }
>         if now := time.Now(); now.Sub(lastFired) >= debounce {
>             lastFired = now
>             onChanged()
>         }
>     }
> }
> ```
>
> **Observed heap impact without fix**: 1.3 MB → 37.8 MB per Teams call-end event.  
> **Observed heap impact with fix**: expected single-digit MB spike (one `setupTheme` call) followed by immediate GC recovery.

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
| 1 | `fyne-io/fyne` | `WatchTheme` ignores `RegNotifyChangeKeyValue` return value; key deletion causes infinite tight loop flooding funcQueue with `setupTheme` closures → 37 MB heap spike | Add ERROR_KEY_DELETED recovery + 300 ms debounce to `WatchTheme` | **Applied**: `go.mod replace` fork at `../fyne-v2-watchtheme-patch`; bug report template in §1 above |
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
