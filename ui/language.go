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

func languageChanged(existingLang, newLang *string) bool {
	return (newLang == nil) != (existingLang == nil) ||
		(newLang != nil && existingLang != nil && *newLang != *existingLang)
}

func updatedConfigForLanguage(existing *config.Config, newLang *string) *config.Config {
	return &config.Config{
		CronExpression:        existing.CronExpression,
		StandUpCronExpression: existing.StandUpCronExpression,
		Language:              newLang,
	}
}

// ShowLanguageDialog opens a standalone language-selection window.
func ShowLanguageDialog(app fyne.App, repo config.Store) {
	S := i18n.Active

	currentCfg, err := repo.Load()
	if err != nil {
		return
	}

	options := i18n.AvailableLanguages
	labels := make([]string, len(options))
	for i, opt := range options {
		if opt.DisplayName == "" {
			labels[i] = S.OptionsLanguageDefault
		} else {
			labels[i] = opt.DisplayName
		}
	}

	currentIdx := 0
	if currentCfg.Language != nil {
		for i, opt := range options {
			if opt.Code == *currentCfg.Language {
				currentIdx = i
				break
			}
		}
	}

	win := app.NewWindow(S.LanguageDialogTitle)
	win.SetFixedSize(true)
	win.Resize(fyne.NewSize(450, 175))
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
		if selected.Code != "" {
			s := selected.Code
			newLang = &s
		}

		langChanged := languageChanged(existing.Language, newLang)

		updated := updatedConfigForLanguage(existing, newLang)
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
		layout.NewSpacer(),
		btnRow,
	)))
	win.Show()
}
