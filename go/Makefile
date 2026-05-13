BINARY   = cornealius-eyeworth
ICON     = ../Properties/Resources/Logo.png

# -s strips the symbol table; -w strips DWARF debug info.
# Together they shave ~35-40% off a Fyne binary with no runtime cost.
LDFLAGS      = -ldflags "-s -w"
LDFLAGS_WIN  = -ldflags "-s -w -H windowsgui"

.PHONY: deps winres build run build-win upx upx-win package-win package-linux package-darwin cross-all clean

deps:
	go mod tidy

winres:
	go run github.com/tc-hib/go-winres@latest make

build:
	go build $(LDFLAGS) -o $(BINARY) .

build-win:
	go build $(LDFLAGS_WIN) -o $(BINARY).exe .

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
	go run fyne.io/fyne/v2/cmd/fyne package -os windows -icon $(ICON) -name "Cornealius Eyeworth"

package-linux: deps
	go run fyne.io/fyne/v2/cmd/fyne package -os linux -icon $(ICON) -name "Cornealius Eyeworth"

package-darwin: deps
	go run fyne.io/fyne/v2/cmd/fyne package -os darwin -icon $(ICON) -name "Cornealius Eyeworth"

cross-all: deps
	fyne-cross windows -arch amd64 -ldflags="-s -w"
	fyne-cross linux   -arch amd64 -ldflags="-s -w"
	fyne-cross darwin  -arch amd64,arm64 -ldflags="-s -w"

clean:
	$(RM) $(BINARY) $(BINARY).exe
