package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

// ShowAboutDialog opens a standalone About window.
func ShowAboutDialog(app fyne.App) {
	S := i18n.Active

	win := app.NewWindow(S.AboutDialogTitle)
	win.SetFixedSize(true)
	win.Resize(fyne.NewSize(450, 225))
	win.CenterOnScreen()

	titleLabel := widget.NewLabelWithStyle(S.AppTitle, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	descLabel := widget.NewLabel(S.AboutDescription)
	descLabel.Wrapping = fyne.TextWrapWord

	versionLabel := widget.NewLabelWithStyle(S.AboutVersion, fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	closeBtn := widget.NewButton(S.ButtonClose, func() { win.Close() })
	btnRow := container.NewHBox(layout.NewSpacer(), closeBtn)

	win.SetContent(container.NewPadded(container.New(layout.NewVBoxLayout(),
		titleLabel,
		widget.NewSeparator(),
		descLabel,
		versionLabel,
		btnRow,
	)))
	win.Show()
}
