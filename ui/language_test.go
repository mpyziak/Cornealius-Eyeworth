package ui

import (
	"testing"

	"github.com/mpyziak/cornealius-eyeworth/config"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

func strPtr(s string) *string { return &s }

// langStr renders a *string for failure messages without printing its
// pointer address.
func langStr(s *string) string {
	if s == nil {
		return "nil"
	}
	return "\"" + *s + "\""
}

func TestUpdatedConfigForLanguage(t *testing.T) {
	tests := []struct {
		name         string
		existingLang *string
		newLang      *string
	}{
		{name: "sets a new language", existingLang: nil, newLang: strPtr("fr")},
		{name: "clears to system default", existingLang: strPtr("es"), newLang: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing := &config.Config{
				CronExpression:        "0 20,40,55 * * * ?",
				StandUpCronExpression: "0 55 * * * ?",
				Language:              tt.existingLang,
			}

			updated := updatedConfigForLanguage(existing, tt.newLang)

			if updated.CronExpression != existing.CronExpression {
				t.Errorf("updatedConfigForLanguage(...).CronExpression = %q, want %q", updated.CronExpression, existing.CronExpression)
			}
			if updated.StandUpCronExpression != existing.StandUpCronExpression {
				t.Errorf("updatedConfigForLanguage(...).StandUpCronExpression = %q, want %q", updated.StandUpCronExpression, existing.StandUpCronExpression)
			}
			gotNil, wantNil := updated.Language == nil, tt.newLang == nil
			if gotNil != wantNil || (!gotNil && *updated.Language != *tt.newLang) {
				t.Errorf("updatedConfigForLanguage(...).Language = %s, want %s", langStr(updated.Language), langStr(tt.newLang))
			}
		})
	}
}

func TestSelectedLanguageIndex(t *testing.T) {
	options := []i18n.LanguageOption{
		{DisplayName: "", Code: ""}, // system default
		{DisplayName: "English", Code: "en"},
		{DisplayName: "Deutsch", Code: "de"},
		{DisplayName: "Polski", Code: "pl"},
	}

	tests := []struct {
		name string
		lang *string
		want int
	}{
		{name: "nil selects system default", lang: nil, want: 0},
		{name: "exact match", lang: strPtr("pl"), want: 3},
		{name: "upper case", lang: strPtr("PL"), want: 3},
		{name: "mixed case", lang: strPtr("De"), want: 2},
		{name: "surrounding space", lang: strPtr("  en  "), want: 1},
		{name: "unknown code falls back to default", lang: strPtr("fr"), want: 0},
		{name: "empty string does not match the default entry", lang: strPtr(""), want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := selectedLanguageIndex(options, tt.lang); got != tt.want {
				t.Fatalf("selectedLanguageIndex(%s) = %d, want %d", langStr(tt.lang), got, tt.want)
			}
		})
	}
}

func TestLanguageChanged(t *testing.T) {
	tests := []struct {
		name        string
		existing    *string
		newLang     *string
		wantChanged bool
	}{
		{name: "both nil", existing: nil, newLang: nil, wantChanged: false},
		{name: "set from nil", existing: nil, newLang: strPtr("fr"), wantChanged: true},
		{name: "clear to nil", existing: strPtr("fr"), newLang: nil, wantChanged: true},
		{name: "same value", existing: strPtr("fr"), newLang: strPtr("fr"), wantChanged: false},
		{name: "different value", existing: strPtr("fr"), newLang: strPtr("de"), wantChanged: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := languageChanged(tt.existing, tt.newLang)
			if got != tt.wantChanged {
				t.Fatalf("languageChanged(%s, %s) = %v, want %v", langStr(tt.existing), langStr(tt.newLang), got, tt.wantChanged)
			}
		})
	}
}
