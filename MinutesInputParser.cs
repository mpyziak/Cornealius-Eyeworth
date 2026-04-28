namespace CornealiusEyeworth;

class MinutesInputParser : IMinutesInputParser
{
    public MinutesParseResult Parse(string input)
    {
        var parts = input.Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);

        if (parts.Length == 0)
            return MinutesParseResult.Fail("Cornealious insists on at least one minute. He has standards.");

        var minutes = new List<int>();
        foreach (var part in parts)
        {
            if (!int.TryParse(part, out var m) || m < 0 || m > 59)
                return MinutesParseResult.Fail(
                    $"'{part}' is not a valid minute. Cornealious expects whole numbers between 0 and 59.");

            minutes.Add(m);
        }

        return MinutesParseResult.Ok([.. minutes.Order()]);
    }
}
