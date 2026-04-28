using Quartz;
using Quartz.Impl;
using Quartz.Spi;

namespace CornealiusEyeworth;

class SchedulerService(Config config, NotificationService notificationService, Action? onJobFired = null)
{
    public async Task RunAsync(CancellationToken cancellationToken)
    {
        var factory = new StdSchedulerFactory();
        var scheduler = await factory.GetScheduler(cancellationToken);

        scheduler.JobFactory = new EyeworthJobFactory(notificationService, config, onJobFired);

        var job = JobBuilder.Create<EyeworthJob>()
            .WithIdentity(EyeworthJob.Key)
            .Build();

        var trigger = BuildCronTrigger(config.MinutesOfHour);

        await scheduler.ScheduleJob(job, trigger, cancellationToken);
        await scheduler.Start(cancellationToken);

        try { await Task.Delay(Timeout.Infinite, cancellationToken); }
        catch (OperationCanceledException) { }

        await scheduler.Shutdown(waitForJobsToComplete: false);
    }

    private static ITrigger BuildCronTrigger(int[] minutes)
    {
        var minuteList = string.Join(",", minutes);
        var cron = $"0 {minuteList} * * * ?";

        return TriggerBuilder.Create()
            .WithIdentity("EyeworthTrigger")
            .WithCronSchedule(cron)
            .Build();
    }
}

class EyeworthJobFactory(NotificationService notificationService, Config config, Action? onJobFired) : IJobFactory
{
    public IJob NewJob(TriggerFiredBundle bundle, IScheduler scheduler) =>
        new EyeworthJob(notificationService, config, onJobFired);

    public void ReturnJob(IJob job) { }
}
