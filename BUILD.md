# Cornealius Eyeworth — Go + Fyne

Cross-platform app for eye-rest reminders. Builds from the repository root and runs on Windows, Linux, and macOS.

---

## Prerequisites

| Tool | Minimum version | Purpose | Install |
|------|----------------|---------|---------|
| **Go** | 1.22 | Build toolchain | https://go.dev/dl/ |
| **C compiler** | any | CGO — required by Fyne's OpenGL backend | Linux: `sudo apt install gcc libgl1-mesa-dev xorg-dev` |
| **MSYS2 / MinGW-w64** | any | CGO for Windows | see below |
| **Xcode CLT** | any | CGO for macOS | `xcode-select --install` |
| `fyne` CLI | optional | Packaging into distributables | `go install fyne.io/tools/cmd/fyne@latest@latest` |

> Fyne uses CGO for OpenGL support. A C compiler is mandatory on every platform.

### Installing MinGW-w64 on Windows

**Option A — winget**:
```powershell
winget install --id MSYS2.MSYS2
# then inside MSYS2 shell:
pacman -S mingw-w64-x86_64-gcc
```

**Option B — portable **:
```powershell
# 1. Download the portable archive (via a local HTTP proxy if needed)
Invoke-WebRequest `
  -Uri "https://github.com/niXman/mingw-builds-binaries/releases/latest/download/x86_64-<ver>-release-posix-seh-ucrt-rt_v14-rev0.7z" `
  -Proxy "<PROXY>" `   # omit if no proxy is required
  -OutFile "$env:USERPROFILE\mingw64.7z"

# 2. Extract (7-Zip required — https://www.7-zip.org)
& "C:\Program Files\7-Zip\7z.exe" x "$env:USERPROFILE\mingw64.7z" -o"$env:USERPROFILE\mingw64" -y

# 3. Add to PATH for the current session
$env:PATH = "$env:USERPROFILE\mingw64\mingw64\bin;$env:PATH"
```
Pick the **posix-seh-ucrt** variant for best compatibility with Go's CGO.

### Network / proxy considerations

Go does not automatically inherit the Windows system proxy. If your environment uses a local HTTP proxy, set these before running `go`:

```powershell
$env:HTTPS_PROXY = "<PROXY>"   # adjust proxy address to match your proxy
$env:HTTP_PROXY  = "<PROXY>"

# If proxy.golang.org returns 403 (proxy intercepts TLS), bypass it entirely:
$env:GOPROXY     = "direct"
$env:GONOSUMDB   = "*"
```


---

## Quick start (development)

```bash
go mod tidy          # download all dependencies (one-time, requires internet)
go build .           # compile
./cornealius-eyeworth
```

## Running tests

```bash
go test ./...
```

**Windows** — run from PowerShell with GCC on `PATH`:

```powershell
# If behind a proxy, set these first (see Prerequisites → Network / proxy):
$env:PATH        = "$env:USERPROFILE\mingw64\mingw64\bin;$env:PATH"
$env:GOPROXY     = "direct"    # only needed if proxy.golang.org is blocked
$env:GONOSUMDB   = "*"         # only needed alongside GOPROXY=direct

go mod tidy
go build -ldflags "-s -w -H windowsgui" -o cornealius-eyeworth.exe .
```

---

## Makefile targets

The repository includes a `Makefile` with these useful targets:

```make
make deps        # go mod tidy
make build       # build stripped native binary
make build-win   # build a Windows GUI binary
make run         # build and run the app
make clean       # remove built binaries
```

---

## Building a standalone bundle

`fyne package` produces a self-contained app bundle:

```bash
# Windows — embed icon + single-file .exe (~25 MB, identical to stripped go build)
go run fyne.io/tools/cmd/fyne@latest package --release --app-id "com.github.mpyziak.cornealius-eyeworth" --os windows --icon assets/Logo.png --name "Cornealius Eyeworth"

# macOS — produces a .app bundle
go run fyne.io/tools/cmd/fyne@latest package --release --app-id "com.github.mpyziak.cornealius-eyeworth" --os darwin --icon assets/Logo.png --name "Cornealius Eyeworth"

# Linux — produces an executable with .desktop file
go run fyne.io/tools/cmd/fyne@latest package --release --app-id "com.github.mpyziak.cornealius-eyeworth" --os linux --icon assets/Logo.png --name "Cornealius Eyeworth"
```

---

## Fyne `WatchTheme` workaround (local module fork)

This repository uses a patched local copy of `fyne.io/fyne/v2` to fix a
Windows-only memory-leak bug in `WatchTheme` (see `LEAKS.md §1`).  The fix is
wired via a `replace` directive in `go.mod`:

```
replace fyne.io/fyne/v2 => ../fyne-v2-watchtheme-patch
```

The patched fork must exist as a sibling directory.  If you clone this
repository to a new machine or the sibling directory is missing, the build will
fail with:

```
go: ../fyne-v2-watchtheme-patch: reading ../fyne-v2-watchtheme-patch/go.mod: open ...: no such file or directory
```

**To recreate the fork:**

> **Note**: verified against `fyne.io/fyne/v2 v2.7.4` — the latest release at
> the time of writing.  Run `go list -m -versions fyne.io/fyne/v2` before
> starting; if a newer version is available, upgrade `go.mod` first (`go get
> fyne.io/fyne/v2@latest ; go mod tidy`) and re-apply the patch to the new
> module cache entry.

```powershell
# 1. Copy the module from the local cache (read-only by default — strip that)
$ver = "v2.7.4"
$src = "$env:USERPROFILE\go\pkg\mod\fyne.io\fyne\$ver"
$dst = "..\fyne-v2-watchtheme-patch"

Copy-Item -Recurse -Force $src $dst
attrib -r "$dst\*.*" /s /d

# 2. Open the single file that needs patching
#    c:\Users\...\Projects\spikes\fyne-v2-watchtheme-patch\internal\app\theme_windows.go
```

Apply the following diff to `internal/app/theme_windows.go` in the copied
directory (the patched version already committed in `../fyne-v2-watchtheme-patch`
can be used as a reference):

```diff
--- a/internal/app/theme_windows.go
+++ b/internal/app/theme_windows.go
@@ -4,6 +4,7 @@ package app
 import (
 	"syscall"
+	"time"
 
 	"golang.org/x/sys/windows/registry"
 
@@ -42,12 +42,41 @@ func isDark() bool {
 // WatchTheme calls the supplied function when the Windows dark/light theme changes.
 func WatchTheme(onChanged func()) {
-	// implementation based on an MIT-licensed Github Gist by Jeremy Black (c) 2022
-	// https://gist.github.com/jerblack/1d05bbcebb50ad55c312e4d7cf1bc909
-	...
+	openKey := func() (registry.Key, error) {
+		return registry.OpenKey(registry.CURRENT_USER, themeRegKey,
+			syscall.KEY_NOTIFY|registry.QUERY_VALUE)
+	}
+	...
 	k, err := registry.OpenKey(...)
 	if err != nil {
 		return
 	}
+
+	const debounce = 300 * time.Millisecond
+	var lastFired time.Time
+
 	for {
-		regNotifyChangeKeyValue.Call(uintptr(k), 0, 0x00000001|0x00000004, 0, 0)
-		onChanged()
+		ret, _, _ := regNotifyChangeKeyValue.Call(uintptr(k), 0, 0x00000001|0x00000004, 0, 0)
+		if ret != 0 { // ERROR_KEY_DELETED (1018) or other error
+			k.Close()
+			for {
+				time.Sleep(50 * time.Millisecond)
+				k, err = openKey()
+				if err == nil { break }
+			}
+			continue
+		}
+		if now := time.Now(); now.Sub(lastFired) >= debounce {
+			lastFired = now
+			onChanged()
+		}
 	}
 }
```

Two changes:
1. **`ERROR_KEY_DELETED` recovery** — Windows deletes and recreates `Themes\Personalize` during a full theme switch (Teams call-end, Focus Assist restore).  The original code ignores the `RegNotifyChangeKeyValue` return value, so when the key is deleted the call returns error 1018 immediately on every iteration — a busy-spin that enqueues thousands of `setupTheme` closures in under 200 ms.  The fix closes the stale handle and polls until the key is recreated.
2. **300 ms debounce** — Teams/Focus Assist writes the key 3–5 more times after recreation.  The debounce collapses the burst into a single `onChanged()` call.

```powershell
# 3. Verify the fork builds
cd ..\fyne-v2-watchtheme-patch
go build ./internal/app/ 2>&1

# 4. Back in the app — confirm both builds compile
cd ..\Cornealius-Eyeworth
go build -ldflags "-s -w -H windowsgui" -o "Cornealius Eyeworth.exe" .
go build -tags sysmon -ldflags "-s -w -H windowsgui" -o "Cornealius Eyeworth.exe" .
```

---

## Diagnostic build — OS event monitoring (`sysmon` tag)

The `sysmon` build tag compiles in additional goroutines that log Windows OS events to help diagnose memory leaks. It is **not** included in production builds.

What `sysmon` adds:
- Registry watcher on `HKCU\...\Themes\Personalize` (tracks burst writes that trigger Fyne's `setupTheme`)
- A pinned OS-thread message pump that logs: HID / network device arrivals and removals (`WM_DEVICECHANGE`), power sleep/resume (`WM_POWERBROADCAST`), session lock/unlock/RDP (`WM_WTSSESSION_CHANGE`), display change, DPI change, theme change, system colour change, font change, `WM_COMPACTING`, and Group Policy `WM_SETTINGCHANGE`
- Registry watcher on `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Group Policy\State` (GP update cycle ~90 min)
- Fyne `Settings.AddListener` to correlate `setupTheme` calls with the registry bursts above

```powershell
# Diagnostic build (Windows, with sysmon monitoring)
go build -tags sysmon -ldflags "-s -w -H windowsgui" -o "Cornealius Eyeworth.exe" .

# Production build (no extra overhead)
go build -ldflags "-s -w -H windowsgui" -o "Cornealius Eyeworth.exe" .
```

Log files are written to `%APPDATA%\Cornealius Eyeworth\logs\` regardless of the build tag (the base `applog` rotating-file logger and 30-second memory sampler are always active).

---

## Project structure

```
./
├── main.go                 Entry point — loads config, sets locale, starts UI
├── assets/
│   ├── assets.go           Embeds Logo.png as a Fyne resource
│   └── Logo.png            App logo (embedded at compile time)
├── config/config.go        Config struct + JSON persistence
├── i18n/strings.go         Localised strings
├── notifications/notifier.go  Cross-platform notifications via fyne.App
├── parsing/parser.go       CRON and minutes parsing
├── scheduling/
│   ├── scheduler.go        Hot-restartable CRON scheduler
│   ├── describer.go        Converts CRON to human-readable text
│   └── nexttrigger.go      Computes next trigger time
├── systray/systray.go      System tray integration and menu handling
├── ui/
│   ├── ui.go               Main window and app orchestration
│   ├── schedule.go         Schedule dialog UI
│   ├── language.go         Language selector UI
│   ├── about.go            About dialog
│   └── help.go             Help dialog
└── winres/                 Windows resource metadata for icon embedding
```

---

## CRON format

The app uses **Quartz-style 6-field** CRON expressions:

```
┌──────────── second (0-59)
│ ┌────────── minute (0-59)
│ │ ┌──────── hour   (0-23)
│ │ │ ┌────── day-of-month (1-31, or *)
│ │ │ │ ┌──── month  (1-12, or *)
│ │ │ │ │ ┌── day-of-week (0-6, or ?)
0 20,40,55 * * * ?
```

`?` is accepted as a wildcard (equivalent to `*`).  
The Simple mode lets you enter just the minute values (e.g. `20, 40, 55`).

---

## Adding a language

1. Add a new `Strings` variable in `i18n/strings.go`.
2. Add a `case "xx":` branch in `SetLanguage()`.
3. Add a row to `ui/language.go`'s `options` slice.
