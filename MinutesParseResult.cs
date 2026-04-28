namespace CornealiusEyeworth;

record MinutesParseResult(int[]? Minutes, string? ErrorMessage)
{
    public bool IsSuccess => Minutes is not null;

    public static MinutesParseResult Ok(int[] minutes) => new(minutes, null);
    public static MinutesParseResult Fail(string error) => new(null, error);
}
