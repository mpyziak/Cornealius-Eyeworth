using Quartz;
using Quartz.Impl;
using Quartz.Spi;

namespace CornealiusEyeworth;

class SchedulerService(Config config, NotificationService notificationService)
{
    public async Task RunAsync(CancellationToken cancellationToken)
    {
        var factory = new StdSchedulerFactory();
        var scheduler = await factory.GetScheduler(cancellationToken);

        scheduler.JobFactory = new EyeworthJobFactory(notificationService, config);

        var job = JobBuilder.Create<EyeworthJob>()
            .WithIdentity(EyeworthJob.Key)
            .Build();

        var trigger = BuildCronTrigger(config.MinutesOfHour);

        await scheduler.ScheduleJob(job, trigger, cancellationToken);
        await scheduler.Start(cancellationToken);

        Console.WriteLine("Cornealius Eyeworth is on duty. Press Ctrl+C to dismiss him.");
        Console.WriteLine($"Scheduled for minutes: {string.Join(", ", config.MinutesOfHour)}");

        await Task.Delay(Timeout.Infinite, cancellationToken);

        await scheduler.Shutdown(cancellationToken);
    }

    private static ITrigger BuildCronTrigger(int[] minutes)
    {
        // Build a cron expression that fires at each specified minute of every hour
        // e.g. minutes [20, 40, 55] => "0 20,40,55 * * * ?"
        var minuteList = string.Join(",", minutes);
        var cron = $"0 {minuteList} * * * ?";

        return TriggerBuilder.Create()
            .WithIdentity("EyeworthTrigger")
            .WithCronSchedule(cron)
            .Build();
    }
}

class EyeworthJobFactory(NotificationService notificationService, Config config) : IJobFactory
{
    public IJob NewJob(TriggerFiredBundle bundle, IScheduler scheduler) =>
        new EyeworthJob(notificationService, config);

    public void ReturnJob(IJob job) { }
}
