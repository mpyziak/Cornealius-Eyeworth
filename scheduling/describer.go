package scheduling

import (
	"fmt"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

// Describe renders cronExpr as user-facing text: the minute list for a
// simple "0 m1,m2 * * * ?" schedule, or the raw expression otherwise.
func Describe(cronExpr string) string {
	if minutes, ok := SimpleMinutes(cronExpr); ok {
		return fmt.Sprintf(i18n.Active.ScheduleDescriptionSimple, minutes)
	}
	return fmt.Sprintf(i18n.Active.ScheduleDescriptionCron, cronExpr)
}
