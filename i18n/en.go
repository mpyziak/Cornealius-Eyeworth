package i18n

// English (default)

var english = Strings{
	AppName:   "Cornealius Eyeworth",
	AppTitle:  "Cornealius Eyeworth",
	GitHubURL: "https://github.com/mpyziak/Cornealius-Eyeworth",

	StatusServing:      "● Serving",
	NextTrigger:        "Next eye care: %s",
	NextTriggerStandUp: "Next back care: %s",

	MenuOptions:      "Options ⚙️",
	MenuLanguage:     "Language... 🌐",
	MenuTriggerTimes: "Schedule... ⏰",
	MenuHelp:         "Help 📖",
	MenuAbout:        "About... ℹ️",
	MenuHowToUse:     "How to use... ❓",
	MenuGitHub:       "GitHub... 📄",
	MenuQuit:         "Quit",

	TrayTooltipShowHide: "Call/Dismiss Cornealius",

	ScheduleDialogTitle:       "Schedule - Cornealius Eyeworth",
	OptionsInstruction:        "Minutes of each hour at which Cornealious shall remind you to rest your eyes (e.g. 20, 40, 55):",
	OptionsStandUpInstruction: "Minutes of each hour at which Cornealius shall remind you to stand and stretch (e.g. 0, 15, 30, 45). Leave empty to disable:",
	ScheduleCronInstruction:   `CRON (e.g. "0 20,40,55 * * * ?"):`,
	ScheduleStandardToggle:    "Standard",
	ScheduleAdvancedToggle:    "Advanced",

	ButtonSave:   "Save",
	ButtonCancel: "Cancel",
	ButtonClose:  "Close",

	AboutDialogTitle: "About - Cornealius Eyeworth",
	AboutDescription: "A distinguished ocular butler who reminds you to rest your eyes at regular intervals.  Following the 20-20-20 Rule, with decorum.",
	AboutVersion:     "Version 1.0  -  © 2026 mpyziak",

	HelpDialogTitle: "How to use - Cornealius Eyeworth",
	HelpBody: `Cornealius Eyeworth is your personal ocular butler. Once running, he discreetly watches the clock and delivers short, unintrusive notifications reminding you to practise basic eye hygiene - look away from your screen, blink, and let your eyes rest.

You are in full control:

• Schedule  (Options › Schedule)
  Choose the minutes of each hour at which Cornealius shall ring his bell. For example, entering 20, 40 will remind you at XX:20 and XX:40 every hour.  Advanced users may enter a CRON expression directly.

• Language  (Options › Language)
  Switch the interface language. A restart is required for the change to take effect.

Cornealius runs quietly in the background. Simply leave the window open (minimising is fine) and he will do the rest - with considerable decorum.`,

	NotificationOnDuty:          "Cornealius is on duty.",
	NotificationMinimizedToTray: "Cornealius is still on duty - find him in the system tray.",
	NotificationDistanceGlanceHeaders: []string{
		"Your eyes deserve an intermission.",
		"A moment of respite for your weary eyes.",
		"The 20-20-20 rule awaits.",
	},
	NotificationMovementHeaders: []string{
		"A little motion would make your day.",
		"Your legs are petitioning for a pause.",
		"A short stretch will improve everything.",
	},
	NotificationMovementQuips: []string{
		"Time to stand and stretch.",
		"Movement awaits your limbs.",
		"Your posture requires attention.",
		"The sedentary life calls for a recess.",
		"A stroll about the office is recommended.",
		"Cornealius suggests a brief constitutional.",
		"Your circulation could use the assistance.",
		"Prolonged sitting is inadvisable. Rise and move.",
		"Your back appreciates vertical orientation.",
		"A moment's ambulation does wonders for the soul.",
		"Movement: nature's most underrated medicine.",
		"The chair is not a throne. Vacate it periodically.",
		"The human spine was not engineered for the seated position.",
		"A brief perambulation costs nothing yet buys considerable vitality.",
	},
	NotificationCombinedHeaders: []string{
		"Some pieces of advice, if I may:",
		"Would you lay your eyes on these:",
		"A packet for your attention:",
	},
	NotificationDistanceGlanceQuips: []string{
		"The human eye was not designed for eternal screen-gazing.",
		"Blinking is free. Use it liberally.",
		"Even the finest monocle requires occasional polishing.",
		"Distance vision: a luxury your eyes richly deserve.",
		"Your focus muscles need stretching too.",
		"Cornealius reminds you: screens are temporary, eyesight should not.",
		"A brief gaze into the distance costs nothing but a moment.",
		"Your eyes have served you well. Return the favour.",
		"Staring contests with your monitor are inadvisable. You will lose.",
		"Your retinas are not subscription-based. Treat them accordingly.",
		"The horizon exists for a reason. Consult it occasionally.",
		"Pixels are plentiful. Good eyesight is finite.",
		"Cornealius notes: blue light is not a substitute for sunlight.",
		"A window is not merely decorative. Gaze through it.",
		"Your ciliary muscles would like a word. Step away from the screen.",
		"Rest now. The scroll bar will still be there when you return.",
		"Focus fatigue is real. The view from your window is also real.",
		"Even the most dedicated clerk must occasionally look up from the ledger.",
	},

	ParseErrorNoMinutes:     "Cornealius insists on at least one minute. He has standards.",
	ParseErrorInvalidMinute: "'%s' is not a valid minute. Cornealius expects whole numbers between 0 and 59.",
	ParseErrorCronEmpty:     "Cornealious requires a CRON expression. Silence is not a schedule.",
	ParseErrorCronInvalid:   "\u201c%s\u201d is not a valid Quartz CRON expression.",

	LanguageDialogTitle:          "Language \u2014 Cornealius Eyeworth",
	OptionsLanguageLabel:         "Language (requires restart):",
	OptionsLanguageDefault:       "System default",
	OptionsLanguageRestartNotice: "The language change will take effect the next time Cornealius starts.",

	ScheduleDescriptionCron:   "CRON: %s",
	ScheduleDescriptionSimple: "Calls at minutes %s of every hour",
}
