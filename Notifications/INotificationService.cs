namespace CornealiusEyeworth.Notifications;

/// <summary>
/// Abstracts the delivery of toast notifications so the underlying
/// platform SDK can be replaced without touching scheduling or UI code.
/// </summary>
internal interface INotificationService
{
    /// <summary>Sends the "on duty" notification shown once at startup.</summary>
    void SendStartupNotification();

    /// <summary>Sends the eye-rest reminder notification.</summary>
    void SendReminderNotification();
}
