using CornealiusEyeworth.Localization;

namespace CornealiusEyeworth.Parsing;

/// <summary>
/// Default implementation of <see cref="IMinutesInputParser"/>.
/// Accepts a comma-separated list of integers in the range 0–59 and returns
/// them sorted in ascending order.
/// </summary>
internal class MinutesInputParser : IMinutesInputParser
{
    /// <inheritdoc/>
    public MinutesParseResult Parse(string input)
    {
        var parts = input.Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);

        if (parts.Length == 0)
            return MinutesParseResult.Fail(Strings.ParseErrorNoMinutes);

        var minutes = new List<int>();
        foreach (var part in parts)
        {
            if (!int.TryParse(part, out var m) || m < 0 || m > 59)
                return MinutesParseResult.Fail(string.Format(Strings.ParseErrorInvalidMinute, part));

            minutes.Add(m);
        }

        return MinutesParseResult.Ok([.. minutes.Order()]);
    }
}
