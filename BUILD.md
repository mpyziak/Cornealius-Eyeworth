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
