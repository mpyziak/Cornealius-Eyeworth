package scheduling

import (
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var nextTriggerParser = cron.NewParser(
	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

// NextTrigger returns cronExpr's next fire time, falling back to one hour
// out on a bad expression.
func NextTrigger(cronExpr string) time.Time {
	normalized := strings.ReplaceAll(cronExpr, "?", "*")
	schedule, err := nextTriggerParser.Parse(normalized)
	if err != nil {
		return time.Now().Add(time.Hour)
	}
	return schedule.Next(time.Now())
}
