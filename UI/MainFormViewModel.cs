namespace CornealiusEyeworth.UI;

/// <summary>
/// Immutable snapshot of the data displayed by <see cref="MainForm"/>.
/// Rebuilt whenever the configuration changes.
/// </summary>
/// <param name="ScheduleDescription">Human-readable summary of the current trigger schedule.</param>
/// <param name="NextTrigger">The next date/time at which a reminder will fire.</param>
internal record MainFormViewModel(string ScheduleDescription, DateTime NextTrigger);
