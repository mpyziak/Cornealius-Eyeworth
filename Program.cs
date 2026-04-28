using CornealiusEyeworth;

var config = new ConfigLoader().Load();
var notificationService = new NotificationService();
var schedulerService = new SchedulerService(config, notificationService);

using var cts = new CancellationTokenSource();
Console.CancelKeyPress += (_, e) =>
{
    e.Cancel = true;
    cts.Cancel();
};

notificationService.SendStartupNotification();
await schedulerService.RunAsync(cts.Token);
