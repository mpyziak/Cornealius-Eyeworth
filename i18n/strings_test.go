package i18n

import "testing"

func localeName(s *Strings) string {
	switch s {
	case &english:
		return "english"
	case &german:
		return "german"
	case &polish:
		return "polish"
	default:
		return "unknown"
	}
}

func TestNormaliseLocale(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"pl-PL", "pl-pl"},
		{"pl_PL", "pl-pl"},
		{"de-DE", "de-de"},
		{"  de-AT  ", "de-at"},
		{"PL", "pl"},
		{"en", "en"},
		{"", ""},
	}

	for _, tt := range tests {
		if got := NormaliseLocale(tt.in); got != tt.want {
			t.Errorf("NormaliseLocale(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestLocaleFor(t *testing.T) {
	tests := []struct {
		raw  string
		want *Strings
	}{
		// Verbatim from Fyne's SystemLocale() on Polish and German Windows.
		{"pl-PL", &polish},
		{"de-DE", &german},
		{"de-AT", &german},
		{"de-CH", &german},

		// LANG on Linux.
		{"pl_PL", &polish},
		{"de_AT", &german},

		// Hand-edited config.json.
		{"PL", &polish},
		{"  De  ", &german},

		{"pl", &polish},
		{"de", &german},
		{"en", &english},
		{"de-li", &german},
		{"de-lu", &german},

		{"fr-FR", &english},
		{"es", &english},
		{"", &english},
	}

	for _, tt := range tests {
		if got := localeFor(tt.raw); got != tt.want {
			t.Errorf("localeFor(%q) = %s, want %s",
				tt.raw, localeName(got), localeName(tt.want))
		}
	}
}

func TestSetLanguageExplicit(t *testing.T) {
	t.Cleanup(func() { Active = &english })

	for _, code := range []string{"pl-PL", "PL", "pl_PL"} {
		Active = &english
		c := code
		SetLanguage(&c)
		if Active != &polish {
			t.Errorf("SetLanguage(%q) = %s, want polish", code, localeName(Active))
		}
	}
}

// Host-dependent, so all this can assert is that Active lands on some known
// locale rather than nil.
func TestSetLanguageSystemDefault(t *testing.T) {
	t.Cleanup(func() { Active = &english })

	SetLanguage(nil)

	if localeName(Active) == "unknown" {
		t.Fatalf("SetLanguage(nil) left Active pointing at an unknown locale table")
	}
	t.Logf("system locale resolved to: %s", localeName(Active))
}
