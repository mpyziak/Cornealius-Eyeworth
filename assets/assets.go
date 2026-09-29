package assets

//go:generate go run ./gen

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed Logo.png
var logoBytes []byte

//go:embed Logo64.png
var iconBytes64 []byte

//go:embed Logo.ico
var iconBytes []byte

// Logo is the full-resolution app logo, for canvas elements only.
var Logo fyne.Resource = fyne.NewStaticResource("Logo.png", logoBytes)

// Icon is the one used at runtime. Regenerate with `go generate ./assets/...`.
var Icon fyne.Resource = fyne.NewStaticResource("Logo64.png", iconBytes64)

// AppIcon is the .ico form of the logo. Unused: nothing in Fyne or the
// stdlib decodes ICO. Kept pending a decision on pointing winres at it.
var AppIcon fyne.Resource = fyne.NewStaticResource("Logo.ico", iconBytes)
