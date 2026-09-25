package i18n

// Polish

var polish = Strings{
	AppName:   "Cornealius Eyeworth",
	AppTitle:  "Cornealius Eyeworth",
	GitHubUrl: "https://github.com/mpyziak/Cornealius-Eyeworth",

	StatusServing:       "● Na służbie",
	ScheduleDescription: "Harmonogram: %s",
	NextTrigger:         "Najbliższa chwila dla oczu: %s",
	NextTriggerStandUp:  "Najbliższa chwila dla pleców: %s",

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
	HelpBody: `Cornealius Eyeworth to Twój osobisty kamerdyner oka. Po uruchomieniu dyskretnie śledzi zegar i wyświetla krótkie, nieuciążliwe powiadomienia przypominające o podstawowej higienie wzroku - oderwij wzrok od ekranu, mrugnij i daj oczom odpocząć.

Masz pełną kontrolę:

• Harmonogram  (Opcje › Harmonogram)
  Wybierz minuty każdej godziny, w których Cornealius zadzwoni swoim dzwonkiem. Na przykład wpisanie 20, 40 przypomni Ci o XX:20 i XX:40 każdej godziny.  Zaawansowani użytkownicy mogą wpisać wyrażenie CRON bezpośrednio.

• Język  (Opcje › Język)
  Zmień język interfejsu. Wymagane jest ponowne uruchomienie, aby zmiana weszła w życie.

Cornealius działa cicho w tle. Wystarczy pozostawić okno otwarte (minimalizacja jest w porządku), a on zrobi resztę - z dużą klasą.`,

	NotificationOnDuty:          "Cornealius jest na służbie.",
	NotificationMinimizedToTray: "Cornealius nadal jest na służbie - znajdziesz go w zasobniku systemowym.",
	NotificationDistanceGlanceHeaders: []string{
		"Twoje oczy zasługują na przerwę.",
		"Chwila wytchnienia dla Twoich zmęczonych oczu.",
		"Oderwij wzrok od ekranu. Twój wzrok nalega.",
		"Zasada 20-20-20 czeka na Ciebie.",
	},
	NotificationMovementHeaders: []string{
		"Trochę ruchu poprawi Ci dzień.",
		"Twoje nogi proszą o chwilę przerwy.",
		"Krótki rozciąg poprawi samopoczucie.",
	},
	NotificationMovementQuips: []string{
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
		"Ludzki kręgosłup nie został zaprojektowany do pozycji siedzącej. Przypomnij mu o alternatywach.",
		"Krótki spacer nic nie kosztuje, a dodaje witalności.",
	},
	NotificationCombinedHeaders: []string{
		"Kilka uwag, jeśli mogę:",
		"Czy zechciałbyś zerknąć na to:",
		"Przesyłki do Twojej uwagi:",
	},
	NotificationDistanceGlanceQuips: []string{
		"Ludzkie oko nie zostało stworzone do wiecznego wpatrywania się w ekran.",
		"Mruganie jest bezpłatne. Używaj go obficie.",
		"Nawet najlepszy monokl wymaga okazjonalnej opieki.",
		"Widzenie w dal - luksus, na który Twoje oczy w pełni zasługują.",
		"Mięśnie akomodacji też potrzebują rozciągnięcia.",
		"20 metrów przestrzeni - balsam dla przepracowanego oka.",
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
		"Zmęczenie wzroku jest realne. Widok z okna również.",
		"Nawet najbardziej sumienny urzędnik musi od czasu do czasu oderwać wzrok od księgi.",
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
	ScheduleDescriptionSimple: "Przypomina w minutach: %s każdej godziny",
}
