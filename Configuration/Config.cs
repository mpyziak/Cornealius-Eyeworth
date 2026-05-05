namespace CornealiusEyeworth.Configuration;

/// <summary>
/// Immutable snapshot of the application user-configurable settings.
/// </summary>
/// <param name="CronExpression">Quartz CRON expression defining when reminders fire (e.g. "0 20,40,55 * * * ?").</param>
/// <param name="Language">Optional BCP-47 culture tag override (e.g. "pl", "de"). Null means use the OS default.</param>
internal record Config(string CronExpression, string? Language = null);