package notifications

import (
	"testing"
	"unicode/utf16"

	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/scheduling"
)

// Not len(s) and not len([]rune(s)) - the German and Polish accents differ.
func utf16Len(s string) int {
	return len(utf16.Encode([]rune(s)))
}

func longest(items []string) (text string, n int) {
	for _, s := range items {
		if l := utf16Len(s); l > n {
			text, n = s, l
		}
	}
	return text, n
}

// The quips are constants, so an overlong one is a build-time defect. This is
// why nothing truncates at runtime. Failures print the headroom.
// Disabled locales included: switching one on must not ship an overlong quip.
func TestNotificationsFitWin32BalloonLimits(t *testing.T) {
	for _, l := range i18n.Locales() {
		str := l.Strings

		t.Run(l.Codes[0], func(t *testing.T) {
			// Every set a title can come from.
			titleSets := map[string][]string{
				"NotificationCombinedHeaders":       str.NotificationCombinedHeaders,
				"NotificationMovementHeaders":       str.NotificationMovementHeaders,
				"NotificationDistanceGlanceHeaders": str.NotificationDistanceGlanceHeaders,
				"NotificationOnDuty":                {str.NotificationOnDuty},
				"AppName":                           {str.AppName},
			}
			for name, set := range titleSets {
				if len(set) == 0 {
					t.Errorf("%s is empty", name)
					continue
				}
				text, n := longest(set)
				if n > MaxBalloonTitleUTF16 {
					t.Errorf("%s: longest entry is %d UTF-16 units, limit %d (over by %d)\n  %q",
						name, n, MaxBalloonTitleUTF16, n-MaxBalloonTitleUTF16, text)
					continue
				}
				t.Logf("%-34s worst %3d / %d units (%d spare)",
					name, n, MaxBalloonTitleUTF16, MaxBalloonTitleUTF16-n)
			}

			// Worst case: both schedules in one Buffer window, longest quip each.
			standUp, _ := longest(str.NotificationMovementQuips)
			eye, _ := longest(str.NotificationDistanceGlanceQuips)
			if standUp == "" || eye == "" {
				t.Fatal("quip sets must not be empty")
			}

			worst := aggregatedContent([]scheduling.Reminder{
				{Category: scheduling.CategoryStandUp, Message: standUp},
				{Category: scheduling.CategoryEye, Message: eye},
			})

			if n := utf16Len(worst); n > MaxBalloonBodyUTF16 {
				t.Errorf("worst-case aggregated body is %d UTF-16 units, limit %d (over by %d)\n%s",
					n, MaxBalloonBodyUTF16, n-MaxBalloonBodyUTF16, worst)
			} else {
				t.Logf("%-34s worst %3d / %d units (%d spare)",
					"aggregated body (standup+eye)", n, MaxBalloonBodyUTF16, MaxBalloonBodyUTF16-n)
			}

			for name, s := range map[string]string{
				"NotificationMinimizedToTray": str.NotificationMinimizedToTray,
			} {
				if n := utf16Len(s); n > MaxBalloonBodyUTF16 {
					t.Errorf("%s is %d UTF-16 units, limit %d", name, n, MaxBalloonBodyUTF16)
				}
			}
		})
	}
}
