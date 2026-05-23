package parsing

import "testing"

func TestParseCron(t *testing.T) {
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
			result := ParseCron(tt.input)
			if result.Valid != tt.wantValid {
				t.Fatalf("ParseCron(%q).Valid = %v, want %v", tt.input, result.Valid, tt.wantValid)
			}
			if tt.wantValid && result.Expression != tt.wantExpr {
				t.Fatalf("ParseCron(%q).Expression = %q, want %q", tt.input, result.Expression, tt.wantExpr)
			}
			if !tt.wantValid && result.Err == "" {
				t.Fatalf("ParseCron(%q) returned no error for invalid input", tt.input)
			}
		})
	}
}

func TestParseMinutes(t *testing.T) {
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
			result := ParseMinutes(tt.input)
			if result.Valid != tt.wantValid {
				t.Fatalf("ParseMinutes(%q).Valid = %v, want %v", tt.input, result.Valid, tt.wantValid)
			}
			if tt.wantValid && result.Expression != tt.wantExpr {
				t.Fatalf("ParseMinutes(%q).Expression = %q, want %q", tt.input, result.Expression, tt.wantExpr)
			}
			if !tt.wantValid && result.Err == "" {
				t.Fatalf("ParseMinutes(%q) returned no error for invalid input", tt.input)
			}
		})
	}
}

func TestTryExtractSimpleMinutes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"valid simple expression", "0 20,40,55 * * * ?", "20, 40, 55"},
		{"valid single minute", "0 5 * * * ?", "5"},
		{"invalid seconds field", "1 20,40,55 * * * ?", ""},
		{"invalid day-of-month field", "0 20,40,55 1 * * * ?", ""},
		{"invalid wildcard field", "0 20,40,55 * * * *", ""},
		{"invalid minute token", "0 20,abc,55 * * * ?", ""},
		{"wrong field count", "0 20,40,55 * * ?", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TryExtractSimpleMinutes(tt.input)
			if got != tt.want {
				t.Fatalf("TryExtractSimpleMinutes(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
