using Microsoft.Toolkit.Uwp.Notifications;

namespace CornealiusEyeworth;

class NotificationService
{
    public void SendStartupNotification()
    {
        new ToastContentBuilder()
            .AddText("Cornealius Eyeworth")
            .AddText("Cornealius is on duty.")
            .Show();
    }

    public void SendReminderNotification(string message)
    {
        new ToastContentBuilder()
            .AddText("Cornealius Eyeworth")
            .AddText(message)
            .Show();
    }
}
