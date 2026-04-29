using CornealiusEyeworth.Configuration;
using Quartz;

namespace CornealiusEyeworth.Scheduling;

/// <summary>
/// Computes the next wall-clock time a reminder will fire from the config CRON expression.
/// </summary>
internal class NextTriggerProvider(Config config) : INextTriggerProvider
{
    public DateTime GetNext()
    {
        var cron = new CronExpression(config.CronExpression);
        var next = cron.GetNextValidTimeAfter(DateTimeOffset.Now);
        return next.HasValue ? next.Value.LocalDateTime : DateTime.Now.AddHours(1);
    }
}