package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

var update = flag.Bool("update", false, "update golden files")

func exportCSV(t *testing.T, exclude string) []byte {
	t.Helper()
	locales, err := selectLocales(i18n.Locales(), exclude)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := writeCSV(&buf, locales); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func export(t *testing.T, exclude string) [][]string {
	t.Helper()
	rows, err := csv.NewReader(bytes.NewReader(exportCSV(t, exclude))).ReadAll()
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

func TestSelectLocales_Exclude(t *testing.T) {
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

// The full exported shape - column order, row order, padding, quoting, the
// multi-line HelpBody, the U+00A0 in French - as one comparison. Replaces
// hand-rolled per-cell lookups, which could look up a missing column as
// index 0 (the key column, not a real miss) and dereference a not-found
// locale's *Strings as nil.
//
// go test ./tools/i18n-export/... -run TestWriteCSV_Golden -update
// regenerates testdata/export.golden after a deliberate strings/locale change.
func TestWriteCSV_Golden(t *testing.T) {
	got := exportCSV(t, "")

	golden := filepath.Join("testdata", "export.golden")
	if *update {
		if err := os.WriteFile(golden, got, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("writeCSV output mismatch (-want +got):\n%s", diff)
	}
}

func TestWriteCSV_ShortPoolsPadded(t *testing.T) {
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
	for j := range longest {
		key := fmt.Sprintf("NotificationDistanceGlanceHeaders[%d]", j)
		if find(rows, key) == nil {
			t.Errorf("missing row %s", key)
		}
	}
}
