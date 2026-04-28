using System.Windows.Forms;
using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Localization;
using CornealiusEyeworth.Parsing;

namespace CornealiusEyeworth.UI;

internal class OptionsDialog : Form
{
    private readonly IConfigRepository _configRepository;
    private readonly IMinutesInputParser _parser;
    private readonly TextBox _minutesInput;
    private readonly ComboBox _languageDropdown;

    public event Action<Config>? ConfigSaved;

    public OptionsDialog(
        IConfigRepository configRepository,
        IMinutesInputParser parser,
        IOptionsDialogControlFactory controlFactory)
    {
        _configRepository = configRepository;
        _parser = parser;

        Text = Strings.OptionsDialogTitle;
        Size = new System.Drawing.Size(320, 260);
        FormBorderStyle = FormBorderStyle.FixedDialog;
        MaximizeBox = false;
        MinimizeBox = false;
        StartPosition = FormStartPosition.CenterParent;
        BackColor = System.Drawing.SystemColors.Window;
        ForeColor = System.Drawing.SystemColors.WindowText;

        var currentConfig = _configRepository.Load();
        _minutesInput = controlFactory.CreateMinutesInput(string.Join(", ", currentConfig.MinutesOfHour));
        _languageDropdown = controlFactory.CreateLanguageDropdown(currentConfig.Language);

        var saveButton = controlFactory.CreateSaveButton();
        saveButton.Click += OnSaveClicked;

        var cancelButton = controlFactory.CreateCancelButton();

        AcceptButton = saveButton;
        CancelButton = cancelButton;

        Controls.AddRange([
            controlFactory.CreateInstructionLabel(),
            _minutesInput,
            controlFactory.CreateLanguageLabel(),
            _languageDropdown,
            saveButton,
            cancelButton
        ]);
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

        var selectedLanguage = (_languageDropdown.SelectedItem as LanguageItem)?.Code;
        var existing = _configRepository.Load();
        var updated = existing with { MinutesOfHour = result.Minutes!, Language = selectedLanguage };
        _configRepository.Save(updated);
        ConfigSaved?.Invoke(updated);

        // If the language changed, prompt the user to restart.
        if (selectedLanguage != existing.Language)
            MessageBox.Show(Strings.OptionsLanguageRestartNotice, Strings.AppName,
                MessageBoxButtons.OK, MessageBoxIcon.Information);
    }
}