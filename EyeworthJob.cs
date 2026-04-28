using Quartz;

namespace CornealiusEyeworth;

[DisallowConcurrentExecution]
class EyeworthJob(NotificationService notificationService, Config config, Action? onFired = null) : IJob
{
    public static readonly JobKey Key = new(nameof(EyeworthJob));

    public Task Execute(IJobExecutionContext context)
    {
        notificationService.SendReminderNotification(config.NotificationMessage);
        onFired?.Invoke();
        return Task.CompletedTask;
    }
}
