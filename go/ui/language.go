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

// ShowLanguageDialog shows the language-selection dialog over parent.
func ShowLanguageDialog(parent fyne.Window, repo *config.Repository) {
	S := i18n.Active

	currentCfg, err := repo.Load()
	if err != nil {
		dialog.ShowError(err, parent)
		return
	}

	// Build option list: first entry is "system default", then explicit locales.
	type langItem struct {
		display string
		code    string // "" = system default
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

	// Find which option matches the current config.
	currentIdx := 0
	if currentCfg.Language != nil {
		for i, o := range options {
			if o.code == *currentCfg.Language {
				currentIdx = i
				break
			}
		}
	}

	instrLabel := widget.NewLabel(S.OptionsLanguageLabel)

	selector := widget.NewSelect(labels, nil)
	selector.SetSelectedIndex(currentIdx)

	var dlg *dialog.CustomDialog

	saveBtn := widget.NewButton(S.ButtonSave, func() {
		idx := selector.SelectedIndex()
		if idx < 0 {
			idx = 0
		}
		selected := options[idx]

		existing, loadErr := repo.Load()
		if loadErr != nil {
			dialog.ShowError(loadErr, parent)
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
			dialog.ShowError(saveErr, parent)
			return
		}
		dlg.Hide()

		if langChanged {
			dialog.ShowInformation(S.AppName, S.OptionsLanguageRestartNotice, parent)
		}
	})
	saveBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton(S.ButtonCancel, func() {
		dlg.Hide()
	})

	btnRow := container.NewHBox(layout.NewSpacer(), saveBtn, cancelBtn)
	content := container.NewPadded(container.New(layout.NewVBoxLayout(),
		instrLabel,
		selector,
		btnRow,
	))

	dlg = dialog.NewCustomWithoutButtons(S.LanguageDialogTitle, content, parent)
	dlg.Resize(fyne.NewSize(320, 180))
	dlg.Show()
}
