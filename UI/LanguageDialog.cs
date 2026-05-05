using System.Windows.Forms;
using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Localization;

namespace CornealiusEyeworth.UI;

/// <summary>
/// Modal dialog for overriding the display language.
/// </summary>
internal class LanguageDialog : Form
{
    private readonly IConfigRepository _configRepository;
    private readonly ComboBox _languageDropdown;

    public event Action<Config>? ConfigSaved;

    public LanguageDialog(
        IConfigRepository configRepository,
        ILanguageDialogControlFactory controlFactory)
    {
        _configRepository = configRepository;

        Text = Strings.LanguageDialogTitle;
        AutoSize = true;
        AutoSizeMode = AutoSizeMode.GrowAndShrink;
        MinimumSize = new System.Drawing.Size(260, 0);
        Padding = new System.Windows.Forms.Padding(0, 0, 16, 12);
        FormBorderStyle = FormBorderStyle.FixedDialog;
        MaximizeBox = false;
        MinimizeBox = false;
        StartPosition = FormStartPosition.CenterParent;
        BackColor = System.Drawing.SystemColors.Window;
        ForeColor = System.Drawing.SystemColors.WindowText;

        var currentConfig = _configRepository.Load();
        _languageDropdown = controlFactory.CreateLanguageDropdown(currentConfig.Language);

        var saveButton = controlFactory.CreateSaveButton();
        saveButton.Click += OnSaveClicked;

        var cancelButton = controlFactory.CreateCancelButton();

        AcceptButton = saveButton;
        CancelButton = cancelButton;

        Controls.AddRange([controlFactory.CreateInstructionLabel(), _languageDropdown, saveButton, cancelButton]);
    }

    private void OnSaveClicked(object? sender, EventArgs e)
    {
        var selectedLanguage = (_languageDropdown.SelectedItem as LanguageItem)?.Code;
        var existing = _configRepository.Load();
        var updated = existing with { Language = selectedLanguage };
        _configRepository.Save(updated);
        ConfigSaved?.Invoke(updated);

        if (selectedLanguage != existing.Language)
            MessageBox.Show(Strings.OptionsLanguageRestartNotice, Strings.AppName,
                MessageBoxButtons.OK, MessageBoxIcon.Information);
    }
}