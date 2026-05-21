# Cornealius Eyeworth — Go + Fyne

Cross-platform rewrite of Cornealius Eyeworth in Go with [Fyne](https://fyne.io/).  
Runs on Windows, Linux, and macOS from a single codebase.

---

## Prerequisites

| Tool | Minimum version | Purpose | Install |
|------|----------------|---------|---------|
| **Go** | 1.22 | Build toolchain | https://go.dev/dl/ |
| **GCC (MinGW-w64)** | any | CGO — required by Fyne's OpenGL bindings **(Windows)** | see below |
| **GCC + OpenGL headers** | any | CGO — required by Fyne's OpenGL bindings **(Linux)** | `sudo apt install gcc libgl1-mesa-dev xorg-dev` |
| **Xcode CLT** | any | CGO — required by Fyne's OpenGL bindings **(macOS)** | `xcode-select --install` |
| `fyne` CLI | optional | Packaging into distributable bundles | `go install fyne.io/fyne/v2/cmd/fyne@latest` |

> Fyne uses CGO to call OpenGL. A C compiler is **mandatory** on every platform.

### Installing MinGW-w64 on Windows

**Option A — winget** (if available):
```powershell
winget install --id MSYS2.MSYS2
# then inside MSYS2 shell:
pacman -S mingw-w64-x86_64-gcc
```

**Option B — portable, no installer needed** (useful in managed/restricted environments):
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

# 4. (Optional) persist permanently
[System.Environment]::SetEnvironmentVariable(
  "PATH",
  "$env:USERPROFILE\mingw64\mingw64\bin;" + [System.Environment]::GetEnvironmentVariable("PATH","User"),
  "User"
)
```
Pick the **posix-seh-ucrt** variant from the [releases page](https://github.com/niXman/mingw-builds-binaries/releases) — it is the most compatible with Go's CGO.

### Network / proxy considerations

Go's build toolchain does **not** read the Windows system proxy automatically.  
If your machine uses a local HTTP proxy (common with Clash, corporate VPN agents, etc.) you must set the environment variables explicitly before running any `go` command:

```powershell
$env:HTTPS_PROXY = "<PROXY>"   # adjust proxy address to match your proxy
$env:HTTP_PROXY  = "<PROXY>"

# If proxy.golang.org returns 403 (proxy intercepts TLS), bypass it entirely:
$env:GOPROXY  = "direct"
$env:GONOSUMDB = "*"
```

To detect what proxy Windows would use for a given URL:
```powershell
[System.Net.WebRequest]::DefaultWebProxy.GetProxy("https://proxy.golang.org")
```

---

## Quick start (development)

```bash
cd go
go mod tidy          # download all dependencies (one-time, requires internet)
go build .           # compile
./cornealius-eyeworth  # run (config.json must be in the same directory)
```

**Windows** — GCC must be on PATH and the console window should be suppressed:

```powershell
cd go

# If behind a proxy, set these first (see Prerequisites → Network / proxy):
$env:PATH        = "$env:USERPROFILE\mingw64\mingw64\bin;$env:PATH"
$env:GOPROXY     = "direct"    # only needed if proxy.golang.org is blocked
$env:GONOSUMDB   = "*"         # only needed alongside GOPROXY=direct

go mod tidy
go build -ldflags "-s -w -H windowsgui" -o cornealius-eyeworth.exe .
```

---

## Building a standalone binary

### Windows (no console window, single .exe)

```powershell
# GCC must be on PATH — see Prerequisites
$env:PATH = "$env:USERPROFILE\mingw64\mingw64\bin;$env:PATH"

go build -ldflags "-s -w -H windowsgui" -o cornealius-eyeworth.exe .
# copy config.json next to the .exe
```

### Linux

```bash
go build -ldflags "-s -w" -o cornealius-eyeworth .
# copy config.json next to the binary
```

### macOS

```bash
go build -ldflags "-s -w" -o cornealius-eyeworth .
# copy config.json next to the binary
```

### Optional: UPX compression

[UPX](https://github.com/upx/upx/releases) compresses a stripped Fyne binary by another ~50%.
Install it and run it over the already-built binary:

```bash
upx --best cornealius-eyeworth        # Linux / macOS
upx --best cornealius-eyeworth.exe    # Windows
```

Or use the Makefile targets:

```bash
make upx        # build + strip + UPX (Linux / macOS)
make upx-win    # build + strip + UPX (Windows)
```

> **Note:** UPX-packed executables may trigger some antivirus scanners as false positives.
> For internal/personal use it is fine; for public distribution prefer the uncompressed stripped binary or `fyne package`.

---

## Packaging with the Fyne CLI (recommended for distribution)

`fyne package` produces a self-contained app bundle (`.exe` on Windows,
`.app` on macOS, executable with icon on Linux).

Run without installing — use `go run` (like `npx` for Node.js):

```bash
# Windows — embed icon + single-file .exe (~25 MB, identical to stripped go build)
go run fyne.io/tools/cmd/fyne@latest package -release -app-id "github.com/mpyziak/cornealius-eyeworth" -os windows -icon assets/Logo.png -name "Cornealius Eyeworth"

# macOS — produces a .app bundle
go run fyne.io/tools/cmd/fyne@latest package -release -app-id "github.com/mpyziak/cornealius-eyeworth" -os darwin  -icon assets/Logo.png -name "Cornealius Eyeworth"

# Linux — produces an executable with .desktop file
go run fyne.io/tools/cmd/fyne@latest package -release -app-id "github.com/mpyziak/cornealius-eyeworth" -os linux   -icon assets/Logo.png -name "Cornealius Eyeworth"
```

> **Note:** `-release` enables dead-code stripping and `-app-id` sets the Windows application identity, keeping the output at ~25 MB — the same size as a stripped `go build`.


Copy `config.json` next to the produced binary / into the bundle before distributing.

---

## Cross-compilation

Use [fyne-cross](https://github.com/fyne-io/fyne-cross) (requires Docker):

```bash
go install github.com/fyne-io/fyne-cross@latest
cd go

fyne-cross windows -arch amd64
fyne-cross linux   -arch amd64
fyne-cross darwin  -arch amd64,arm64
```

---

## Project structure

```
go/
├── main.go                 Entry point — loads config, sets locale, starts UI
├── config.json             Default runtime config (copied next to the binary)
├── assets/
│   ├── assets.go           Embeds Logo.png as a Fyne resource
│   └── Logo.png            App logo (embedded at compile time)
├── config/config.go        Config struct + JSON persistence
├── i18n/strings.go         Localised strings: English, Deutsch, Polski
├── notifications/notifier.go  Cross-platform notifications via fyne.App
├── parsing/parser.go       Cron + minutes input validation and conversion
├── scheduling/
│   ├── scheduler.go        Hot-restartable CRON scheduler (robfig/cron)
│   ├── describer.go        Converts CRON to human-readable text
│   └── nexttrigger.go      Computes next trigger time
└── ui/
    ├── ui.go               Run() — orchestrates scheduler, window, notifications
    ├── schedule.go         Schedule dialog (Standard / Advanced tabs)
    ├── language.go         Language selector dialog
    ├── about.go            About dialog
    └── help.go             Help dialog
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

1. Add a new `Strings` variable in `i18n/strings.go` (copy from `english`).
2. Add a `case "xx":` branch in `SetLanguage()`.
3. Add a row to `ui/language.go`'s `options` slice.
