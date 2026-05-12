// Package assets embeds static resources (logo image) for use throughout the app.
package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed Logo.png
var logoBytes []byte

// Logo is the application logo as a Fyne static resource.
var Logo fyne.Resource = fyne.NewStaticResource("Logo.png", logoBytes)
