namespace CornealiusEyeworth.Configuration;

/// <summary>
/// Immutable snapshot of the application's user-configurable settings.
/// </summary>
/// <param name="NotificationMessage">The text shown in the eye-rest toast notification.</param>
/// <param name="MinutesOfHour">The minutes within each hour at which reminders are fired.</param>
internal record Config(string NotificationMessage, int[] MinutesOfHour);
