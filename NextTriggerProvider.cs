namespace CornealiusEyeworth;

class NextTriggerProvider(Config config) : INextTriggerProvider
{
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
