package i18n

import (
	"reflect"
	"regexp"
	"slices"
	"testing"
)

func localeName(s *Strings) string {
	for _, l := range locales {
		if l.Strings == s {
			return l.Codes[0]
		}
	}
	return "unknown"
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
		t.Run(tt.in, func(t *testing.T) {
			if got := NormaliseLocale(tt.in); got != tt.want {
				t.Errorf("NormaliseLocale(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
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
		{"fr-FR", &french},
		{"fr-CA", &french},
		{"el-GR", &greek},

		// LANG on Linux.
		{"pl_PL", &polish},
		{"de_AT", &german},
		{"fr_BE", &french},
		{"el_CY", &greek},

		// Hand-edited config.json.
		{"PL", &polish},
		{"  De  ", &german},
		{"FR", &french},
		{"El", &greek},

		{"pl", &polish},
		{"de", &german},
		{"en", &english},
		{"fr", &french},
		{"el", &greek},
		{"de-li", &german},
		{"de-lu", &german},

		{"es", &english},
		{"gr", &english}, // country code, not the language code
		{"", &english},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			if got := localeFor(locales, tt.raw); got != tt.want {
				t.Errorf("localeFor(%q) = %s, want %s",
					tt.raw, localeName(got), localeName(tt.want))
			}
		})
	}
}

func TestSetLanguageExplicit(t *testing.T) {
	t.Cleanup(func() { Active = &english })

	for _, code := range []string{"pl-PL", "PL", "pl_PL"} {
		Active = &english
		SetLanguage(&code)
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

func TestLocaleFor_Disabled(t *testing.T) {
	registry := Locales()
	fr := slices.IndexFunc(registry, func(l Locale) bool { return l.Strings == &french })
	if fr < 0 {
		t.Fatal("french is not registered in locales")
	}
	registry[fr].Enabled = false

	for _, code := range []string{"fr", "fr-FR", "fr_CA"} {
		if got := localeFor(registry, code); got != &english {
			t.Errorf("localeFor(%q) with french disabled = %s, want en", code, localeName(got))
		}
	}

	registry[fr].Enabled = true
	if got := localeFor(registry, "fr-FR"); got != &french {
		t.Errorf("localeFor(\"fr-FR\") with french enabled = %s, want fr", localeName(got))
	}
}

func TestAvailableLanguages_ExcludesDisabled(t *testing.T) {
	registry := Locales()
	fr := slices.IndexFunc(registry, func(l Locale) bool { return l.Strings == &french })
	if fr < 0 {
		t.Fatal("french is not registered in locales")
	}
	registry[fr].Enabled = false

	for _, opt := range languageOptions(registry) {
		if opt.Code == "fr" {
			t.Errorf("languageOptions() offers disabled french: %+v", opt)
		}
	}
}

// Locales hands out a copy, so a caller flipping Enabled cannot reach the
// registry the app resolves against.
func TestLocalesReturnsACopy(t *testing.T) {
	registry := Locales()
	for i := range registry {
		registry[i].Enabled = false
		registry[i].Codes[0] = "xx"
	}
	if got := localeFor(locales, "fr"); got != &french {
		t.Errorf("localeFor(\"fr\") after mutating Locales() = %s, want fr", localeName(got))
	}
}

func TestAvailableLanguages_SystemDefaultFirst(t *testing.T) {
	options := AvailableLanguages()
	if len(options) == 0 {
		t.Fatal("AvailableLanguages() returned no entries, want the system-default entry first")
	}
	if options[0] != (LanguageOption{"", ""}) {
		t.Fatalf("AvailableLanguages()[0] = %+v, want the system-default entry", options[0])
	}
	for _, opt := range options[1:] {
		if opt.Code == "" || opt.DisplayName == "" {
			t.Errorf("AvailableLanguages() has a blank entry after the default: %+v", opt)
		}
	}
}

func TestLocaleCodesAreNormalisedAndUnique(t *testing.T) {
	owner := map[string]string{}
	for _, l := range locales {
		if len(l.Codes) == 0 {
			t.Errorf("%s has no codes", l.DisplayName)
			continue
		}
		for _, c := range l.Codes {
			if c != NormaliseLocale(c) {
				t.Errorf("%s: code %q is not normalised, localeFor would never match it", l.DisplayName, c)
			}
			if prev, dup := owner[c]; dup {
				t.Errorf("code %q is claimed by both %s and %s", c, prev, l.DisplayName)
			}
			owner[c] = l.DisplayName
		}
	}
}

var formatVerb = regexp.MustCompile(`%[-+ #0]*[0-9.]*[a-zA-Z%]`)

// A struct literal that omits a field compiles fine and shows the user an
// empty label, so every locale is checked field by field against English.
func TestLocalesAreComplete(t *testing.T) {
	en := reflect.ValueOf(english)
	typ := en.Type()

	for _, l := range locales {
		v := reflect.ValueOf(*l.Strings)
		t.Run(l.Codes[0], func(t *testing.T) {
			for i := range typ.NumField() {
				name := typ.Field(i).Name
				switch f := v.Field(i); f.Kind() {
				case reflect.String:
					if f.String() == "" {
						t.Errorf("%s is empty", name)
						continue
					}
					want := formatVerb.FindAllString(en.Field(i).String(), -1)
					if got := formatVerb.FindAllString(f.String(), -1); !slices.Equal(got, want) {
						t.Errorf("%s verbs = %v, want %v: %q", name, got, want, f.String())
					}
				case reflect.Slice:
					if f.Len() == 0 {
						t.Errorf("%s is empty", name)
					}
					for j := range f.Len() {
						if f.Index(j).String() == "" {
							t.Errorf("%s[%d] is empty", name, j)
						}
					}
				default:
					t.Errorf("%s: %s fields are not checked by this test", name, f.Kind())
				}
			}
		})
	}
}
