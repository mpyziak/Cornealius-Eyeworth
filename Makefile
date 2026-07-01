BINARY   = cornealius-eyeworth
APP_ID   = com.github.mpyziak.cornealius-eyeworth
ICON     = assets/Logo.png
APP_NAME = "Cornealius Eyeworth"

# -s strips the symbol table; -w strips DWARF debug info.
# Together they shave ~35-40% off a Fyne binary with no runtime cost.
LDFLAGS      = -ldflags "-s -w"
LDFLAGS_WIN  = -ldflags "-s -w -H windowsgui"

# Diagnostic build tags (not included in default release builds).
# diagnostics: enables rotating-file logging and memory sampler (dev-diagnostics/).
# diagnostics + sysmon: also enables OS-level event monitoring (Windows only).
DIAG_TAGS        = -tags diagnostics
DIAG_SYSMON_TAGS = -tags diagnostics,sysmon

# Fyne fork paths for the WatchTheme patch (see patches/fyne-theme_windows.go).
FYNE_VER  = v2@v2.7.4
FYNE_SRC  = $(shell go env GOPATH)/pkg/mod/fyne.io/fyne/$(FYNE_VER)
FYNE_FORK = ../fyne-v2-watchtheme-patch

.PHONY: deps winres build build-win build-dev build-dev-win build-dev-sysmon build-dev-sysmon-win upx upx-win run patch-fyne package-win package-linux package-darwin cross-all clean

deps:
	go mod tidy

winres:
	go run github.com/tc-hib/go-winres@latest make

# ── Release builds (no diagnostics compiled in) ───────────────────────────────

build:
	go build $(LDFLAGS) -o $(BINARY) .

build-win:
	go build $(LDFLAGS_WIN) -o $(BINARY).exe .

# ── Diagnostic builds ────────────────────────────────────────────────────────

build-dev:
	go build $(LDFLAGS) $(DIAG_TAGS) -o $(BINARY) .

build-dev-win:
	go build $(LDFLAGS_WIN) $(DIAG_TAGS) -o $(BINARY).exe .

build-dev-sysmon:
	go build $(LDFLAGS) $(DIAG_SYSMON_TAGS) -o $(BINARY) .

build-dev-sysmon-win:
	go build $(LDFLAGS_WIN) $(DIAG_SYSMON_TAGS) -o $(BINARY).exe .

# ── Fyne WatchTheme fork setup ───────────────────────────────────────────────
# Creates ../fyne-v2-watchtheme-patch from the local module cache and applies
# the debounce + ERROR_KEY_DELETED-recovery patch to internal/app/theme_windows.go.
# Run once after cloning on a new machine. Safe to re-run if FYNE_FORK is absent.

patch-fyne:
	@if [ -d "$(FYNE_FORK)" ]; then echo "$(FYNE_FORK) already exists; delete it first to re-patch"; exit 1; fi
	cp -r "$(FYNE_SRC)" "$(FYNE_FORK)"
	chmod -R u+w "$(FYNE_FORK)"
	cp external-patches/fyne-theme_windows.go "$(FYNE_FORK)/internal/app/theme_windows.go"
	@echo "Fyne fork ready at $(FYNE_FORK)"

# Run UPX over the already-stripped binary for another ~50% reduction.
# Requires UPX to be on PATH: https://github.com/upx/upx/releases
# --best is slower to pack but produces the smallest file; safe to omit.
upx: build
	upx --best $(BINARY)

upx-win: build-win
	upx --best $(BINARY).exe

run: build
	./$(BINARY)

# fyne package strips by default; pass extra ldflags via FyneFlags if needed.
package-win: deps
	go run fyne.io/tools/cmd/fyne@latest package --release --app-id $(APP_ID) --os windows --icon $(ICON) --name $(APP_NAME)
	mkdir -p dist
	mv $(APP_NAME).exe dist/$(APP_NAME).exe

package-linux: deps
	go run fyne.io/tools/cmd/fyne@latest package --release --app-id $(APP_ID) --os linux --icon $(ICON) --name $(APP_NAME)
	mkdir -p dist
	mv $(APP_NAME).tar.xz dist/$(APP_NAME).tar.xz

package-darwin: deps
	go run fyne.io/tools/cmd/fyne@latest package --release --app-id $(APP_ID) --os darwin --icon $(ICON) --name $(APP_NAME)
	mkdir -p dist
	mv $(APP_NAME).app dist/$(APP_NAME).app

clean:
	$(RM) $(BINARY) $(BINARY).exe
	$(RM) -rf dist
