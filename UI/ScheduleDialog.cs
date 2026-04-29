using System.Windows.Forms;
using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Localization;
using CornealiusEyeworth.Parsing;

namespace CornealiusEyeworth.UI;

/// <summary>
/// Modal dialog for editing the reminder schedule.
/// Simple mode accepts a comma-separated minutes list; Advanced mode accepts a raw Quartz CRON.
/// Exactly one mode is active at a time — the radio buttons disable the inactive panel.
/// </summary>
internal class ScheduleDialog : Form
{
    private readonly IConfigRepository _configRepository;
    private readonly ICronExpressionParser _parser;

    private readonly RadioButton _simpleRadio;
    private readonly TextBox _minutesInput;
    private readonly RadioButton _advancedRadio;
    private readonly TextBox _cronInput;

    public event Action<Config>? ConfigSaved;

    public ScheduleDialog(
        IConfigRepository configRepository,
        ICronExpressionParser parser,
        IScheduleDialogControlFactory controlFactory)
    {
        _configRepository = configRepository;
        _parser = parser;

        Text = Strings.ScheduleDialogTitle;
        Size = new System.Drawing.Size(304, 262);
        FormBorderStyle = FormBorderStyle.FixedDialog;
        MaximizeBox = false;
        MinimizeBox = false;
        StartPosition = FormStartPosition.CenterParent;
        BackColor = System.Drawing.SystemColors.Window;
        ForeColor = System.Drawing.SystemColors.WindowText;

        var currentConfig = _configRepository.Load();
        var currentCron = currentConfig.CronExpression;
        var simpleMinutes = TryExtractSimpleMinutes(currentCron);
        var isSimple = simpleMinutes is not null;

        _simpleRadio    = controlFactory.CreateSimpleRadio();
        var simpleLabel = controlFactory.CreateSimpleInstructionLabel();
        _minutesInput   = controlFactory.CreateMinutesInput(simpleMinutes ?? string.Empty);

        _advancedRadio  = controlFactory.CreateAdvancedRadio();
        var advLabel    = controlFactory.CreateAdvancedInstructionLabel();
        _cronInput      = controlFactory.CreateCronInput(currentCron);

        // Start in the right mode
        _simpleRadio.Checked   = isSimple;
        _advancedRadio.Checked = !isSimple;
        SetSimpleEnabled(isSimple);

        _simpleRadio.CheckedChanged   += (_, _) => SetSimpleEnabled(_simpleRadio.Checked);
        _advancedRadio.CheckedChanged += (_, _) => SetSimpleEnabled(!_advancedRadio.Checked);

        var saveButton = controlFactory.CreateSaveButton();
        saveButton.Click += OnSaveClicked;
        var cancelButton = controlFactory.CreateCancelButton();

        AcceptButton = saveButton;
        CancelButton = cancelButton;

        Controls.AddRange([
            _simpleRadio, simpleLabel, _minutesInput,
            _advancedRadio, advLabel, _cronInput,
            saveButton, cancelButton
        ]);
    }

    private void SetSimpleEnabled(bool simpleActive)
    {
        _minutesInput.Enabled = simpleActive;
        _cronInput.Enabled    = !simpleActive;
    }

    private void OnSaveClicked(object? sender, EventArgs e)
    {
        CronParseResult result;

        if (_simpleRadio.Checked)
            result = _parser.ParseMinutes(_minutesInput.Text);
        else
            result = _parser.Parse(_cronInput.Text);

        if (!result.IsSuccess)
        {
            MessageBox.Show(result.ErrorMessage, Strings.AppName, MessageBoxButtons.OK, MessageBoxIcon.Warning);
            DialogResult = DialogResult.None;
            return;
        }

        var existing = _configRepository.Load();
        var updated  = existing with { CronExpression = result.CronExpression! };
        _configRepository.Save(updated);
        ConfigSaved?.Invoke(updated);
    }

    /// <summary>
    /// Returns a comma-separated minutes string if the CRON matches the simple
    /// "0 m1,m2,... * * * ?" pattern, or null if it is a free-form expression.
    /// </summary>
    private static string? TryExtractSimpleMinutes(string cron)
    {
        var parts = cron.Trim().Split(' ');
        if (parts.Length != 6) return null;
        if (parts[0] != "0") return null;
        if (parts[2] != "*" || parts[3] != "*" || parts[4] != "*" || parts[5] != "?") return null;
        // Validate each token in the minutes field is a valid integer 0-59
        var tokens = parts[1].Split(',');
        if (tokens.Any(t => !int.TryParse(t, out var m) || m < 0 || m > 59)) return null;
        return string.Join(", ", tokens);
    }
}