using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Notifications;
using Quartz;
using Quartz.Impl;

namespace CornealiusEyeworth.Scheduling;

internal class SchedulerService(Config config, INotificationService notificationService, Action? onJobFired = null)
{
    public async Task RunAsync(CancellationToken cancellationToken)
    {
        var factory = new StdSchedulerFactory();
        var scheduler = await factory.GetScheduler(cancellationToken);

        scheduler.JobFactory = new EyeworthJobFactory(notificationService, config, onJobFired);

        var job = JobBuilder.Create<EyeworthJob>()
            .WithIdentity(EyeworthJob.Key)
            .Build();

        var trigger = TriggerBuilder.Create()
            .WithIdentity("EyeworthTrigger")
            .WithCronSchedule(config.CronExpression)
            .Build();

        await scheduler.ScheduleJob(job, trigger, cancellationToken);
        await scheduler.Start(cancellationToken);

        try { await Task.Delay(Timeout.Infinite, cancellationToken); }
        catch (OperationCanceledException) { }

        await scheduler.Shutdown(waitForJobsToComplete: false);
    }
}