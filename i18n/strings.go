// Package i18n holds the translated UI strings. Call SetLanguage once at
// startup, then read Active.
//
// Each language's table lives in its own file (en.go, de.go, ...); the
// registry returned by Locales is the one place that decides which of them
// the app can use.
package i18n

import (
	"slices"
	"strings"

	"fyne.io/fyne/v2/lang"
)

// Strings is one language's complete set of UI text. Every locale must fill
// every field; TestLocalesAreComplete enforces it.
type Strings struct {
	AppName   string
	AppTitle  string
	GitHubURL string

	StatusServing      string
	NextTrigger        string // %s = HH:mm
	NextTriggerStandUp string // %s = HH:mm

	MenuOptions      string
	MenuLanguage     string
	MenuTriggerTimes string
	MenuHelp         string
	MenuAbout        string
	MenuHowToUse     string
	MenuGitHub       string
	MenuQuit         string

	TrayTooltipShowHide string

	ScheduleDialogTitle       string
	OptionsInstruction        string
	OptionsStandUpInstruction string
	ScheduleCronInstruction   string
	ScheduleStandardToggle    string
	ScheduleAdvancedToggle    string

	ButtonSave   string
	ButtonCancel string
	ButtonClose  string

	AboutDialogTitle string
	AboutDescription string
	AboutVersion     string

	HelpDialogTitle string
	HelpBody        string

	NotificationOnDuty                string
	NotificationMinimizedToTray       string
	NotificationDistanceGlanceHeaders []string
	NotificationMovementHeaders       []string
	NotificationMovementQuips         []string
	NotificationCombinedHeaders       []string
	NotificationDistanceGlanceQuips   []string

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

// Locale is one entry in the language registry.
type Locale struct {
	// Codes[0] is what the language dialog writes to config.json. The rest are
	// the region-qualified forms that also resolve here; the list is
	// enumerated on purpose, see doc/potential-fixes.md §5. All normalised.
	Codes       []string
	DisplayName string
	// A disabled locale is neither offered in the dialog nor resolved from
	// config.json or the OS; those users get English.
	Enabled bool
	Strings *Strings
}

// Unexported so Enabled stays a compile-time switch: nothing outside this
// file can flip it, reorder entries or swap a table at runtime.
var locales = []Locale{
	{
		Codes:       []string{"en"},
		DisplayName: "English",
		Enabled:     true,
		Strings:     &english,
	},
	{
		Codes:       []string{"de", "de-de", "de-at", "de-ch", "de-li", "de-lu"},
		DisplayName: "Deutsch",
		Enabled:     true,
		Strings:     &german,
	},
	{
		Codes:       []string{"pl", "pl-pl"},
		DisplayName: "Polski",
		Enabled:     true,
		Strings:     &polish,
	},
	{
		Codes:       []string{"fr", "fr-fr", "fr-be", "fr-ch", "fr-ca", "fr-lu", "fr-mc"},
		DisplayName: "Français",
		Enabled:     true,
		Strings:     &french,
	},
	{
		Codes:       []string{"el", "el-gr", "el-cy"},
		DisplayName: "Ελληνικά",
		Enabled:     true,
		Strings:     &greek,
	},
}

// Locales returns a copy of the registry, disabled locales included, in
// dialog order. The Strings tables it points to are shared and must not be
// modified.
func Locales() []Locale {
	out := slices.Clone(locales)
	for i := range out {
		out[i].Codes = slices.Clone(out[i].Codes)
	}
	return out
}

// Active is the table the UI reads from. SetLanguage assigns it.
var Active = &english

// NormaliseLocale lower-cases code, trims it and turns "_" into "-", so
// "pl_PL" and " pl-PL " both become "pl-pl".
func NormaliseLocale(code string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(code)), "_", "-")
}

// SetLanguage points Active at the table for language, or at the OS locale's
// table when language is nil. Anything unrecognised or disabled gets English.
func SetLanguage(language *string) {
	var raw string
	if language == nil {
		raw = lang.SystemLocale().LanguageString()
	} else {
		raw = *language
	}
	Active = localeFor(locales, raw)
}

// Split out of SetLanguage so it is testable without the host's locale, and
// takes the registry so tests can pass a modified copy.
// Normalise first: Windows hands over "pl-PL", upper-case region and all.
func localeFor(registry []Locale, raw string) *Strings {
	code := NormaliseLocale(raw)
	for _, l := range registry {
		if l.Enabled && slices.Contains(l.Codes, code) {
			return l.Strings
		}
	}
	return &english
}

// LanguageOption is one row of the language dialog.
type LanguageOption struct {
	DisplayName string
	Code        string // empty string = system default
}

// AvailableLanguages lists the dialog's choices: the system default first,
// then every enabled locale.
func AvailableLanguages() []LanguageOption {
	return languageOptions(locales)
}

func languageOptions(registry []Locale) []LanguageOption {
	options := make([]LanguageOption, 0, len(registry)+1)
	// Display name filled at runtime from Active.OptionsLanguageDefault.
	options = append(options, LanguageOption{DisplayName: "", Code: ""})
	for _, l := range registry {
		if l.Enabled {
			options = append(options, LanguageOption{DisplayName: l.DisplayName, Code: l.Codes[0]})
		}
	}
	return options
}
