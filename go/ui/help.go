package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

// ShowHelpDialog displays the "How to use" content over parent.
func ShowHelpDialog(parent fyne.Window) {
	S := i18n.Active

	titleLabel := widget.NewLabelWithStyle(S.AppTitle, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	bodyLabel := widget.NewLabel(S.HelpBody)
	bodyLabel.Wrapping = fyne.TextWrapWord

	scroll := container.NewScroll(bodyLabel)
	scroll.SetMinSize(fyne.NewSize(400, 300))

	content := container.NewPadded(container.New(layout.NewVBoxLayout(),
		titleLabel,
		widget.NewSeparator(),
		scroll,
	))

	dialog.ShowCustom(S.HelpDialogTitle, S.ButtonClose, content, parent)
}
