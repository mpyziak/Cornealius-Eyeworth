package scheduling

import "testing"

func TestParseCron(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     string
		wantValid bool
		wantExpr  string
	}{
		{"valid quartz", "0 0 9 ? * MON-FRI", true, "0 0 9 ? * MON-FRI"},
		{"valid with question mark", "0 15 10 ? * *", true, "0 15 10 ? * *"},
		{"valid with whitespace", "  0 30 12 ? * *  ", true, "0 30 12 ? * *"},
		{"empty input", "", false, ""},
		{"invalid cron", "not-a-cron", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			expr, err := ParseCron(tt.input)
			if (err == nil) != tt.wantValid {
				t.Fatalf("ParseCron(%q) error = %v, wantValid %v", tt.input, err, tt.wantValid)
			}
			if tt.wantValid && expr != tt.wantExpr {
				t.Fatalf("ParseCron(%q) = %q, want %q", tt.input, expr, tt.wantExpr)
			}
		})
	}
}

func TestParseMinutes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     string
		wantValid bool
		wantExpr  string
	}{
		{"single minute", "20", true, "0 20 * * * ?"},
		{"multiple minutes", "20,40,55", true, "0 20,40,55 * * * ?"},
		{"duplicates removed", "20, 40,20,55", true, "0 20,40,55 * * * ?"},
		{"whitespace separators", "20 40 55", true, "0 20,40,55 * * * ?"},
		{"empty input", "", false, ""},
		{"invalid minute", "70", false, ""},
		{"invalid token", "abc", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			expr, err := ParseMinutes(tt.input)
			if (err == nil) != tt.wantValid {
				t.Fatalf("ParseMinutes(%q) error = %v, wantValid %v", tt.input, err, tt.wantValid)
			}
			if tt.wantValid && expr != tt.wantExpr {
				t.Fatalf("ParseMinutes(%q) = %q, want %q", tt.input, expr, tt.wantExpr)
			}
		})
	}
}

func TestSimpleMinutes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		input  string
		want   string
		wantOk bool
	}{
		{"valid simple expression", "0 20,40,55 * * * ?", "20, 40, 55", true},
		{"valid single minute", "0 5 * * * ?", "5", true},
		{"invalid seconds field", "1 20,40,55 * * * ?", "", false},
		{"invalid day-of-month field", "0 20,40,55 1 * * * ?", "", false},
		{"invalid wildcard field", "0 20,40,55 * * * *", "", false},
		{"invalid minute token", "0 20,abc,55 * * * ?", "", false},
		{"wrong field count", "0 20,40,55 * * ?", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := SimpleMinutes(tt.input)
			if ok != tt.wantOk || got != tt.want {
				t.Fatalf("SimpleMinutes(%q) = (%q, %v), want (%q, %v)", tt.input, got, ok, tt.want, tt.wantOk)
			}
		})
	}
}
