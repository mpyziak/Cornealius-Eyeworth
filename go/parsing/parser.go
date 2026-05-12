// Package parsing validates and converts schedule inputs.
package parsing

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/robfig/cron/v3"
)

// ParseResult is the outcome of a parse operation.
type ParseResult struct {
	Valid      bool
	Expression string // set when Valid == true; canonical Quartz CRON expression
	Err        string // set when Valid == false
}

func ok(expr string) ParseResult  { return ParseResult{Valid: true, Expression: expr} }
func fail(msg string) ParseResult { return ParseResult{Valid: false, Err: msg} }

// quartzToRobfig replaces Quartz '?' wildcards with '*' for robfig/cron compatibility.
func quartzToRobfig(expr string) string {
	return strings.ReplaceAll(expr, "?", "*")
}

// cronParser is a robfig/cron parser that understands seconds-first 6-field expressions.
var cronParser = cron.NewParser(
	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

// ParseCron validates a Quartz-style 6-field CRON expression.
// '?' is accepted and treated as '*'.
func ParseCron(input string) ParseResult {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return fail(i18n.Active.ParseErrorCronEmpty)
	}
	if _, err := cronParser.Parse(quartzToRobfig(trimmed)); err != nil {
		return fail(fmt.Sprintf(i18n.Active.ParseErrorCronInvalid, trimmed))
	}
	return ok(trimmed)
}

// ParseMinutes converts a comma-separated list of minute values (e.g. "20,40,55")
// into a canonical Quartz CRON expression such as "0 20,40,55 * * * ?".
func ParseMinutes(input string) ParseResult {
	parts := strings.FieldsFunc(strings.TrimSpace(input), func(r rune) bool {
		return r == ',' || r == ' '
	})

	if len(parts) == 0 {
		return fail(i18n.Active.ParseErrorNoMinutes)
	}

	seen := map[int]bool{}
	var minutes []int
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		m, err := strconv.Atoi(p)
		if err != nil || m < 0 || m > 59 {
			return fail(fmt.Sprintf(i18n.Active.ParseErrorInvalidMinute, p))
		}
		if !seen[m] {
			seen[m] = true
			minutes = append(minutes, m)
		}
	}

	if len(minutes) == 0 {
		return fail(i18n.Active.ParseErrorNoMinutes)
	}

	sort.Ints(minutes)

	strs := make([]string, len(minutes))
	for i, m := range minutes {
		strs[i] = strconv.Itoa(m)
	}
	return ok("0 " + strings.Join(strs, ",") + " * * * ?")
}

// TryExtractSimpleMinutes returns the human-readable comma-separated minutes string
// from a simple "0 m1,m2 * * * ?" expression (e.g. "20, 40, 55"), or an empty
// string if the expression is not in that simple form.
func TryExtractSimpleMinutes(cronExpr string) string {
	parts := strings.Fields(cronExpr)
	if len(parts) != 6 {
		return ""
	}
	if parts[0] != "0" || parts[2] != "*" || parts[3] != "*" || parts[4] != "*" || parts[5] != "?" {
		return ""
	}
	for _, t := range strings.Split(parts[1], ",") {
		if m, err := strconv.Atoi(t); err != nil || m < 0 || m > 59 {
			return ""
		}
	}
	return strings.ReplaceAll(parts[1], ",", ", ")
}
