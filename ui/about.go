package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

// ShowAboutDialog opens the About window, or focuses it if already open.
func ShowAboutDialog(app fyne.App) {
	str := i18n.Active

	if focusExisting(dialogAbout) {
		return
	}

	win := app.NewWindow(str.AboutDialogTitle)
	win.SetFixedSize(true)
	win.Resize(fyne.NewSize(450, 225))
	win.CenterOnScreen()

	titleLabel := widget.NewLabelWithStyle(str.AppTitle, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	descLabel := widget.NewLabel(str.AboutDescription)
	descLabel.Wrapping = fyne.TextWrapWord

	versionLabel := widget.NewLabelWithStyle(str.AboutVersion, fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	closeBtn := widget.NewButton(str.ButtonClose, func() { win.Close() })
	btnRow := container.NewHBox(layout.NewSpacer(), closeBtn)

	win.SetContent(container.NewPadded(container.New(layout.NewVBoxLayout(),
		titleLabel,
		widget.NewSeparator(),
		descLabel,
		versionLabel,
		btnRow,
	)))
	registerDialog(dialogAbout, win)
	win.Show()
}
