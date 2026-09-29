# Cornealius-Eyeworth - Project Context

## Overview

Eye-rest and stand-up reminder app. Lives in the system tray; fires balloon
notifications on a CRON schedule. Cross-platform source; primary target is Windows.
This is a **Go / Fyne v2** rewrite of an older C# WinForms version. The C# version
no longer exists in this repo.

---

## Quick facts

| Key | Value |
|-----|-------|
| Language | Go 1.22 |
| Module | `github.com/mpyziak/cornealius-eyeworth` |
| UI framework | Fyne v2 (`fyne.io/fyne/v2 v2.7.4`) |
| Systray | `fyne.io/systray v1.12.1` |
| Scheduler | `github.com/robfig/cron/v3` (Quartz-style 6-field, seconds first) |
| Config file | `config.json` next to the executable |
| Entry point | `main.go` → `ui.Run()` |

---

## Package map

```
main.go              Entry point: loads config, sets locale, starts Fyne app + ui.Run()
assets/              Embeds Logo.png (Fyne resource) and Logo.ico (systray + window icon)
config/              Config struct (CronExpression, StandUpCronExpression, Language) + JSON repo
i18n/                Strings struct, locales registry (Enabled flag), read-only via Locales(), SetLanguage(), Active global
  en/de/pl/fr/el.go  One table per language: English / Deutsch / Polski / Français / Ελληνικά
notifications/       Platform-split notification senders
  notifier.go        !windows: wraps fyne.App.SendNotification()
  notifier_windows.go windows: Shell_NotifyIcon balloon tips via direct Win32 syscalls
  flusher.go         AppFlusher - bridges scheduling.Flusher → SendReminders()
parsing/             ParseCron(), ParseMinutes(), TryExtractSimpleMinutes()
scheduling/
  scheduler.go       Scheduler (hot-restartable robfig/cron wrapper), ApplySchedule()
  buffer.go          Buffer - coalesces rapid firings into one notification (2 s window)
  describer.go       Describe() - CRON → human-readable text
  nexttrigger.go     NextTrigger() - next wall-clock fire time
systray/             Manager: tray icon + on-demand status window lifecycle
ui/
  ui.go              Run() - wires everything together; anchor window + window factory
  schedule.go        ShowScheduleDialog()
  language.go        ShowLanguageDialog()
  about.go           ShowAboutDialog()
  help.go            ShowHelpDialog()
tools/i18n-export/   All locales → CSV for Google Sheets (make i18n-export)
winres/winres.json   Windows resource manifest (icon, version info) for go-winres
rsrc_windows_amd64.syso  Compiled Windows resources (auto-linked by go build)
```

---

## Architecture highlights

### Window lifecycle (important for memory)
- A silent **anchor window** (`app.NewWindow("")`, never shown) is created first to
  keep Fyne's event loop alive. It is never `Show()`'d so Windows never registers its
  HWND as a visible top-level window.
- The status window is **on-demand**: `systray.Manager` holds a `windowFactory func()
  fyne.Window`. Each tray click calls the factory and shows a fresh window;
  closing it calls `win.Close()` (not `win.Hide()`), which destroys the HWND and
  frees the Fyne/GLFW OpenGL context entirely.
- Rationale: `win.Hide()` keeps the HWND registered with the Windows window manager.
  On a locked/screen-shared session Windows routes display events (`WM_DPICHANGED`,
  etc.) to all registered HWNDs. Fyne handles these by re-allocating GL framebuffers
  and texture mipmaps (CGO/native memory, invisible to Go GC), causing multi-GB leaks
  over hours. `Close()` removes the HWND and prevents this.

### Systray window patching (Windows - `systray/systraypatch_win.go`)
On some machines, endpoint-security shell extensions call `EnumWindows` and
inject an extra "Quit" context-menu item into any process that has a notification-
area icon and an enumerable top-level HWND.  Two HWNDs are created by Fyne's
systray subsystem:

- **`SystrayClass`** - plain Win32 message window owned by `fyne.io/systray`.
  `SetParent(hwnd, HWND_MESSAGE)` is safe here: the window has no GL context,
  its message loop (`GetMessage` in a background goroutine) continues to work
  after reparenting, and `Shell_NotifyIcon` holds the HWND directly so tray
  callbacks are unaffected.  **This window IS reparented.**

- **`SystrayMonitor`** - hidden GLFW window created by Fyne's driver as a quit-
  detection hook.  It has an OpenGL context managed by GLFW on the main thread.
  **This window must NEVER be reparented or have its Win32 style mutated
  cross-thread.**  Calling `SetParent` on it from a background goroutine causes
  Windows to deliver `WM_WINDOWPOSCHANGING/CHANGED/SIZE/MOVE` synchronously to
  the main thread.  GLFW's wndProc processes these and queries the window's
  geometry; a message-only window returns non-standard values, leaving GLFW in
  an inconsistent state that produces a continuous polling/recorrection loop
  (100% CPU) and unbounded GL allocation.  `SystrayMonitor` has no
  `Shell_NotifyIcon` entry so no shell extension can target it regardless of
  its parent.

Fyne also auto-injects a second Quit into every systray menu it builds
(`addMissingQuitForMenu` in `driver/glfw/driver_desktop.go`) unless the last
menu item already has `IsQuit = true`.  The quit item in `doSetup()` must
always set `quitItem.IsQuit = true` to prevent this.

### Notifications (Windows)
- `fyne.App.SendNotification()` shells out to PowerShell and is blocked on
  corporate machines. Windows implementation uses `Shell_NotifyIconW` (balloon
  tip) via direct Win32 syscalls instead.
- The systray window (`SystrayClass` HWND, id=100) must exist before calling
  `Shell_NotifyIconW`. `SendStartup` is therefore called from inside the
  `systray.Run` `onReady` callback, not before it.

### Scheduling
- Two independent `Scheduler` instances (eye, stand-up).
- `scheduling.Buffer` coalesces events within a 2-second window, then calls
  `AppFlusher.Flush()` → `SendReminders()`. Prevents duplicate popups when
  both schedules fire in the same second.
- `ApplySchedule()` centralises start/stop logic; used from both `ui.Run()` and
  `ShowScheduleDialog` save handler.

---

## Build

Recommended (matches fyne package tooling, ~25 MB, proper icon on all machines):
```powershell
go run fyne.io/tools/cmd/fyne@latest package -release -app-id "github.com/mpyziak/cornealius-eyeworth" -os windows -icon assets/Logo.png -name "Cornealius Eyeworth"
```

Plain go build (dev/quick):
```powershell
go build -ldflags "-s -w -H windowsgui" -o "Cornealius Eyeworth.exe" .
```

Windows resources (icon, manifest) are embedded via `rsrc_windows_amd64.syso`
generated by `go-winres` (`go generate` in root, or `make winres`).

---

## Memory leak - investigation log

### Symptoms
- Observed exclusively on a **constrained corporate Windows PC**.
- App starts at ~60 MB, drops to 10–15 MB after GC. Stable for many hours.
- Then bloats to 2–6+ GB and keeps growing until killed.
- Reproduced without Teams calls and during a locked-screen session, ruling
  out screen-share hooks as the sole trigger.

### What was ruled out
- The `scheduling.Buffer` and `Scheduler` have no unbounded growth.
- `AppFlusher` and `SendReminders` allocate nothing long-lived.
- No goroutine leak found in production paths.
- `debug.FreeOSMemory()` on a 5-minute ticker was tried but removed (not
  appropriate for release builds; doesn't fix the root cause).

### Confirmed root causes (resolved 2026-06-23/24)

The initial hypothesis (CGO/native memory, DWM/DPI events, GL framebuffer
churn) was ruled out once the `applog` + `sysmon` diagnostics were in place.
The actual causes are Go-heap, not CGO, and are both triggered by
Windows environment events:

**Root cause 1 - `fyne.io/fyne/v2` `WatchTheme` busy-spin on key deletion**  
When a Teams call ends (or Focus Assist restores), Windows briefly *deletes and
recreates* `HKCU\...\Themes\Personalize`.  Fyne's `WatchTheme` calls
`RegNotifyChangeKeyValue` synchronously and ignores its return value.  While
the key is deleted, RNCV returns `ERROR_KEY_DELETED` (1018) immediately every
time, spinning the loop at CPU speed and posting a `setupTheme` closure to
Fyne's unbounded `funcQueue` on every iteration.  In the 119 ms deletion
window, thousands of closures can be enqueued; they all drain in one event-loop
frame, each running a full theme-cache reset + widget-tree walk.  Observed:
heap 1.3 MB → 37.8 MB, 7 GC cycles, in under 10 ms. Left unattended,
app can within minutes bloat in memory to levels causing other apps crashing.

**Root cause 2 - GLFW HID device enumeration via DirectInput on Bluetooth events**  
GLFW registers for `GUID_DEVINTERFACE_HID` notifications and on each arrival
calls `IDirectInput8_EnumDevices(DIEDFL_ALLDEVICES)`.  A Bluetooth headset
switching HFP↔A2DP generates 3–8 such arrivals per switch.  Corporate
endpoint-security agents (CrowdStrike, SentinelOne) hook DirectInput and
allocate forensic scan context per invocation; that memory is not freed until
the agent's internal timeout (minutes to hours).

### Applied mitigations

| Cause | Fix | Status |
|-------|-----|--------|
| `WatchTheme` busy-spin | `go.mod replace` fork at `../fyne-v2-watchtheme-patch`; adds RNCV return-value check + 300 ms debounce | **Applied** |
| On-demand window (eliminates long-lived HWND receiving `WM_DPICHANGED` etc.) | Window factory in `systray.Manager` - `Close()` on hide instead of `Hide()` | **Applied** |
| GLFW HID DirectInput re-scan | Requires C-level GLFW patch; no app-level workaround possible | **Pending upstream** |

Full root-cause analysis, exact log sequences, and upstream bug report templates
are in `LEAKS.md`.
