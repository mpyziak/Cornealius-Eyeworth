namespace CornealiusEyeworth.Configuration;

/// <summary>
/// Immutable snapshot of the application user-configurable settings.
/// </summary>
/// <param name="NotificationMessage">The text shown in the eye-rest toast notification.</param>
/// <param name="MinutesOfHour">The minutes within each hour at which reminders are fired.</param>
/// <param name="Language">Optional BCP-47 culture tag override (e.g. "pl", "de"). Null means use the OS default.</param>
internal record Config(string NotificationMessage, int[] MinutesOfHour, string? Language = null);
