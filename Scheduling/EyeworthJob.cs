using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Notifications;
using Quartz;

namespace CornealiusEyeworth.Scheduling;

/// <summary>
/// The Quartz.NET job that fires at each scheduled minute and dispatches
/// a reminder notification. Marked <see cref="DisallowConcurrentExecutionAttribute"/>
/// to prevent overlapping executions.
/// </summary>
[DisallowConcurrentExecution]
internal class EyeworthJob(INotificationService notificationService, Action? onFired = null) : IJob
{
    /// <summary>The stable key used to register this job with the Quartz scheduler.</summary>
    public static readonly JobKey Key = new(nameof(EyeworthJob));

    /// <inheritdoc/>
    public Task Execute(IJobExecutionContext context)
    {
        notificationService.SendReminderNotification();
        onFired?.Invoke();
        return Task.CompletedTask;
    }
}
