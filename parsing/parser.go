package parsing

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/robfig/cron/v3"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

type ParseResult struct {
	Valid      bool
	Expression string // set when Valid == true; canonical Quartz CRON expression
	Err        string // set when Valid == false
}

func ok(expr string) ParseResult  { return ParseResult{Valid: true, Expression: expr} }
func fail(msg string) ParseResult { return ParseResult{Valid: false, Err: msg} }

// robfig/cron has no '?'.
func quartzToRobfig(expr string) string {
	return strings.ReplaceAll(expr, "?", "*")
}

// Seconds-first, 6 fields.
var cronParser = cron.NewParser(
	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

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

// "20,40,55" -> "0 20,40,55 * * * ?".
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

// Inverse of ParseMinutes. Empty string if the expression is not that shape.
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
