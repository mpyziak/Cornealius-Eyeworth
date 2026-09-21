# Cornealius Eyeworth - Go + Fyne

Cross-platform app for eye-rest reminders. Builds from the repository root and runs on Windows, Linux, and macOS.

---

## Prerequisites

| Tool | Minimum version | Purpose | Install |
|------|----------------|---------|---------|
| **Go** | 1.22 | Build toolchain | https://go.dev/dl/ |
| **C compiler** | any | CGO - required by Fyne's OpenGL backend | Linux: `sudo apt install gcc libgl1-mesa-dev xorg-dev` |
| **MSYS2 / MinGW-w64** | any | CGO for Windows | see below |
| **Xcode CLT** | any | CGO for macOS | `xcode-select --install` |
| `fyne` CLI | optional | Packaging into bundles | `go install fyne.io/tools/cmd/fyne@latest` |

> Fyne uses CGO for OpenGL support. A C compiler is mandatory on every platform.

### Installing MinGW-w64 on Windows

**Option A - winget**:
```powershell
winget install --id MSYS2.MSYS2
# then inside MSYS2 shell:
pacman -S mingw-w64-x86_64-gcc
```

**Option B - portable**:
```powershell
# 1. Download the portable archive (via a local HTTP proxy if needed)
Invoke-WebRequest `
  -Uri "https://github.com/niXman/mingw-builds-binaries/releases/latest/download/x86_64-<ver>-release-posix-seh-ucrt-rt_v14-rev0.7z" `
  -Proxy "<PROXY>" `   # omit if no proxy is required
  -OutFile "$env:USERPROFILE\mingw64.7z"

# 2. Extract (7-Zip required - https://www.7-zip.org)
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

## Quick start (Linux / macOS)

```bash
go mod tidy          # download all dependencies (one-time, requires internet)
go build .           # compile
./cornealius-eyeworth
```

## Quick start (Windows)

From PowerShell, with GCC on `PATH`:

```powershell
# If behind a proxy, set these first (see Prerequisites → Network / proxy):
$env:PATH        = "$env:USERPROFILE\mingw64\mingw64\bin;$env:PATH"
$env:GOPROXY     = "direct"    # only needed if proxy.golang.org is blocked
$env:GONOSUMDB   = "*"         # only needed alongside GOPROXY=direct

go mod tidy
go build -ldflags "-s -w -H windowsgui" -o cornealius-eyeworth.exe .
```

`-H windowsgui` is what keeps a console window from opening behind the app. Drop
it while debugging if you want `stdout` back.

## Running tests

```bash
go test ./...
```

Nothing in the suite opens a window or a tray icon - it passes on Linux with
`DISPLAY` unset, so CI needs no display server.

---

## Makefile targets

```make
make deps                # go mod tidy
make winres              # regenerate rsrc_windows_amd64.syso (icon + manifest)
make build               # release build - no logging compiled in
make build-win           # release build, Windows GUI subsystem
make build-dev           # logging enabled (diagnostics tag)
make build-dev-win       # logging enabled, Windows GUI subsystem
make build-dev-sysmon    # logging + OS event monitors (Windows only)
make build-dev-sysmon-win # the same, Windows GUI subsystem
make patch-fyne          # recreate the Fyne WatchTheme fork (run once per machine)
make run                 # build and run the app
make upx / upx-win       # UPX-compress the built binary (needs upx on PATH)
make package-linux       # fyne package -> dist/
make package-darwin
make package-win         # go build + winres resources -> dist/
make clean               # remove built binaries and dist/
```

The `package-*` targets must run on the platform they name - none of this
cross-compiles, because Fyne needs CGO and a matching C toolchain.

---

## Building a standalone bundle

`make package-linux` and `make package-darwin` wrap `fyne package`, which produces a
`.tar.xz` with a `.desktop` file and a `.app` bundle respectively. Without `make`:

```bash
go run fyne.io/tools/cmd/fyne@latest package --release \
  --app-id "com.github.mpyziak.cornealius-eyeworth" \
  --os linux --icon assets/Logo.png --name "Cornealius Eyeworth"
```

**Not on Windows:** `make package-win` uses `go build`, not `fyne
package`, and writes `dist/Cornealius-Eyeworth.exe`. Without `make`:

```powershell
go build -ldflags "-s -w -H windowsgui" -o "dist\Cornealius-Eyeworth.exe" .
```

The icon and manifest come from `rsrc_windows_amd64.syso` (`make winres`), so
`fyne package` would embed a *second* manifest and the linker rejects duplicate. See also `RELEASES.md`.

---

## Fyne `WatchTheme` workaround (local module fork)

`go.mod` points Fyne at a patched sibling checkout that fixes a Windows-only
memory leak (`LEAKS.md §1`):

```
replace fyne.io/fyne/v2 => ../fyne-v2-watchtheme-patch
```

A fresh clone does not have it, and the build fails:

```
go: ../fyne-v2-watchtheme-patch: reading ../fyne-v2-watchtheme-patch/go.mod: open ...: no such file or directory
```

Fix it with `make patch-fyne`, or by hand on a box without `make`:

```powershell
# 1. Clone upstream Fyne at the pinned tag, beside this repository
git clone --depth 1 --branch v2.7.4 https://github.com/fyne-io/fyne.git ..\fyne-v2-watchtheme-patch

# 2. Overwrite the single file that carries the patch
copy external-patches\fyne-theme_windows.go ..\fyne-v2-watchtheme-patch\internal\app\theme_windows.go
```

That is all of it. To confirm:

```powershell
# The fork itself compiles
cd ..\fyne-v2-watchtheme-patch
go build ./internal/app/ 2>&1

# And the app builds against it, in all three profiles
cd ..\Cornealius-Eyeworth
go build -ldflags "-s -w -H windowsgui" -o "Cornealius Eyeworth.exe" .
go build -tags diagnostics -ldflags "-s -w -H windowsgui" -o "Cornealius Eyeworth.exe" .
go build -tags diagnostics,sysmon -ldflags "-s -w -H windowsgui" -o "Cornealius Eyeworth.exe" .
```

> **`v2.7.4` is written in four places** - the clone command above, `FYNE_TAG` in
> the `Makefile`, `FYNE_TAG` in `.github/workflows/release.yml`, and `go.mod`.
> Nothing checks that they agree. Change all four together, and re-check the patch
> against upstream: it replaces `theme_windows.go` wholesale, so an upstream
> rewrite is discarded silently instead of conflicting.

> Clone, don't copy out of the module cache. The `replace` in `go.mod` means Go
> never downloads Fyne, so the cache is empty exactly when you need it.

What the patch changes, for reference - step 2 already applied it:

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

1. **`ERROR_KEY_DELETED` recovery.** A full theme switch deletes and recreates
   `Themes\Personalize`. Upstream ignores the `RegNotifyChangeKeyValue` return, so
   while the key is gone it spins at CPU speed, queueing thousands of `setupTheme`
   closures. The patch reopens the handle instead.
2. **300 ms debounce.** Teams writes the key another 3–5 times after recreating it.

`LEAKS.md §1` has the log trace and the upstream bug report.

---

## Diagnostic builds

Logging and OS event monitoring are not in release builds at all. They live in
`dev-diagnostics/`, behind two tags:

| Tag | What it adds |
|-----|--------------|
| `diagnostics` | Rotating-file logger + 30-second memory sampler |
| `diagnostics,sysmon` | Plus OS-level event monitors (Windows only) |

`sysmon` alone does nothing - it only gates code inside `dev-diagnostics/`.

`sysmon` watches `Themes\Personalize` and the Group Policy `State` key, runs a
pinned message pump logging `WM_DEVICECHANGE`, `WM_POWERBROADCAST`,
`WM_WTSSESSION_CHANGE`, DPI/display/theme/font changes and `WM_COMPACTING`, and
hooks `Settings.AddListener` to line `setupTheme` calls up against the registry
bursts. That combination is what identified both leaks in `LEAKS.md`.

```powershell
# Standard build (no logging, no OS monitors)
go build -ldflags "-s -w -H windowsgui" -o "Cornealius Eyeworth.exe" .

# Diagnostic build - logging only
go build -tags diagnostics -ldflags "-s -w -H windowsgui" -o "Cornealius Eyeworth.exe" .

# Diagnostic build - logging + OS event monitoring (Windows only)
go build -tags diagnostics,sysmon -ldflags "-s -w -H windowsgui" -o "Cornealius Eyeworth.exe" .
```

Log files are written next to the executable, named by UTC timestamp (e.g. `cornealius-2026-06-22T14-30-00Z.log`).

---

## Project structure

```
./
├── main.go                 Entry point - loads config, sets locale and icon, starts UI
├── assets/
│   ├── assets.go           Embeds the logo, the 64x64 icon and the .ico
│   ├── gen/main.go         go:generate - rescales Logo.png to Logo64.png
│   └── Logo.png            Source logo, 800x800
├── config/                 Config struct, defaults, JSON persistence
├── i18n/strings.go         Localised strings - English, German, Polish
├── diagnostics/            Always-present logging interface; no-ops in release builds
├── dev-diagnostics/        The real logger + memory sampler, behind build tags
├── notifications/
│   ├── notifier.go         !windows - sends via fyne.App
│   ├── notifier_windows.go windows - Shell_NotifyIcon balloons via syscall
│   ├── flusher.go          Connects scheduling.Buffer to the senders above
│   ├── formatting.go       Batch → one title and one bulleted body
│   └── limits.go           Win32 balloon buffer sizes, asserted by tests
├── parsing/parser.go       CRON and minutes parsing
├── scheduling/
│   ├── scheduler.go        Hot-restartable CRON scheduler + ApplySchedule
│   ├── buffer.go           Coalesces firings inside a 2s window
│   ├── describer.go        Converts CRON to human-readable text
│   └── nexttrigger.go      Computes next trigger time
├── systray/
│   ├── systray.go          Tray icon, menu, on-demand status window
│   └── systraypatch_win.go HWND_MESSAGE reparent - see LEAKS.md
├── ui/
│   ├── ui.go               Run() - wires everything; anchor + window factory
│   ├── schedule.go         Schedule dialog
│   ├── language.go         Language selector
│   ├── dialogs.go          Single-instance registry for the dialogs above
│   ├── about.go            About dialog
│   └── help.go             Help dialog
├── external-patches/       The one patched Fyne file, copied into the fork
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
