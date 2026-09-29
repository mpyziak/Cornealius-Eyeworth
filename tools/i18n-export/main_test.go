package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"testing"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

func export(t *testing.T, exclude string) [][]string {
	t.Helper()
	locales, err := selectLocales(i18n.Locales(), exclude)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := writeCSV(&buf, locales); err != nil {
		t.Fatal(err)
	}
	rows, err := csv.NewReader(&buf).ReadAll()
	if err != nil {
		t.Fatalf("output is not valid CSV: %v", err)
	}
	return rows
}

func find(rows [][]string, key string) []string {
	for _, r := range rows {
		if r[0] == key {
			return r
		}
	}
	return nil
}

func TestExcludeDropsTheColumn(t *testing.T) {
	rows := export(t, "EL")
	for _, h := range rows[0] {
		if h == "el" {
			t.Fatalf("header still has el: %v", rows[0])
		}
	}
	if len(rows[0]) != len(i18n.Locales()) { // key column + all but one locale
		t.Fatalf("header = %v, want key + %d locales", rows[0], len(i18n.Locales())-1)
	}
}

func TestUnknownExcludeCodeIsAnError(t *testing.T) {
	if _, err := selectLocales(i18n.Locales(), "el,gr"); err == nil {
		t.Fatal(`selectLocales(..., "el,gr") accepted "gr"`)
	}
}

// Multi-line HelpBody and the U+00A0 in French must survive a CSV round trip.
func TestCellsRoundTrip(t *testing.T) {
	rows := export(t, "")
	col := map[string]int{}
	for i, h := range rows[0] {
		col[h] = i
	}

	var english, french *i18n.Strings
	for _, l := range i18n.Locales() {
		switch l.Codes[0] {
		case "en":
			english = l.Strings
		case "fr":
			french = l.Strings
		}
	}

	help := find(rows, "HelpBody")
	if help == nil {
		t.Fatal("HelpBody row missing from export")
	}
	if got := help[col["en"]]; got != english.HelpBody {
		t.Errorf("HelpBody[en] = %q, want %q", got, english.HelpBody)
	}

	next := find(rows, "NextTrigger")
	if next == nil {
		t.Fatal("NextTrigger row missing from export")
	}
	if got := next[col["fr"]]; got != french.NextTrigger {
		t.Errorf("NextTrigger[fr] = %q, want %q", got, french.NextTrigger)
	}
}

func TestShorterPoolsArePaddedWithBlanks(t *testing.T) {
	rows := export(t, "")
	width := len(rows[0])
	for _, r := range rows {
		if len(r) != width {
			t.Fatalf("row %q has %d cells, header has %d", r[0], len(r), width)
		}
	}

	longest := 0
	for _, l := range i18n.Locales() {
		longest = max(longest, len(l.Strings.NotificationDistanceGlanceHeaders))
	}
	for j := 0; j < longest; j++ {
		key := fmt.Sprintf("NotificationDistanceGlanceHeaders[%d]", j)
		if find(rows, key) == nil {
			t.Errorf("missing row %s", key)
		}
	}
}
