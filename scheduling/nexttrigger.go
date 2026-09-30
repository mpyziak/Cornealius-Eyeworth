package scheduling

import (
	"time"
)

// NextTrigger returns cronExpr's next fire time, falling back to one hour
// out on a bad expression.
func NextTrigger(cronExpr string) time.Time {
	schedule, err := cronParser.Parse(quartzToRobfig(cronExpr))
	if err != nil {
		return time.Now().Add(time.Hour)
	}
	return schedule.Next(time.Now())
}
