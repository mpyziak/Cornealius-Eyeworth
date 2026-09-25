package i18n

// German

var german = Strings{
	AppName:   "Cornealius Eyeworth",
	AppTitle:  "Cornealius Eyeworth",
	GitHubUrl: "https://github.com/mpyziak/Cornealius-Eyeworth",

	StatusServing:       "● Im Dienst",
	ScheduleDescription: "Zeitplan: %s",
	NextTrigger:         "Nächste Augenpflege: %s",
	NextTriggerStandUp:  "Nächste Rückenpflege: %s",

	MenuOptions:      "Optionen ⚙️",
	MenuLanguage:     "Sprache... 🌐",
	MenuTriggerTimes: "Zeitplan... ⏰",
	MenuHelp:         "Hilfe 📖",
	MenuAbout:        "Über... ℹ️",
	MenuHowToUse:     "Verwendung... ❓",
	MenuGitHub:       "GitHub... 📄",
	MenuQuit:         "Beenden",

	TrayTooltipShowHide:    "Cornealius rufen/entlassen",
	TrayTooltipNextTrigger: "Nächste Erinnerungszeit",
	TrayTooltipQuit:        "Beenden",

	ScheduleDialogTitle:       "Zeitplan - Cornealius Eyeworth",
	OptionsInstruction:        "Minuten jeder Stunde, in denen Cornealius Sie an die Augenpause erinnern soll (z.B. 20, 40, 55):",
	OptionsStandUpInstruction: "Minuten jeder Stunde, zu denen Cornealius dich daran erinnern soll, aufzustehen und zu dehnen (z. B. 0, 15, 30, 45). Leer lassen zum Deaktivieren:",
	ScheduleCronInstruction:   `CRON-Ausdruck (z. B. "0 20,40,55 * * * ?"):`,
	ScheduleStandardToggle:    "Standard",
	ScheduleAdvancedToggle:    "Erweitert",
	ScheduleEyeToggle:         "Augenpflege",
	ScheduleStandUpToggle:     "Aufstehen",

	ButtonSave:   "Speichern",
	ButtonCancel: "Abbrechen",
	ButtonClose:  "Schließen",

	AboutDialogTitle: "Über - Cornealius Eyeworth",
	AboutDescription: "Ein distinguierter Augenbutler, der Sie in regelmäßigen Abständen an die Augenpause erinnert.  Der 20-20-20-Regel folgend, mit Stil.",
	AboutVersion:     "Version 1.0  -  © 2026 mpyziak",

	HelpDialogTitle: "Verwendung - Cornealius Eyeworth",
	HelpBody: `Cornealius Eyeworth ist Ihr persönlicher Augenbutler. Einmal gestartet, beobachtet er diskret die Uhr und liefert kurze, unaufdringliche Benachrichtigungen, die Sie an grundlegende Augenhygiene erinnern - schauen Sie vom Bildschirm weg, blinzeln Sie und lassen Sie Ihre Augen ausruhen.

Sie haben die volle Kontrolle:

• Zeitplan  (Optionen › Zeitplan)
  Wählen Sie die Minuten jeder Stunde, zu denen Cornealius seine Glocke läutet. Zum Beispiel erinnert Sie die Eingabe von 20, 40 jede Stunde um XX:20 und XX:40.  Fortgeschrittene Benutzer können direkt einen CRON-Ausdruck eingeben.

• Sprache  (Optionen › Sprache)
  Wechseln Sie die Oberflächensprache. Ein Neustart ist erforderlich, damit die Änderung wirksam wird.

Cornealius läuft still im Hintergrund. Lassen Sie das Fenster einfach offen (Minimieren ist in Ordnung) und er erledigt den Rest - mit beachtlichem Stil.`,

	NotificationOnDuty:          "Cornealius ist im Dienst.",
	NotificationMinimizedToTray: "Cornealius ist weiterhin im Dienst - Sie finden ihn in der Taskleiste.",
	NotificationDistanceGlanceHeaders: []string{
		"Ihre Augen verdienen eine Pause.",
		"Ein Moment der Erholung für Ihre müden Augen.",
		"Die 20-20-20-Regel erwartet Sie.",
	},
	NotificationMovementHeaders: []string{
		"Ein wenig Bewegung würde Ihrem Tag guttun.",
		"Deine Beine bitten um eine Pause.",
		"Ein kurzes Dehnen verbessert alles.",
	},
	NotificationMovementQuips: []string{
		"Zeit aufzustehen und zu dehnen.",
		"Deine Gliedmaßen brauchen Bewegung.",
		"Deine Körperhaltung verdient Aufmerksamkeit.",
		"Das sitzende Leben ruft nach einer Pause.",
		"Ein Spaziergang durch das Büro wird empfohlen.",
		"Cornealius empfiehlt einen kurzen Spaziergang.",
		"Dein Kreislauf könnte Hilfe gebrauchen.",
		"Langes Sitzen ist nicht ratsam. Stehe auf und bewege dich.",
		"Dein Rücken schätzt die aufrechte Körperhaltung.",
		"Ein Moment Bewegung wirkt Wunder für die Seele.",
		"Bewegung: Natur's unterschätztes Heilmittel.",
		"Der Stuhl ist kein Thron. Verlasse ihn regelmäßig.",
		"Die menschliche Wirbelsäule wurde nicht für die Sitzposition geschaffen. Erinnern Sie sie an Alternativen.",
		"Ein kurzer Spaziergang kostet nichts doch schenkt Vitalität.",
	},
	NotificationCombinedHeaders: []string{
		"Ein paar Hinweise, wenn ich bitten darf:",
		"Würden Sie einen Blick darauf werfen:",
		"Ein Paket zu Ihrer Aufmerksamkeit:",
	},
	NotificationDistanceGlanceQuips: []string{
		"Das menschliche Auge wurde nicht für ewiges Bildschirmstarren geschaffen.",
		"Blinzeln ist kostenlos. Nutzen Sie es reichlich.",
		"Selbst das feinste Monokel bedarf gelegentlicher Pflege.",
		"Fernsicht: ein Luxus, den Ihre Augen wirklich verdienen.",
		"Auch Ihre Fokusmuskulatur braucht Dehnung.",
		"20 Meter ins Nichts - ein Tonikum für die überarbeitete Netzhaut.",
		"Cornealius erinnert Sie: Bildschirme sind vergänglich, die Sehkraft soll nicht.",
		"Ein Blick in die Ferne kostet nur einen Moment.",
		"Ihre Augen haben Ihnen gut gedient. Erwidern Sie die Gunst.",
		"Starre-Wettbewerbe mit dem Monitor sind nicht ratsam. Sie werden verlieren.",
		"Ihre Netzhaut ist kein Abonnementdienst. Behandeln Sie sie entsprechend.",
		"Der Horizont existiert aus gutem Grund. Konsultieren Sie ihn gelegentlich.",
		"Pixel sind reichlich. Gute Sehkraft ist endlich.",
		"Cornealius merkt an: Blaulicht ist kein Ersatz für Sonnenlicht.",
		"Ein Fenster ist nicht nur dekorativ. Schauen Sie hindurch.",
		"Ihre Ziliarmuskel hätten gerne ein Wort. Treten Sie vom Bildschirm zurück.",
		"Augenermüdung ist real. Die Aussicht aus Ihrem Fenster ebenfalls.",
		"Selbst der gewissenhafteste Schreiber muss gelegentlich vom Hauptbuch aufsehen.",
	},

	ParseErrorNoMinutes:     "Cornealius besteht auf mindestens einer Minute. Er hat Ansprüche.",
	ParseErrorInvalidMinute: "„%s\" ist keine gültige Minute. Cornealius erwartet ganze Zahlen zwischen 0 und 59.",
	ParseErrorCronEmpty:     "Cornealius erwartet einen Ausdruck. Stille ist kein Zeitplan.",
	ParseErrorCronInvalid:   "„%s\" ist kein gültiger CRON-Ausdruck im Quartz-Format.",

	LanguageDialogTitle:          "Sprache - Cornealius Eyeworth",
	OptionsLanguageLabel:         "Sprache (erfordert Neustart):",
	OptionsLanguageDefault:       "Systemstandard",
	OptionsLanguageRestartNotice: "Die Sprachänderung wird beim nächsten Start wirksam.",

	ScheduleDescriptionCron:   "CRON: %s",
	ScheduleDescriptionSimple: "Ruft an den Minuten: %s jeder Stunde",
}
