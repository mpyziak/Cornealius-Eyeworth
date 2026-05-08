using CornealiusEyeworth.Localization;
using Microsoft.Toolkit.Uwp.Notifications;

namespace CornealiusEyeworth.Notifications;

/// <summary>
/// Delivers Windows toast notifications via the UWP Notifications SDK.
/// </summary>
internal class NotificationService : INotificationService
{
    private static readonly Random _random = Random.Shared;

    private static string PickRandom(string[] entries) =>
        entries[_random.Next(entries.Length)];

    /// <inheritdoc/>
    public void SendStartupNotification()
    {
        new ToastContentBuilder()
            .AddText(Strings.NotificationOnDuty)
            .AddText(PickRandom(Strings.NotificationQuips))
            .Show();
    }

    /// <inheritdoc/>
    public void SendReminderNotification()
    {
        new ToastContentBuilder()
            .AddText(PickRandom(Strings.NotificationReminders))
            .AddText(PickRandom(Strings.NotificationQuips))
            .Show();
    }
}
