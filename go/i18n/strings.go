// Package i18n provides localised strings for Cornealius Eyeworth.
// Call SetLanguage once at startup; then reference Active anywhere.
package i18n

import (
	"strings"

	"fyne.io/fyne/v2/lang"
)

// Strings holds every user-visible string for one locale.
type Strings struct {
	AppName   string
	AppTitle  string
	GitHubUrl string

	StatusServing string
	// Format args: %s = description (e.g. "20, 40, 55")
	ScheduleDescription string
	// Format args: %s = next trigger time formatted as HH:mm
	NextTrigger        string
	NextTriggerStandUp string

	MenuOptions      string
	MenuLanguage     string
	MenuTriggerTimes string
	MenuHelp         string
	MenuAbout        string
	MenuHowToUse     string
	MenuGitHub       string
	MenuQuit         string

	TrayTooltipShowHide    string
	TrayTooltipNextTrigger string
	TrayTooltipQuit        string

	ScheduleDialogTitle       string
	OptionsInstruction        string
	OptionsStandUpInstruction string
	ScheduleCronInstruction   string
	ScheduleStandardToggle    string
	ScheduleAdvancedToggle    string
	ScheduleEyeToggle         string
	ScheduleStandUpToggle     string

	ButtonSave   string
	ButtonCancel string
	ButtonClose  string

	AboutDialogTitle string
	AboutDescription string
	AboutVersion     string

	HelpDialogTitle string
	HelpBody        string

	NotificationOnDuty           string
	NotificationMinimizedToTray  string
	NotificationReminders        []string
	NotificationStandUpReminders []string
	NotificationEyeReminder      string
	NotificationStandUpReminder  string
	NotificationCombinedTitle    string
	NotificationQuips            []string

	ParseErrorNoMinutes     string
	ParseErrorInvalidMinute string // format: %s = invalid token
	ParseErrorCronEmpty     string
	ParseErrorCronInvalid   string // format: %s = invalid expression

	LanguageDialogTitle          string
	OptionsLanguageLabel         string
	OptionsLanguageDefault       string
	OptionsLanguageRestartNotice string

	ScheduleDescriptionCron   string // format: %s = raw cron
	ScheduleDescriptionSimple string // format: %s = "20, 40, 55"
}

// Active points to the currently selected locale. Defaults to English.
var Active = &english

// SetLanguage switches the active locale.
// Pass nil for the OS default locale; if the system locale is unsupported, default to English.
func SetLanguage(language *string) {
	localeCode := ""
	if language == nil {
		localeCode = lang.SystemLocale().LanguageString()
	} else {
		localeCode = strings.ToLower(strings.TrimSpace(*language))
	}

	switch localeCode {
	case "de", "de-de", "de-at", "de-ch", "de-li", "de-lu":
		Active = &german
	case "pl", "pl-pl":
		Active = &polish
	default:
		Active = &english
	}
}

// LanguageOption represents one entry in the language selector.
type LanguageOption struct {
	DisplayName string
	Code        string // empty string = system default
}

// AvailableLanguages is the list shown in the Language dialog.
var AvailableLanguages = []LanguageOption{
	{"", ""}, // display name filled at runtime from Active.OptionsLanguageDefault
	{"English", "en"},
	{"Deutsch", "de"},
	{"Polski", "pl"},
}

// ---------------------------------------------------------------------------
// English (default)
// ---------------------------------------------------------------------------

var english = Strings{
	AppName:   "Cornealius Eyeworth",
	AppTitle:  "Cornealius Eyeworth",
	GitHubUrl: "https://github.com/mpyziak/Cornealius-Eyeworth",

	StatusServing:       "● Serving",
	ScheduleDescription: "Schedule: %s",
	NextTrigger:         "Next trigger: %s",
	NextTriggerStandUp:  "Next stand-up: %s",

	MenuOptions:      "Options ⚙️",
	MenuLanguage:     "Language... 🌐",
	MenuTriggerTimes: "Schedule... ⏰",
	MenuHelp:         "Help 📖",
	MenuAbout:        "About... ℹ️",
	MenuHowToUse:     "How to use... ❓",
	MenuGitHub:       "GitHub... 📄",
	MenuQuit:         "Quit",

	TrayTooltipShowHide:    "Call/Dismiss Cornealius",
	TrayTooltipNextTrigger: "Next reminder time",
	TrayTooltipQuit:        "Quit the application",

	ScheduleDialogTitle:       "Schedule — Cornealius Eyeworth",
	OptionsInstruction:        "Minutes of each hour at which Cornealious shall remind you to rest your eyes (e.g. 20, 40, 55):",
	OptionsStandUpInstruction: "Minutes of each hour at which Cornealius shall remind you to stand and stretch (e.g. 0, 15, 30, 45). Leave empty to disable:",
	ScheduleCronInstruction:   `CRON (e.g. "0 20,40,55 * * * ?"):`,
	ScheduleStandardToggle:    "Standard",
	ScheduleAdvancedToggle:    "Advanced",
	ScheduleEyeToggle:         "Eye Care",
	ScheduleStandUpToggle:     "Stand-up",

	ButtonSave:   "Save",
	ButtonCancel: "Cancel",
	ButtonClose:  "Close",

	AboutDialogTitle: "About — Cornealius Eyeworth",
	AboutDescription: "A distinguished ocular butler who reminds you to rest your eyes at regular intervals.  Following the 20-20-20 Rule, with decorum.",
	AboutVersion:     "Version 1.0  -  © 2026 mpyziak",

	HelpDialogTitle: "How to use — Cornealius Eyeworth",
	HelpBody: `Cornealius Eyeworth is your personal ocular butler. Once running, he discreetly watches the clock and delivers short, unintrusive notifications reminding you to practise basic eye hygiene — look away from your screen, blink, and let your eyes rest.

You are in full control:

• Schedule  (Options › Schedule)
  Choose the minutes of each hour at which Cornealius shall ring his bell. For example, entering 20, 40 will remind you at XX:20 and XX:40 every hour.  Advanced users may enter a CRON expression directly.

• Language  (Options › Language)
  Switch the interface language. A restart is required for the change to take effect.

Cornealius runs quietly in the background. Simply leave the window open (minimising is fine) and he will do the rest — with considerable decorum.`,

	NotificationOnDuty:          "Cornealius is on duty.",
	NotificationMinimizedToTray: "Cornealius is still on duty — find him in the system tray.",
	NotificationReminders: []string{
		"Your eyes deserve an intermission.",
		"A moment of respite for your weary eyes.",
		"The 20-20-20 rule awaits.",
	},
	NotificationStandUpReminders: []string{
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
	},
	NotificationEyeReminder:     "Rest your eyes.",
	NotificationStandUpReminder: "Move about.",
	NotificationCombinedTitle:   "Wellness Check",
	NotificationQuips: []string{
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
	ScheduleDescriptionSimple: "Minutes %s past every hour",
}

// ---------------------------------------------------------------------------
// German
// ---------------------------------------------------------------------------

var german = Strings{
	AppName:   "Cornealius Eyeworth",
	AppTitle:  "Cornealius Eyeworth",
	GitHubUrl: "https://github.com/mpyziak/Cornealius-Eyeworth",

	StatusServing:       "● Im Dienst",
	ScheduleDescription: "Zeitplan: %s",
	NextTrigger:         "Nächste Erinnerung: %s",
	NextTriggerStandUp:  "Nächstes Aufstehen: %s",

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
	HelpBody: `Cornealius Eyeworth ist Ihr persönlicher Augenbutler. Einmal gestartet, beobachtet er diskret die Uhr und liefert kurze, unaufdringliche Benachrichtigungen, die Sie an grundlegende Augenhygiene erinnern — schauen Sie vom Bildschirm weg, blinzeln Sie und lassen Sie Ihre Augen ausruhen.

Sie haben die volle Kontrolle:

• Zeitplan  (Optionen › Zeitplan)
  Wählen Sie die Minuten jeder Stunde, zu denen Cornealius seine Glocke läutet. Zum Beispiel erinnert Sie die Eingabe von 20, 40 jede Stunde um XX:20 und XX:40.  Fortgeschrittene Benutzer können direkt einen CRON-Ausdruck eingeben.

• Sprache  (Optionen › Sprache)
  Wechseln Sie die Oberflächensprache. Ein Neustart ist erforderlich, damit die Änderung wirksam wird.

Cornealius läuft still im Hintergrund. Lassen Sie das Fenster einfach offen (Minimieren ist in Ordnung) und er erledigt den Rest — mit beachtlichem Stil.`,

	NotificationOnDuty:          "Cornealius ist im Dienst.",
	NotificationMinimizedToTray: "Cornealius ist weiterhin im Dienst — Sie finden ihn in der Taskleiste.",
	NotificationReminders: []string{
		"Ihre Augen verdienen eine Pause.",
		"Ein Moment der Erholung für Ihre müden Augen.",
		"Die 20-20-20-Regel erwartet Sie.",
	},
	NotificationStandUpReminders: []string{
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
	},
	NotificationEyeReminder:     "Ruhe deine Augen aus.",
	NotificationStandUpReminder: "Bewege dich herum.",
	NotificationCombinedTitle:   "Gesundheitsprüfung",
	NotificationQuips: []string{
		"Das menschliche Auge wurde nicht für ewiges Bildschirmstarren geschaffen.",
		"Blinzeln ist kostenlos. Nutzen Sie es reichlich.",
		"Selbst das feinste Monokel bedarf gelegentlicher Pflege.",
		"Fernsicht: ein Luxus, den Ihre Augen wirklich verdienen.",
		"Auch Ihre Fokusmuskulatur braucht Dehnung.",
		"20 Meter ins Nichts — ein Tonikum für die überarbeitete Netzhaut.",
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
	ScheduleDescriptionSimple: "Minuten %s jeder Stunde",
}

// ---------------------------------------------------------------------------
// Polish
// ---------------------------------------------------------------------------

var polish = Strings{
	AppName:   "Cornealius Eyeworth",
	AppTitle:  "Cornealius Eyeworth",
	GitHubUrl: "https://github.com/mpyziak/Cornealius-Eyeworth",

	StatusServing:       "● Na służbie",
	ScheduleDescription: "Harmonogram: %s",
	NextTrigger:         "Następne przypomnienie: %s",
	NextTriggerStandUp:  "Następne wstanie: %s",

	MenuOptions:      "Opcje ⚙️",
	MenuLanguage:     "Język... 🌐",
	MenuTriggerTimes: "Harmonogram... ⏰",
	MenuHelp:         "Pomoc 📖",
	MenuAbout:        "O programie... ℹ️",
	MenuHowToUse:     "Jak używać... ❓",
	MenuGitHub:       "GitHub... 📄",
	MenuQuit:         "Zamknij",

	TrayTooltipShowHide:    "Zawołaj/Odeślij Cornealiusa",
	TrayTooltipNextTrigger: "Czas następnego przypomnienia",
	TrayTooltipQuit:        "Zamknij aplikację",

	ScheduleDialogTitle:       "Harmonogram - Cornealius Eyeworth",
	OptionsInstruction:        "Minuty każdej godziny, w których Cornealius przypomni Ci o odpoczynku dla oczu (np. 20, 40, 55):",
	OptionsStandUpInstruction: "Minuty każdej godziny, w których Cornealius ma Cię przypomnieć o wstaniu i rozciągnięciu (np. 0, 15, 30, 45). Pozostaw puste, aby wyłączyć:",
	ScheduleCronInstruction:   `Wyrażenie CRON (np. "0 20,40,55 * * * ?"):`,
	ScheduleStandardToggle:    "Standardowe",
	ScheduleAdvancedToggle:    "Zaawansowane",
	ScheduleEyeToggle:         "Opieka oczna",
	ScheduleStandUpToggle:     "Wstań",

	ButtonSave:   "Zapisz",
	ButtonCancel: "Anuluj",
	ButtonClose:  "Zamknij",

	AboutDialogTitle: "O programie - Cornealius Eyeworth",
	AboutDescription: "Wybitny kamerdyner, który przypomina o regularnym odpoczynku dla oczu.  Zgodnie z zasadą 20-20-20, z klasą.",
	AboutVersion:     "Wersja 1.0  -  © 2026 mpyziak",

	HelpDialogTitle: "Jak używać - Cornealius Eyeworth",
	HelpBody: `Cornealius Eyeworth to Twój osobisty kamerdyner oka. Po uruchomieniu dyskretnie śledzi zegar i wyświetla krótkie, nieuciążliwe powiadomienia przypominające o podstawowej higienie wzroku — oderwij wzrok od ekranu, mrugnij i daj oczom odpocząć.

Masz pełną kontrolę:

• Harmonogram  (Opcje › Harmonogram)
  Wybierz minuty każdej godziny, w których Cornealius zadzwoni swoim dzwonkiem. Na przykład wpisanie 20, 40 przypomni Ci o XX:20 i XX:40 każdej godziny.  Zaawansowani użytkownicy mogą wpisać wyrażenie CRON bezpośrednio.

• Język  (Opcje › Język)
  Zmień język interfejsu. Wymagane jest ponowne uruchomienie, aby zmiana weszła w życie.

Cornealius działa cicho w tle. Wystarczy pozostawić okno otwarte (minimalizacja jest w porządku), a on zrobi resztę — z dużą klasą.`,

	NotificationOnDuty:          "Cornealius jest na służbie.",
	NotificationMinimizedToTray: "Cornealius nadal jest na służbie — znajdziesz go w zasobniku systemowym.",
	NotificationReminders: []string{
		"Twoje oczy zasługują na przerwę.",
		"Chwila wytchnienia dla Twoich zmęczonych oczu.",
		"Oderwij wzrok od ekranu. Twój wzrok nalega.",
		"Zasada 20-20-20 czeka na Ciebie.",
	},
	NotificationStandUpReminders: []string{
		"Czas wstać i się rozciągnąć.",
		"Twoje kończyny potrzebują ruchu.",
		"Twoja postawa zasługuje na uwagę.",
		"Siedząca praca wymaga przerwy.",
		"Spacer po biurze jest zalecany.",
		"Cornealius sugeruje krótki spacer.",
		"Twój układ krążenia potrzebuje pomocy.",
		"Długie siedzenie jest niewskazane. Wstań i się poruszaj.",
		"Twoje plecy doceniają pozycję pionową.",
		"Moment ruchu robi cuda dla duszy.",
		"Ruch: niedoceniany lek natury.",
		"Krzesło to nie tron. Opuszczaj je okresowo.",
	},
	NotificationEyeReminder:     "Odpoczną twoje oczy.",
	NotificationStandUpReminder: "Poruś się.",
	NotificationCombinedTitle:   "Kontrola zdrowia",
	NotificationQuips: []string{
		"Ludzkie oko nie zostało stworzone do wiecznego wpatrywania się w ekran.",
		"Mruganie jest bezpłatne. Używaj go obficie.",
		"Nawet najlepszy monokl wymaga okazjonalnej opieki.",
		"Widzenie w dal — luksus, na który Twoje oczy w pełni zasługują.",
		"Mięśnie akomodacji też potrzebują rozciągnięcia.",
		"20 metrów przestrzeni — balsam dla przepracowanego oka.",
		"Cornealius przypomina: ekrany przemijają, wzrok nie powinien.",
		"Krótkie spojrzenie w dal kosztuje jedynie chwilę.",
		"Wpatrywanie się w monitor to zawód z góry przesądzony. Przegrasz.",
		"Twoja siatkówka nie jest usługą abonamentową. Traktuj ją odpowiednio.",
		"Horyzont istnieje z powodu. Zaglądaj do niego od czasu do czasu.",
		"Piksele są w nadmiarze. Dobry wzrok jest ograniczony.",
		"Cornealius zauważa: niebieskie światło nie zastąpi słonecznego.",
		"Okno to nie tylko dekoracja. Sprawdź, co za nim.",
		"Twoje mięśnie rzęskowe mają coś do powiedzenia. Odejdź od ekranu.",
		"Odpocznij teraz. Pasek przewijania będzie czekał na Twój powrót.",
	},

	ParseErrorNoMinutes:     "Cornealius nalega na co najmniej jedną minutę. Ma swoje standardy.",
	ParseErrorInvalidMinute: "„%s\" to nieprawidłowa minuta. Cornealius oczekuje liczb całkowitych od 0 do 59.",
	ParseErrorCronEmpty:     "Cornealius wymaga wyrażenia CRON. Cisza to nie harmonogram.",
	ParseErrorCronInvalid:   "\"%s\" nie jest prawidłowym wyrażeniem CRON w formacie Quartz.",

	LanguageDialogTitle:          "Język - Cornealius Eyeworth",
	OptionsLanguageLabel:         "Język (wymaga ponownego uruchomienia):",
	OptionsLanguageDefault:       "Domyślny systemowy",
	OptionsLanguageRestartNotice: "Zmiana języka zostanie zastosowana przy następnym uruchomieniu.",

	ScheduleDescriptionCron:   "CRON: %s",
	ScheduleDescriptionSimple: "Minuty %s każdej godziny",
}
