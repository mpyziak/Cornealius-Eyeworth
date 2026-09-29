// Writes every locale's strings as one CSV - a row per key, a column per
// language - for Google Sheets' File › Import.
//
//	go run ./tools/i18n-export -exclude el -o dist/i18n-strings.csv
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() (err error) {
	out := flag.String("o", "", "output file (default stdout)")
	exclude := flag.String("exclude", "", "comma-separated locale codes to leave out, e.g. el")
	flag.Parse()

	locales, err := selectLocales(i18n.Locales(), *exclude)
	if err != nil {
		return err
	}

	var w io.Writer = os.Stdout
	if *out != "" {
		f, ferr := os.Create(*out)
		if ferr != nil {
			return ferr
		}
		defer func() {
			if cerr := f.Close(); err == nil {
				err = cerr
			}
		}()
		w = f
	}

	return writeCSV(w, locales)
}

// Disabled locales are exported too; the flag only governs the running app.
// An unknown code is an error so a typo ("gr") cannot silently keep a column.
func selectLocales(all []i18n.Locale, exclude string) ([]i18n.Locale, error) {
	skip := map[string]bool{}
	for _, c := range strings.Split(exclude, ",") {
		if c = i18n.NormaliseLocale(c); c != "" {
			skip[c] = true
		}
	}

	var kept []i18n.Locale
	for _, l := range all {
		if skip[l.Codes[0]] {
			delete(skip, l.Codes[0])
			continue
		}
		kept = append(kept, l)
	}
	if len(skip) > 0 {
		unknown := make([]string, 0, len(skip))
		for c := range skip {
			unknown = append(unknown, c)
		}
		slices.Sort(unknown)
		return nil, fmt.Errorf("-exclude: no locale with code(s) %q", unknown)
	}
	return kept, nil
}

// A []string field gets a row per index, blank where a locale's pool is
// shorter. The pools are picked from at random, so row N lines up across
// columns only as far as each translation kept the English order.
func writeCSV(w io.Writer, locales []i18n.Locale) error {
	cw := csv.NewWriter(w)

	header := []string{"key"}
	tables := make([]reflect.Value, len(locales))
	for i, l := range locales {
		header = append(header, l.Codes[0])
		tables[i] = reflect.ValueOf(*l.Strings)
	}
	if err := cw.Write(header); err != nil {
		return err
	}

	typ := reflect.TypeOf(i18n.Strings{})
	for f := range typ.NumField() {
		field := typ.Field(f)
		switch field.Type.Kind() {
		case reflect.String:
			row := []string{field.Name}
			for _, t := range tables {
				row = append(row, t.Field(f).String())
			}
			if err := cw.Write(row); err != nil {
				return err
			}

		case reflect.Slice:
			n := 0
			for _, t := range tables {
				n = max(n, t.Field(f).Len())
			}
			for j := range n {
				row := []string{fmt.Sprintf("%s[%d]", field.Name, j)}
				for _, t := range tables {
					cell := ""
					if pool := t.Field(f); j < pool.Len() {
						cell = pool.Index(j).String()
					}
					row = append(row, cell)
				}
				if err := cw.Write(row); err != nil {
					return err
				}
			}

		default:
			return fmt.Errorf("i18n.Strings.%s: %s fields are not exportable", field.Name, field.Type)
		}
	}

	cw.Flush()
	return cw.Error()
}
