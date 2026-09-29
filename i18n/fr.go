package i18n

// French
//
// The spaces before : ; ? ! and inside « » are U+00A0 (no-break), per French
// typography. They look like ordinary spaces here; keep them when editing.

var french = Strings{
	AppName:   "Cornealius Eyeworth",
	AppTitle:  "Cornealius Eyeworth",
	GitHubURL: "https://github.com/mpyziak/Cornealius-Eyeworth",

	StatusServing:      "● En service",
	NextTrigger:        "Prochain soin des yeux : %s",
	NextTriggerStandUp: "Prochain soin du dos : %s",

	MenuOptions:      "Options ⚙️",
	MenuLanguage:     "Langue... 🌐",
	MenuTriggerTimes: "Horaires... ⏰",
	MenuHelp:         "Aide 📖",
	MenuAbout:        "À propos... ℹ️",
	MenuHowToUse:     "Mode d’emploi... ❓",
	MenuGitHub:       "GitHub... 📄",
	MenuQuit:         "Quitter",

	TrayTooltipShowHide: "Appeler/Congédier Cornealius",

	ScheduleDialogTitle:       "Horaires - Cornealius Eyeworth",
	OptionsInstruction:        "Minutes de chaque heure auxquelles Cornealius vous rappellera de reposer vos yeux (p. ex. 20, 40, 55) :",
	OptionsStandUpInstruction: "Minutes de chaque heure auxquelles Cornealius vous rappellera de vous lever et de vous étirer (p. ex. 0, 15, 30, 45). Laisser vide pour désactiver :",
	ScheduleCronInstruction:   "Expression CRON (p. ex. « 0 20,40,55 * * * ? ») :",
	ScheduleStandardToggle:    "Standard",
	ScheduleAdvancedToggle:    "Avancé",

	ButtonSave:   "Enregistrer",
	ButtonCancel: "Annuler",
	ButtonClose:  "Fermer",

	AboutDialogTitle: "À propos - Cornealius Eyeworth",
	AboutDescription: "Un majordome oculaire distingué qui vous rappelle de reposer vos yeux à intervalles réguliers.  Selon la règle des 20-20-20, avec décorum.",
	AboutVersion:     "Version 1.0  -  © 2026 mpyziak",

	HelpDialogTitle: "Mode d’emploi - Cornealius Eyeworth",
	HelpBody: `Cornealius Eyeworth est votre majordome oculaire personnel. Une fois lancé, il surveille discrètement l’horloge et vous adresse de courtes notifications, jamais importunes, qui vous rappellent les règles élémentaires d’hygiène visuelle - détournez le regard de l’écran, clignez des yeux et laissez-les se reposer.

Vous gardez le contrôle total :

• Horaires  (Options › Horaires)
  Choisissez les minutes de chaque heure auxquelles Cornealius fera tinter sa clochette. Par exemple, saisir 20, 40 vous rappellera à XX:20 et XX:40 chaque heure.  Les utilisateurs avancés peuvent saisir directement une expression CRON.

• Langue  (Options › Langue)
  Changez la langue de l’interface. Un redémarrage est nécessaire pour que la modification prenne effet.

Cornealius travaille discrètement en arrière-plan. Laissez simplement la fenêtre ouverte (la réduire ne pose aucun problème) et il s’occupera du reste - avec un décorum considérable.`,

	NotificationOnDuty:          "Cornealius est en service.",
	NotificationMinimizedToTray: "Cornealius est toujours en service - retrouvez-le dans la zone de notification.",
	NotificationDistanceGlanceHeaders: []string{
		"Vos yeux méritent un entracte.",
		"Un moment de répit pour vos yeux fatigués.",
		"La règle des 20-20-20 vous attend.",
	},
	NotificationMovementHeaders: []string{
		"Un peu de mouvement égaierait votre journée.",
		"Vos jambes sollicitent une pause.",
		"Un court étirement arrangera bien des choses.",
	},
	NotificationMovementQuips: []string{
		"Il est temps de vous lever et de vous étirer.",
		"Le mouvement attend vos membres.",
		"Votre posture requiert votre attention.",
		"La vie sédentaire réclame une suspension de séance.",
		"Un petit tour du bureau est recommandé.",
		"Cornealius suggère une brève promenade de santé.",
		"Votre circulation sanguine apprécierait un coup de main.",
		"Rester longtemps assis est déconseillé. Levez-vous et bougez.",
		"Votre dos apprécie la position verticale.",
		"Quelques pas font des merveilles pour l’âme.",
		"Le mouvement : le remède le plus sous-estimé de la nature.",
		"La chaise n’est pas un trône. Quittez-la de temps à autre.",
		"La colonne vertébrale humaine n’a pas été conçue pour la position assise.",
		"Une brève déambulation ne coûte rien et procure une vitalité considérable.",
	},
	NotificationCombinedHeaders: []string{
		"Quelques conseils, si vous me permettez :",
		"Daignez poser les yeux sur ceci :",
		"Un pli à votre attention :",
	},
	NotificationDistanceGlanceQuips: []string{
		"L’œil humain n’a pas été conçu pour fixer un écran éternellement.",
		"Cligner des yeux est gratuit. Faites-en un usage généreux.",
		"Même le plus beau des monocles doit être poli de temps à autre.",
		"La vision de loin : un luxe que vos yeux méritent amplement.",
		"Vos muscles de mise au point ont eux aussi besoin de s’étirer.",
		"Cornealius vous le rappelle : les écrans passent, la vue devrait rester.",
		"Un bref regard au loin ne coûte qu’un instant.",
		"Vos yeux vous ont bien servi. Rendez-leur la pareille.",
		"Les duels de regards avec votre écran sont déconseillés. Vous perdrez.",
		"Vos rétines ne s’achètent pas sur abonnement. Traitez-les en conséquence.",
		"L’horizon existe pour une raison. Consultez-le de temps à autre.",
		"Les pixels sont innombrables. Une bonne vue ne l’est pas.",
		"Cornealius le note : la lumière bleue ne remplace pas celle du soleil.",
		"Une fenêtre n’est pas qu’un élément décoratif. Regardez à travers.",
		"Vos muscles ciliaires aimeraient vous toucher un mot. Éloignez-vous de l’écran.",
		"Reposez-vous. La barre de défilement vous attendra à votre retour.",
		"La fatigue visuelle est bien réelle. La vue depuis votre fenêtre aussi.",
		"Même le commis le plus dévoué doit parfois lever les yeux de son grand livre.",
	},

	ParseErrorNoMinutes:     "Cornealius insiste pour avoir au moins une minute. Il a des principes.",
	ParseErrorInvalidMinute: "« %s » n’est pas une minute valide. Cornealius attend des nombres entiers entre 0 et 59.",
	ParseErrorCronEmpty:     "Cornealius requiert une expression CRON. Le silence n’est pas un emploi du temps.",
	ParseErrorCronInvalid:   "« %s » n’est pas une expression CRON Quartz valide.",

	LanguageDialogTitle:          "Langue - Cornealius Eyeworth",
	OptionsLanguageLabel:         "Langue (redémarrage requis) :",
	OptionsLanguageDefault:       "Langue du système",
	OptionsLanguageRestartNotice: "Le changement de langue prendra effet au prochain démarrage de Cornealius.",

	ScheduleDescriptionCron:   "CRON : %s",
	ScheduleDescriptionSimple: "Se présente aux minutes %s de chaque heure",
}
