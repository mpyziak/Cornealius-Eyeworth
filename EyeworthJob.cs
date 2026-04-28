using Quartz;

namespace CornealiusEyeworth;

[DisallowConcurrentExecution]
class EyeworthJob(NotificationService notificationService, Config config) : IJob
{
    public static readonly JobKey Key = new(nameof(EyeworthJob));

    public Task Execute(IJobExecutionContext context)
    {
        notificationService.SendReminderNotification(config.NotificationMessage);
        return Task.CompletedTask;
    }
}
