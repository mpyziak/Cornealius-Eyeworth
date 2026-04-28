namespace CornealiusEyeworth.Parsing;

/// <summary>
/// Parses a raw, user-supplied string of comma-separated minute values
/// into a validated <see cref="MinutesParseResult"/>.
/// </summary>
internal interface IMinutesInputParser
{
    /// <summary>
    /// Parses <paramref name="input"/> and returns either a successful result
    /// with sorted, validated minutes or a failure result with an error message.
    /// </summary>
    MinutesParseResult Parse(string input);
}
