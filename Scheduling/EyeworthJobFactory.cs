using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Notifications;
using Quartz;
using Quartz.Spi;

namespace CornealiusEyeworth.Scheduling;

/// <summary>
/// Quartz <see cref="IJobFactory"/> that constructs <see cref="EyeworthJob"/> instances
/// with their required dependencies injected manually, avoiding a DI container dependency.
/// </summary>
internal class EyeworthJobFactory(INotificationService notificationService, Config config, Action? onJobFired) : IJobFactory
{
    /// <inheritdoc/>
    public IJob NewJob(TriggerFiredBundle bundle, IScheduler scheduler) =>
        new EyeworthJob(notificationService, config, onJobFired);

    /// <inheritdoc/>
    public void ReturnJob(IJob job) { }
}
