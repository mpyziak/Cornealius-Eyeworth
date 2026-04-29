using CornealiusEyeworth.Localization;

namespace CornealiusEyeworth.Scheduling;

/// <summary>
/// Converts a Quartz CRON expression into a human-readable description.
/// Simple "0 m1,m2 * * * ?" patterns are expanded to a localised sentence;
/// anything more complex falls back to the raw expression.
/// </summary>
internal static class CronDescriber
{
    public static string Describe(string cron)
    {
        var parts = cron.Trim().Split(' ');
        if (parts.Length == 6
            && parts[0] == "0"
            && parts[2] == "*" && parts[3] == "*" && parts[4] == "*" && parts[5] == "?"
            && parts[1].Split(',').All(t => int.TryParse(t, out var m) && m >= 0 && m <= 59))
        {
            return string.Format(Strings.ScheduleDescriptionSimple, parts[1].Replace(",", ", "));
        }

        return string.Format(Strings.ScheduleDescriptionCron, cron);
    }
}