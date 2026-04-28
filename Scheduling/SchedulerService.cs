using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Notifications;
using Quartz;
using Quartz.Impl;

namespace CornealiusEyeworth.Scheduling;

/// <summary>
/// Starts a Quartz.NET scheduler that fires <see cref="EyeworthJob"/> at the
/// minutes-of-the-hour defined in <see cref="Config"/>.
/// Create a new instance whenever the configuration changes, and cancel the
/// previous instance's <see cref="CancellationToken"/> before doing so.
/// </summary>
internal class SchedulerService(Config config, INotificationService notificationService, Action? onJobFired = null)
{
    /// <summary>
    /// Runs the scheduler until <paramref name="cancellationToken"/> is cancelled,
    /// then shuts down cleanly.
    /// </summary>
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

    /// <summary>
    /// Builds a Quartz cron trigger that fires at the given
    /// <paramref name="minutes"/> of every hour, every day.
    /// </summary>
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
