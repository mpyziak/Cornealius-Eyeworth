using CornealiusEyeworth.Localization;
using Microsoft.Toolkit.Uwp.Notifications;

namespace CornealiusEyeworth.Notifications;

/// <summary>
/// Delivers Windows toast notifications via the UWP Notifications SDK.
/// </summary>
internal class NotificationService : INotificationService
{
    /// <inheritdoc/>
    public void SendStartupNotification()
    {
        new ToastContentBuilder()
            .AddText(Strings.AppName)
            .AddText(Strings.NotificationOnDuty)
            .Show();
    }

    /// <inheritdoc/>
    public void SendReminderNotification(string message)
    {
        new ToastContentBuilder()
            .AddText(Strings.AppName)
            .AddText(message)
            .Show();
    }
}
