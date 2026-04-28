using CornealiusEyeworth.Configuration;

namespace CornealiusEyeworth.Scheduling;

/// <summary>
/// Computes the next wall-clock time a reminder will fire, given the
/// minutes-of-the-hour defined in <see cref="Config"/>.
/// </summary>
internal class NextTriggerProvider(Config config) : INextTriggerProvider
{
    /// <inheritdoc/>
    public DateTime GetNext()
    {
        var now = DateTime.Now;
        return config.MinutesOfHour
            .Select(m => new DateTime(now.Year, now.Month, now.Day, now.Hour, m, 0))
            .Select(t => t <= now ? t.AddHours(1) : t)
            .OrderBy(t => t)
            .First();
    }
}
