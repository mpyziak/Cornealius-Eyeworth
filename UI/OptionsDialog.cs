using System.Windows.Forms;
using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Localization;
using CornealiusEyeworth.Parsing;

namespace CornealiusEyeworth.UI;

/// <summary>
/// Modal dialog that allows the user to edit the minutes-of-the-hour
/// at which reminders are scheduled. Validates input before persisting.
/// </summary>
internal class OptionsDialog : Form
{
    private readonly IConfigRepository _configRepository;
    private readonly IMinutesInputParser _parser;
    private readonly TextBox _minutesInput;

    /// <summary>Raised after the user successfully saves new settings.</summary>
    public event Action<Config>? ConfigSaved;

    /// <summary>
    /// Initialises the dialog with its required dependencies and pre-populates
    /// the minutes input from the current persisted configuration.
    /// </summary>
    public OptionsDialog(
        IConfigRepository configRepository,
        IMinutesInputParser parser,
        IOptionsDialogControlFactory controlFactory)
    {
        _configRepository = configRepository;
        _parser = parser;

        Text = Strings.OptionsDialogTitle;
        Size = new System.Drawing.Size(320, 200);
        FormBorderStyle = FormBorderStyle.FixedDialog;
        MaximizeBox = false;
        MinimizeBox = false;
        StartPosition = FormStartPosition.CenterParent;
        BackColor = System.Drawing.SystemColors.Window;
        ForeColor = System.Drawing.SystemColors.WindowText;

        var currentConfig = _configRepository.Load();
        _minutesInput = controlFactory.CreateMinutesInput(string.Join(", ", currentConfig.MinutesOfHour));

        var saveButton = controlFactory.CreateSaveButton();
        saveButton.Click += OnSaveClicked;

        var cancelButton = controlFactory.CreateCancelButton();

        AcceptButton = saveButton;
        CancelButton = cancelButton;

        Controls.AddRange([controlFactory.CreateInstructionLabel(), _minutesInput, saveButton, cancelButton]);
    }

    private void OnSaveClicked(object? sender, EventArgs e)
    {
        var result = _parser.Parse(_minutesInput.Text);

        if (!result.IsSuccess)
        {
            MessageBox.Show(result.ErrorMessage, Strings.AppName, MessageBoxButtons.OK, MessageBoxIcon.Warning);
            DialogResult = DialogResult.None;
            return;
        }

        var existing = _configRepository.Load();
        var updated = existing with { MinutesOfHour = result.Minutes! };
        _configRepository.Save(updated);
        ConfigSaved?.Invoke(updated);
    }
}
