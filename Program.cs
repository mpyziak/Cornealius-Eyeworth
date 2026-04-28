using System.Text.Json;
using Microsoft.Toolkit.Uwp.Notifications;

var configPath = Path.Combine(AppContext.BaseDirectory, "config.json");
var config = JsonSerializer.Deserialize<Config>(File.ReadAllText(configPath))!;

Console.WriteLine("Cornealius Eyeworth is on duty. Press Ctrl+C to dismiss him.");

while (true)
{
    var now = DateTime.Now;
    var next = GetNextTrigger(now, config.MinutesOfHour);
    Console.WriteLine($"Next notification at {next:HH:mm}");
    await Task.Delay(next - now);

    new ToastContentBuilder()
        .AddText("Cornealius Eyeworth")
        .AddText(config.NotificationMessage)
        .Show();
}

static DateTime GetNextTrigger(DateTime now, int[] minutes)
{
    var candidate = minutes
        .Select(m => new DateTime(now.Year, now.Month, now.Day, now.Hour, m, 0))
        .Where(t => t > now)
        .OrderBy(t => t)
        .FirstOrDefault();

    if (candidate == default)
    {
        var nextHour = now.AddHours(1);
        candidate = new DateTime(nextHour.Year, nextHour.Month, nextHour.Day, nextHour.Hour, minutes.Min(), 0);
    }

    return candidate;
}

record Config(string NotificationMessage, int[] MinutesOfHour);
