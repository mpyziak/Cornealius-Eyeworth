// Package assets embeds static resources (logo image) for use throughout the app.
package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed Logo.png
var logoBytes []byte

//go:embed Logo.ico
var iconBytes []byte

// Logo is the application logo as a Fyne static resource for UI elements.
var Logo fyne.Resource = fyne.NewStaticResource("Logo.png", logoBytes)

// AppIcon is the application icon used for the window/taskbar/notification icon.
var AppIcon fyne.Resource = fyne.NewStaticResource("Logo.ico", iconBytes)

// IconBytes returns the raw icon bytes for systray integration.
func IconBytes() []byte {
	return iconBytes
}
