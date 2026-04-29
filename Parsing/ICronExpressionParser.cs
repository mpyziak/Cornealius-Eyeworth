namespace CornealiusEyeworth.Parsing;

/// <summary>Validates and normalises CRON expression input.</summary>
internal interface ICronExpressionParser
{
    /// <summary>Validates a raw Quartz CRON expression entered by the user.</summary>
    CronParseResult Parse(string input);

    /// <summary>Converts a comma-separated minutes list (e.g. "20, 40, 55") to a Quartz CRON expression.</summary>
    CronParseResult ParseMinutes(string minutesInput);
}