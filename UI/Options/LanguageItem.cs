namespace CornealiusEyeworth.UI.Options;

/// <summary>Represents a language choice in a dropdown.</summary>
internal record LanguageItem(string DisplayName, string? Code)
{
    public override string ToString() => DisplayName;
}
