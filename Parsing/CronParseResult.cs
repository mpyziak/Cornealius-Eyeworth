namespace CornealiusEyeworth.Parsing;

/// <summary>Represents the outcome of validating a CRON expression string.</summary>
internal record CronParseResult(string? CronExpression, string? ErrorMessage)
{
    public bool IsSuccess => CronExpression is not null;
    public static CronParseResult Ok(string cron) => new(cron, null);
    public static CronParseResult Fail(string error) => new(null, error);
}