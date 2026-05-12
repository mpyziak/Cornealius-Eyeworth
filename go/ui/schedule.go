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

// ShowScheduleDialog shows a modal schedule-editing dialog over parent.
// onSaved is called with the updated *config.Config when the user saves.
func ShowScheduleDialog(parent fyne.Window, repo *config.Repository, onSaved func(*config.Config)) {
	S := i18n.Active

	currentCfg, err := repo.Load()
	if err != nil {
		dialog.ShowError(err, parent)
		return
	}

	simpleMinutes := parsing.TryExtractSimpleMinutes(currentCfg.CronExpression)
	startAdvanced := simpleMinutes == ""

	// ── Standard tab ───────────────────────────────────────────────────────
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

	// ── Advanced tab ───────────────────────────────────────────────────────
	cronInstrLabel := widget.NewLabel(S.ScheduleCronInstruction)

	cronEntry := widget.NewEntry()
	cronEntry.SetText(currentCfg.CronExpression)

	advancedContent := container.NewPadded(container.New(layout.NewVBoxLayout(),
		cronInstrLabel,
		cronEntry,
	))

	// ── Tabs ───────────────────────────────────────────────────────────────
	tabs := container.NewAppTabs(
		container.NewTabItem(S.ScheduleStandardToggle, standardContent),
		container.NewTabItem(S.ScheduleAdvancedToggle, advancedContent),
	)
	if startAdvanced {
		tabs.SelectIndex(1)
	}

	// ── Error label ────────────────────────────────────────────────────────
	errorLabel := widget.NewLabel("")
	errorLabel.Importance = widget.DangerImportance

	// ── Buttons ────────────────────────────────────────────────────────────
	var dlg *dialog.CustomDialog

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
			dialog.ShowError(saveErr, parent)
			return
		}
		dlg.Hide()
		onSaved(updated)
	})
	saveBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton(S.ButtonCancel, func() {
		dlg.Hide()
	})

	btnRow := container.NewHBox(layout.NewSpacer(), saveBtn, cancelBtn)
	bottom := container.New(layout.NewVBoxLayout(), errorLabel, btnRow)

	content := container.NewBorder(nil, bottom, nil, nil, tabs)

	dlg = dialog.NewCustomWithoutButtons(S.ScheduleDialogTitle, content, parent)
	dlg.Resize(fyne.NewSize(420, 260))
	dlg.Show()
}
