package scheduling

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/robfig/cron/v3"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

// robfig/cron has no '?'.
func quartzToRobfig(expr string) string {
	return strings.ReplaceAll(expr, "?", "*")
}

// cronParser is seconds-first, 6 fields.
var cronParser = cron.NewParser(
	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

// ParseCron validates input as a Quartz CRON expression.
func ParseCron(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New(i18n.Active.ParseErrorCronEmpty)
	}
	if _, err := cronParser.Parse(quartzToRobfig(trimmed)); err != nil {
		return "", fmt.Errorf(i18n.Active.ParseErrorCronInvalid, trimmed)
	}
	return trimmed, nil
}

// ParseMinutes turns a comma/space-separated minute list, e.g. "20,40,55",
// into the equivalent Quartz CRON expression "0 20,40,55 * * * ?".
func ParseMinutes(input string) (string, error) {
	parts := strings.FieldsFunc(strings.TrimSpace(input), func(r rune) bool {
		return r == ',' || r == ' '
	})

	if len(parts) == 0 {
		return "", errors.New(i18n.Active.ParseErrorNoMinutes)
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
			return "", fmt.Errorf(i18n.Active.ParseErrorInvalidMinute, p)
		}
		if !seen[m] {
			seen[m] = true
			minutes = append(minutes, m)
		}
	}

	if len(minutes) == 0 {
		return "", errors.New(i18n.Active.ParseErrorNoMinutes)
	}

	slices.Sort(minutes)

	strs := make([]string, len(minutes))
	for i, m := range minutes {
		strs[i] = strconv.Itoa(m)
	}
	return "0 " + strings.Join(strs, ",") + " * * * ?", nil
}

// SimpleMinutes is the inverse of ParseMinutes: it returns the minute list
// an expression of the shape "0 m1,m2 * * * ?" was built from, and whether
// cronExpr was that shape at all.
func SimpleMinutes(cronExpr string) (minutes string, ok bool) {
	parts := strings.Fields(cronExpr)
	if len(parts) != 6 {
		return "", false
	}
	if parts[0] != "0" || parts[2] != "*" || parts[3] != "*" || parts[4] != "*" || parts[5] != "?" {
		return "", false
	}
	for t := range strings.SplitSeq(parts[1], ",") {
		if m, err := strconv.Atoi(t); err != nil || m < 0 || m > 59 {
			return "", false
		}
	}
	return strings.ReplaceAll(parts[1], ",", ", "), true
}
