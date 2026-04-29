using CornealiusEyeworth.Localization;
using Quartz;

namespace CornealiusEyeworth.Parsing;

/// <summary>
/// Validates Quartz CRON expressions and converts simple minutes input to CRON.
/// </summary>
internal class CronExpressionParser : ICronExpressionParser
{
    public CronParseResult Parse(string input)
    {
        var trimmed = input.Trim();
        if (string.IsNullOrWhiteSpace(trimmed))
            return CronParseResult.Fail(Strings.ParseErrorCronEmpty);

        if (!CronExpression.IsValidExpression(trimmed))
            return CronParseResult.Fail(string.Format(Strings.ParseErrorCronInvalid, trimmed));

        return CronParseResult.Ok(trimmed);
    }

    public CronParseResult ParseMinutes(string minutesInput)
    {
        var parts = minutesInput.Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);

        if (parts.Length == 0)
            return CronParseResult.Fail(Strings.ParseErrorNoMinutes);

        var minutes = new List<int>();
        foreach (var part in parts)
        {
            if (!int.TryParse(part, out var m) || m < 0 || m > 59)
                return CronParseResult.Fail(string.Format(Strings.ParseErrorInvalidMinute, part));
            minutes.Add(m);
        }

        var cron = $"0 {string.Join(",", minutes.Order())} * * * ?";
        return CronParseResult.Ok(cron);
    }
}