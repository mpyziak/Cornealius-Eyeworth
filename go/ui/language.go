package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mpyziak/cornealius-eyeworth/config"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

// ShowLanguageDialog opens a standalone language-selection window.
func ShowLanguageDialog(app fyne.App, repo *config.Repository) {
	S := i18n.Active

	currentCfg, err := repo.Load()
	if err != nil {
		return
	}

	type langItem struct {
		display string
		code    string
	}
	options := []langItem{
		{S.OptionsLanguageDefault, ""},
		{"English", "en"},
		{"Deutsch", "de"},
		{"Polski", "pl"},
	}

	labels := make([]string, len(options))
	for i, o := range options {
		labels[i] = o.display
	}

	currentIdx := 0
	if currentCfg.Language != nil {
		for i, o := range options {
			if o.code == *currentCfg.Language {
				currentIdx = i
				break
			}
		}
	}

	win := app.NewWindow(S.LanguageDialogTitle)
	win.SetFixedSize(true)
	win.Resize(fyne.NewSize(450, 150))
	win.CenterOnScreen()

	instrLabel := widget.NewLabel(S.OptionsLanguageLabel)

	selector := widget.NewSelect(labels, nil)
	selector.SetSelectedIndex(currentIdx)

	saveBtn := widget.NewButton(S.ButtonSave, func() {
		idx := selector.SelectedIndex()
		if idx < 0 {
			idx = 0
		}
		selected := options[idx]

		existing, loadErr := repo.Load()
		if loadErr != nil {
			dialog.ShowError(loadErr, win)
			return
		}

		var newLang *string
		if selected.code != "" {
			s := selected.code
			newLang = &s
		}

		langChanged := (newLang == nil) != (existing.Language == nil) ||
			(newLang != nil && existing.Language != nil && *newLang != *existing.Language)

		updated := &config.Config{
			CronExpression: existing.CronExpression,
			Language:       newLang,
		}
		if saveErr := repo.Save(updated); saveErr != nil {
			dialog.ShowError(saveErr, win)
			return
		}

		if langChanged {
			d := dialog.NewInformation(S.AppName, S.OptionsLanguageRestartNotice, win)
			d.SetOnClosed(func() { win.Close() })
			d.Show()
		} else {
			win.Close()
		}
	})
	saveBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton(S.ButtonCancel, func() { win.Close() })

	btnRow := container.NewHBox(layout.NewSpacer(), saveBtn, cancelBtn)
	win.SetContent(container.NewPadded(container.New(layout.NewVBoxLayout(),
		instrLabel,
		selector,
		widget.NewLabel(""), // spacer
		btnRow,
	)))
	win.Show()
}
