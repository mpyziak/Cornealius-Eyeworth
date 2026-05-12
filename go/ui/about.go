package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

// ShowAboutDialog displays application information over parent.
func ShowAboutDialog(parent fyne.Window) {
	S := i18n.Active

	titleLabel := widget.NewLabelWithStyle(S.AppTitle, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	descLabel := widget.NewLabel(S.AboutDescription)
	descLabel.Wrapping = fyne.TextWrapWord

	versionLabel := widget.NewLabelWithStyle(S.AboutVersion, fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	content := container.NewPadded(container.New(layout.NewVBoxLayout(),
		titleLabel,
		widget.NewSeparator(),
		descLabel,
		versionLabel,
	))

	dialog.ShowCustom(S.AboutDialogTitle, S.ButtonClose, content, parent)
}
