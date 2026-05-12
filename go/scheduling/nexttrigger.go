package scheduling

import (
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var nextTriggerParser = cron.NewParser(
	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

// NextTrigger returns the next wall-clock time the given Quartz CRON expression
// will fire. Falls back to one hour from now if the expression cannot be parsed.
func NextTrigger(cronExpr string) time.Time {
	normalized := strings.ReplaceAll(cronExpr, "?", "*")
	schedule, err := nextTriggerParser.Parse(normalized)
	if err != nil {
		return time.Now().Add(time.Hour)
	}
	return schedule.Next(time.Now())
}
