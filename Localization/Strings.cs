namespace CornealiusEyeworth.Localization;

/// <summary>
/// Centralised repository of all user-facing strings.
/// Replace this class (or swap it for a resource file) to support additional languages.
/// </summary>
internal static class Strings
{
    // ── Application ──────────────────────────────────────────────────────────
    public const string AppName        = "Cornealius Eyeworth";
    public const string AppTitle       = "\U0001F441\uFE0F  Cornealius Eyeworth";
    public const string GitHubUrl      = "https://github.com/mpyziak/Cornealius-Eyeworth";

    // ── Main window ──────────────────────────────────────────────────────────
    public const string StatusServing          = "\u25CF Serving";
    /// <summary>Format arg {0}: comma-separated list of minutes.</summary>
    public const string ScheduleDescription    = "Speaks out at minutes: {0} of every hour";
    /// <summary>Format arg {0}: next trigger time (HH:mm).</summary>
    public const string NextTrigger            = "Next trigger: {0:HH\\:mm}";

    // ── Menu ─────────────────────────────────────────────────────────────────
    public const string MenuOptions      = "Options";
    public const string MenuTriggerTimes = "Trigger times...";
    public const string MenuHelp         = "Help";
    public const string MenuAbout        = "About...";
    public const string MenuGitHub       = "GitHub...";

    // ── Options dialog ───────────────────────────────────────────────────────
    public const string OptionsDialogTitle   = "Options \u2014 Cornealius Eyeworth";
    public const string OptionsInstruction   = "Minutes of each hour at which Cornealious shall\nremind you to rest your eyes (e.g. 20, 40, 55):";
    public const string ButtonSave           = "Save";
    public const string ButtonCancel         = "Cancel";

    // ── About dialog ─────────────────────────────────────────────────────────
    public const string AboutDialogTitle  = "About \u2014 Cornealius Eyeworth";
    public const string AboutDescription  = "A distinguished ocular butler who reminds you\nto rest your eyes at regular intervals.\n\nFollowing the 20-20-20 Rule, with decorum.";
    public const string AboutVersion      = "Version 1.0  \u2014  \u00A9 2026 mpyziak";
    public const string ButtonClose       = "Close";

    // ── Notifications ────────────────────────────────────────────────────────
    public const string NotificationOnDuty = "Cornealius is on duty.";

    // ── Input parsing errors ─────────────────────────────────────────────────
    public const string ParseErrorNoMinutes      = "Cornealious insists on at least one minute. He has standards.";
    /// <summary>Format arg {0}: the invalid token entered by the user.</summary>
    public const string ParseErrorInvalidMinute  = "'{0}' is not a valid minute. Cornealious expects whole numbers between 0 and 59.";

    // ── Fatal error ───────────────────────────────────────────────────────────
    public const string FatalErrorTitle   = "Fatal Error";
    /// <summary>Format args: {0} = ex.Message, {1} = ex.StackTrace.</summary>
    public const string FatalErrorMessage = "Cornealius Eyeworth failed to start:\n\n{0}\n\n{1}";
}
