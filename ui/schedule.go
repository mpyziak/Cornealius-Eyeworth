package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mpyziak/cornealius-eyeworth/config"
	"github.com/mpyziak/cornealius-eyeworth/diagnostics"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// ShowScheduleDialog opens the schedule-editing window, or focuses it if
// already open. deps.onFire must be the same callback ui.Run passed, not a
// local one.
func ShowScheduleDialog(app fyne.App, deps scheduleDeps, onSaved func(*config.Config)) {
	str := i18n.Active

	if focusExisting(dialogSchedule) {
		return
	}

	win := app.NewWindow(str.ScheduleDialogTitle)
	win.SetFixedSize(true)
	win.Resize(fyne.NewSize(500, 350))
	win.CenterOnScreen()

	currentCfg, err := deps.repo.Load()
	if err != nil {
		diagnostics.Err("open schedule dialog: cannot load config.json (%v)", err)
		win.Show()
		dialog.ShowError(err, win)
		return
	}

	simpleEyeMinutes, simpleEyeOk := scheduling.SimpleMinutes(currentCfg.CronExpression)
	simpleStandUpMinutes, simpleStandUpOk := scheduling.SimpleMinutes(currentCfg.StandUpCronExpression)
	startAdvanced := !simpleEyeOk || (currentCfg.StandUpCronExpression != "" && !simpleStandUpOk)

	eyeInstrLabel := widget.NewLabel(str.OptionsInstruction)
	eyeInstrLabel.Wrapping = fyne.TextWrapWord

	minutesEntry := widget.NewEntry()
	minutesEntry.SetPlaceHolder("20, 40, 55")
	if simpleEyeOk {
		minutesEntry.SetText(simpleEyeMinutes)
	}

	standUpInstrLabel := widget.NewLabel(str.OptionsStandUpInstruction)
	standUpInstrLabel.Wrapping = fyne.TextWrapWord

	standUpMinutesEntry := widget.NewEntry()
	standUpMinutesEntry.SetPlaceHolder("0, 15, 30, 45")
	if simpleStandUpOk {
		standUpMinutesEntry.SetText(simpleStandUpMinutes)
	}

	standardContent := container.NewPadded(container.New(layout.NewVBoxLayout(),
		eyeInstrLabel,
		minutesEntry,
		widget.NewSeparator(),
		standUpInstrLabel,
		standUpMinutesEntry,
	))

	eyeCronInstrLabel := widget.NewLabel(str.ScheduleCronInstruction)
	eyeCronEntry := widget.NewEntry()
	eyeCronEntry.SetText(currentCfg.CronExpression)

	standUpCronInstrLabel := widget.NewLabel(str.ScheduleCronInstruction)
	standUpCronEntry := widget.NewEntry()
	standUpCronEntry.SetText(currentCfg.StandUpCronExpression)

	advancedContent := container.NewPadded(container.New(layout.NewVBoxLayout(),
		eyeCronInstrLabel,
		eyeCronEntry,
		widget.NewSeparator(),
		standUpCronInstrLabel,
		standUpCronEntry,
	))

	tabs := container.NewAppTabs(
		container.NewTabItem(str.ScheduleStandardToggle, standardContent),
		container.NewTabItem(str.ScheduleAdvancedToggle, advancedContent),
	)
	if startAdvanced {
		tabs.SelectIndex(1)
	}

	errorLabel := widget.NewLabel("")
	errorLabel.Importance = widget.DangerImportance

	saveBtn := widget.NewButton(str.ButtonSave, func() {
		var eyeExpr, standUpExpr string
		var err error

		if tabs.SelectedIndex() == 0 {
			if eyeExpr, err = scheduling.ParseMinutes(minutesEntry.Text); err == nil && standUpMinutesEntry.Text != "" {
				standUpExpr, err = scheduling.ParseMinutes(standUpMinutesEntry.Text)
			}
		} else {
			if eyeExpr, err = scheduling.ParseCron(eyeCronEntry.Text); err == nil && strings.TrimSpace(standUpCronEntry.Text) != "" {
				standUpExpr, err = scheduling.ParseCron(standUpCronEntry.Text)
			}
		}

		if err != nil {
			errorLabel.SetText(err.Error())
			return
		}

		updated := &config.Config{
			CronExpression:        eyeExpr,
			StandUpCronExpression: standUpExpr,
			Language:              currentCfg.Language,
		}
		if saveErr := deps.repo.Save(updated); saveErr != nil {
			dialog.ShowError(saveErr, win)
			return
		}
		diagnostics.Event("schedule saved - eye=%s standUp=%s", updated.CronExpression, updated.StandUpCronExpression)

		if applyErr := scheduling.ApplySchedule(
			deps.eyeSched,
			deps.standUpSched,
			scheduling.Spec{
				EyeCron:     updated.CronExpression,
				StandUpCron: updated.StandUpCronExpression,
			},
			deps.onFire,
		); applyErr != nil {
			// Should be unreachable - both were validated above.
			diagnostics.Err("saved schedule rejected by scheduler: %v", applyErr)
			dialog.ShowError(applyErr, win)
			return
		}

		win.Close()
		onSaved(updated)
	})
	saveBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton(str.ButtonCancel, func() { win.Close() })

	btnRow := container.NewHBox(layout.NewSpacer(), saveBtn, cancelBtn)
	bottom := container.New(layout.NewVBoxLayout(), errorLabel, btnRow)

	win.SetContent(container.NewBorder(nil, bottom, nil, nil, tabs))
	registerDialog(dialogSchedule, win)
	win.Show()
}
