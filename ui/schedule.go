package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mpyziak/cornealius-eyeworth/config"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/notifications"
	"github.com/mpyziak/cornealius-eyeworth/parsing"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// ShowScheduleDialog opens a standalone schedule-editing window.
// onSaved is called with the updated *config.Config when the user saves.
func ShowScheduleDialog(app fyne.App, repo config.Store, eyeSched, standUpSched scheduling.Runner, onSaved func(*config.Config), buffer scheduling.ReminderAggregator) {
	S := i18n.Active

	currentCfg, err := repo.Load()
	if err != nil {
		return
	}

	simpleEyeMinutes := parsing.TryExtractSimpleMinutes(currentCfg.CronExpression)
	simpleStandUpMinutes := parsing.TryExtractSimpleMinutes(currentCfg.StandUpCronExpression)
	startAdvanced := simpleEyeMinutes == "" || (currentCfg.StandUpCronExpression != "" && simpleStandUpMinutes == "")

	eyeInstrLabel := widget.NewLabel(S.OptionsInstruction)
	eyeInstrLabel.Wrapping = fyne.TextWrapWord

	minutesEntry := widget.NewEntry()
	minutesEntry.SetPlaceHolder("20, 40, 55")
	if simpleEyeMinutes != "" {
		minutesEntry.SetText(simpleEyeMinutes)
	}

	standUpInstrLabel := widget.NewLabel(S.OptionsStandUpInstruction)
	standUpInstrLabel.Wrapping = fyne.TextWrapWord

	standUpMinutesEntry := widget.NewEntry()
	standUpMinutesEntry.SetPlaceHolder("0, 15, 30, 45")
	if simpleStandUpMinutes != "" {
		standUpMinutesEntry.SetText(simpleStandUpMinutes)
	}

	standardContent := container.NewPadded(container.New(layout.NewVBoxLayout(),
		eyeInstrLabel,
		minutesEntry,
		widget.NewSeparator(),
		standUpInstrLabel,
		standUpMinutesEntry,
	))

	eyeCronInstrLabel := widget.NewLabel(S.ScheduleCronInstruction)
	eyeCronEntry := widget.NewEntry()
	eyeCronEntry.SetText(currentCfg.CronExpression)

	standUpCronInstrLabel := widget.NewLabel(S.ScheduleCronInstruction)
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
	win.Resize(fyne.NewSize(500, 350))
	win.CenterOnScreen()

	saveBtn := widget.NewButton(S.ButtonSave, func() {
		var eyeResult parsing.ParseResult
		var standUpResult parsing.ParseResult

		if tabs.SelectedIndex() == 0 {
			eyeResult = parsing.ParseMinutes(minutesEntry.Text)
			if parsed := parsing.ParseMinutes(standUpMinutesEntry.Text); parsed.Valid {
				standUpResult = parsed
			} else if standUpMinutesEntry.Text == "" {
				standUpResult = parsing.ParseResult{Valid: true, Expression: ""}
			} else {
				standUpResult = parsed
			}
		} else {
			eyeResult = parsing.ParseCron(eyeCronEntry.Text)
			if strings.TrimSpace(standUpCronEntry.Text) == "" {
				standUpResult = parsing.ParseResult{Valid: true, Expression: ""}
			} else {
				standUpResult = parsing.ParseCron(standUpCronEntry.Text)
			}
		}

		if !eyeResult.Valid {
			errorLabel.SetText(eyeResult.Err)
			return
		}

		if !standUpResult.Valid {
			errorLabel.SetText(standUpResult.Err)
			return
		}

		updated := &config.Config{
			CronExpression:        eyeResult.Expression,
			StandUpCronExpression: standUpResult.Expression,
			Language:              currentCfg.Language,
		}
		if saveErr := repo.Save(updated); saveErr != nil {
			dialog.ShowError(saveErr, win)
			return
		}

		scheduling.ApplySchedule(
			eyeSched,
			standUpSched,
			scheduling.ScheduleSpec{
				EyeCron:     updated.CronExpression,
				StandUpCron: updated.StandUpCronExpression,
			},
			func(notificationCategory string) {
				buffer.Add(notifications.NewReminder(notificationCategory))
			},
		)

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
