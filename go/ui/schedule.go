package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mpyziak/cornealius-eyeworth/config"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/parsing"
)

// ShowScheduleDialog opens a standalone schedule-editing window.
// onSaved is called with the updated *config.Config when the user saves.
func ShowScheduleDialog(app fyne.App, repo *config.Repository, onSaved func(*config.Config)) {
	S := i18n.Active

	currentCfg, err := repo.Load()
	if err != nil {
		return
	}

	simpleMinutes := parsing.TryExtractSimpleMinutes(currentCfg.CronExpression)
	startAdvanced := simpleMinutes == ""

	// Standard tab
	instrLabel := widget.NewLabel(S.OptionsInstruction)
	instrLabel.Wrapping = fyne.TextWrapWord

	minutesEntry := widget.NewEntry()
	minutesEntry.SetPlaceHolder("20, 40, 55")
	if simpleMinutes != "" {
		minutesEntry.SetText(simpleMinutes)
	}

	standardContent := container.NewPadded(container.New(layout.NewVBoxLayout(),
		instrLabel,
		minutesEntry,
	))

	// Advanced tab
	cronInstrLabel := widget.NewLabel(S.ScheduleCronInstruction)

	cronEntry := widget.NewEntry()
	cronEntry.SetText(currentCfg.CronExpression)

	advancedContent := container.NewPadded(container.New(layout.NewVBoxLayout(),
		cronInstrLabel,
		cronEntry,
	))

	tabs := container.NewAppTabs(
		container.NewTabItem(S.ScheduleStandardToggle, standardContent),
		container.NewTabItem(S.ScheduleAdvancedToggle, advancedContent),
	)
	if startAdvanced {
		tabs.SelectIndex(1)
	}

	errorLabel := widget.NewLabel("")
	errorLabel.Importance = widget.DangerImportance

	win := app.NewWindow(S.ScheduleDialogTitle)
	win.SetFixedSize(true)
	win.Resize(fyne.NewSize(450, 200))
	win.CenterOnScreen()

	saveBtn := widget.NewButton(S.ButtonSave, func() {
		var result parsing.ParseResult
		if tabs.SelectedIndex() == 0 {
			result = parsing.ParseMinutes(minutesEntry.Text)
		} else {
			result = parsing.ParseCron(cronEntry.Text)
		}

		if !result.Valid {
			errorLabel.SetText(result.Err)
			return
		}

		updated := &config.Config{
			CronExpression: result.Expression,
			Language:       currentCfg.Language,
		}
		if saveErr := repo.Save(updated); saveErr != nil {
			dialog.ShowError(saveErr, win)
			return
		}
		win.Close()
		onSaved(updated)
	})
	saveBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton(S.ButtonCancel, func() { win.Close() })

	btnRow := container.NewHBox(layout.NewSpacer(), saveBtn, cancelBtn)
	bottom := container.New(layout.NewVBoxLayout(), errorLabel, btnRow)

	win.SetContent(container.NewBorder(nil, bottom, nil, nil, tabs))
	win.Show()
}
