package ui

import (
	"testing"

	"github.com/mpyziak/cornealius-eyeworth/config"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

func strPtr(s string) *string { return &s }

func TestUpdatedConfigForLanguagePreservesCronExpressions(t *testing.T) {
	existing := &config.Config{
		CronExpression:        "0 20,40,55 * * * ?",
		StandUpCronExpression: "0 55 * * * ?",
		Language:              nil,
	}

	updated := updatedConfigForLanguage(existing, strPtr("fr"))

	if updated.CronExpression != existing.CronExpression {
		t.Fatalf("CronExpression changed: got %q, want %q", updated.CronExpression, existing.CronExpression)
	}
	if updated.StandUpCronExpression != existing.StandUpCronExpression {
		t.Fatalf("StandUpCronExpression changed: got %q, want %q", updated.StandUpCronExpression, existing.StandUpCronExpression)
	}
	if updated.Language == nil || *updated.Language != "fr" {
		t.Fatalf("Language not set correctly: got %v", updated.Language)
	}
}

func TestUpdatedConfigForLanguageClearingToDefaultPreservesCronExpressions(t *testing.T) {
	existing := &config.Config{
		CronExpression:        "0 20,40,55 * * * ?",
		StandUpCronExpression: "0 55 * * * ?",
		Language:              strPtr("es"),
	}

	updated := updatedConfigForLanguage(existing, nil)

	if updated.CronExpression != existing.CronExpression {
		t.Fatalf("CronExpression changed: got %q, want %q", updated.CronExpression, existing.CronExpression)
	}
	if updated.StandUpCronExpression != existing.StandUpCronExpression {
		t.Fatalf("StandUpCronExpression changed: got %q, want %q", updated.StandUpCronExpression, existing.StandUpCronExpression)
	}
	if updated.Language != nil {
		t.Fatalf("Language should be nil, got %v", updated.Language)
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
				t.Fatalf("selectedLanguageIndex(%v) = %d, want %d", tt.lang, got, tt.want)
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
				t.Fatalf("languageChanged() = %v, want %v", got, tt.wantChanged)
			}
		})
	}
}
