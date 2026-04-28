namespace CornealiusEyeworth.Scheduling;

/// <summary>
/// Calculates the next scheduled trigger time based on the current configuration.
/// </summary>
internal interface INextTriggerProvider
{
    /// <summary>Returns the next <see cref="DateTime"/> at which a reminder will fire.</summary>
    DateTime GetNext();
}
