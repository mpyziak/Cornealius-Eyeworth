package scheduling

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
)

// Describe renders cronExpr as user-facing text: the minute list for a
// simple "0 m1,m2 * * * ?" schedule, or the raw expression otherwise.
func Describe(cronExpr string) string {
	parts := strings.Fields(cronExpr)
	if len(parts) == 6 &&
		parts[0] == "0" &&
		parts[2] == "*" && parts[3] == "*" && parts[4] == "*" && parts[5] == "?" {
		allValid := true
		for _, t := range strings.Split(parts[1], ",") {
			if m, err := strconv.Atoi(t); err != nil || m < 0 || m > 59 {
				allValid = false
				break
			}
		}
		if allValid {
			human := strings.ReplaceAll(parts[1], ",", ", ")
			return fmt.Sprintf(i18n.Active.ScheduleDescriptionSimple, human)
		}
	}
	return fmt.Sprintf(i18n.Active.ScheduleDescriptionCron, cronExpr)
}
