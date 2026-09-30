package ui

import (
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mpyziak/cornealius-eyeworth/config"
	"github.com/mpyziak/cornealius-eyeworth/diagnostics"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

// languageIndexForLabel finds the index of label in labels, the same order
// ShowLanguageDialog builds the radio group's options in. -1 (label not
// found, e.g. nothing selected yet) falls back to 0, the system default.
func languageIndexForLabel(labels []string, label string) int {
	idx := slices.Index(labels, label)
	if idx == -1 {
		return 0
	}
	return idx
}

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

// Normalising both sides, so a hand-edited "PL" or "pl_PL" picks the same entry
// i18n.SetLanguage will. 0 is the system default.
func selectedLanguageIndex(options []i18n.LanguageOption, lang *string) int {
	if lang == nil {
		return 0
	}
	want := i18n.NormaliseLocale(*lang)
	for i, opt := range options {
		if opt.Code != "" && i18n.NormaliseLocale(opt.Code) == want {
			return i
		}
	}
	return 0
}

// ShowLanguageDialog opens the language-selection window, or focuses it if
// already open.
func ShowLanguageDialog(app fyne.App, repo *config.Repository) {
	str := i18n.Active

	if focusExisting(dialogLanguage) {
		return
	}

	win := app.NewWindow(str.LanguageDialogTitle)
	win.SetFixedSize(true)

	currentCfg, err := repo.Load()
	if err != nil {
		diagnostics.Err("open language dialog: cannot load config.json (%v)", err)
		win.Resize(fyne.NewSize(450, 175))
		win.CenterOnScreen()
		win.Show()
		dialog.ShowError(err, win)
		return
	}

	options := i18n.AvailableLanguages()
	labels := make([]string, len(options))
	for i, opt := range options {
		if opt.DisplayName == "" {
			labels[i] = str.OptionsLanguageDefault
		} else {
			labels[i] = opt.DisplayName
		}
	}

	currentIdx := selectedLanguageIndex(options, currentCfg.Language)

	instrLabel := widget.NewLabel(str.OptionsLanguageLabel)

	radio := widget.NewRadioGroup(labels, nil)
	radio.Required = true
	radio.SetSelected(labels[currentIdx])

	// Cap the visible list at ~8 rows; beyond that it scrolls instead of
	// growing the window past the screen.
	radioHeight := radio.MinSize().Height
	rowHeight := radioHeight / float32(len(labels))
	if maxHeight := rowHeight * 8; radioHeight > maxHeight {
		radioHeight = maxHeight
	}
	radioScroll := container.NewVScroll(radio)
	radioScroll.SetMinSize(fyne.NewSize(400, radioHeight))

	saveBtn := widget.NewButton(str.ButtonSave, func() {
		idx := languageIndexForLabel(labels, radio.Selected)
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
			d := dialog.NewInformation(str.AppName, str.OptionsLanguageRestartNotice, win)
			d.SetOnClosed(func() { win.Close() })
			d.Show()
		} else {
			win.Close()
		}
	})
	saveBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton(str.ButtonCancel, func() { win.Close() })

	btnRow := container.NewHBox(layout.NewSpacer(), saveBtn, cancelBtn)
	content := container.NewPadded(container.New(layout.NewVBoxLayout(),
		instrLabel,
		radioScroll,
		layout.NewSpacer(),
		btnRow,
	))
	win.SetContent(content)
	win.Resize(fyne.NewSize(450, content.MinSize().Height))
	win.CenterOnScreen()
	registerDialog(dialogLanguage, win)
	win.Show()
}
