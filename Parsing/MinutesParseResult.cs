namespace CornealiusEyeworth.Parsing;

/// <summary>
/// Represents the outcome of parsing a user-supplied minutes string.
/// Use <see cref="Ok"/> or <see cref="Fail"/> factory methods to construct.
/// </summary>
internal record MinutesParseResult(int[]? Minutes, string? ErrorMessage)
{
    /// <summary>Returns <c>true</c> when parsing succeeded and <see cref="Minutes"/> is populated.</summary>
    public bool IsSuccess => Minutes is not null;

    /// <summary>Creates a successful result containing the parsed <paramref name="minutes"/>.</summary>
    public static MinutesParseResult Ok(int[] minutes) => new(minutes, null);

    /// <summary>Creates a failed result carrying the human-readable <paramref name="error"/> message.</summary>
    public static MinutesParseResult Fail(string error) => new(null, error);
}
