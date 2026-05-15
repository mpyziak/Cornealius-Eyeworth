package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

// ShowHelpDialog opens a standalone Help window.
func ShowHelpDialog(app fyne.App) {
	S := i18n.Active

	win := app.NewWindow(S.HelpDialogTitle)
	win.SetFixedSize(true)
	win.Resize(fyne.NewSize(450, 450))
	win.CenterOnScreen()

	titleLabel := widget.NewLabelWithStyle(S.AppTitle, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	bodyLabel := widget.NewLabel(S.HelpBody)
	bodyLabel.Wrapping = fyne.TextWrapWord

	scroll := container.NewScroll(bodyLabel)
	scroll.SetMinSize(fyne.NewSize(400, 350))

	closeBtn := widget.NewButton(S.ButtonClose, func() { win.Close() })
	btnRow := container.NewHBox(layout.NewSpacer(), closeBtn)

	win.SetContent(container.NewPadded(container.New(layout.NewVBoxLayout(),
		titleLabel,
		widget.NewSeparator(),
		scroll,
		btnRow,
	)))
	win.Show()
}
