using System.Windows.Forms;

namespace CornealiusEyeworth;

class OptionsDialog : Form
{
    private readonly IConfigRepository _configRepository;
    private readonly IMinutesInputParser _parser;
    private readonly TextBox _minutesInput;

    public event Action<Config>? ConfigSaved;

    public OptionsDialog(IConfigRepository configRepository, IMinutesInputParser parser, OptionsDialogControlFactory controlFactory)
    {
        _configRepository = configRepository;
        _parser = parser;

        Text = "Options — Cornealius Eyeworth";
        Size = new System.Drawing.Size(320, 200);
        FormBorderStyle = FormBorderStyle.FixedDialog;
        MaximizeBox = false;
        MinimizeBox = false;
        StartPosition = FormStartPosition.CenterParent;
        BackColor = SystemColors.Window;
        ForeColor = SystemColors.WindowText;

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
            MessageBox.Show(result.ErrorMessage, "Cornealius Eyeworth", MessageBoxButtons.OK, MessageBoxIcon.Warning);
            DialogResult = DialogResult.None;
            return;
        }

        var existing = _configRepository.Load();
        var updated = existing with { MinutesOfHour = result.Minutes! };
        _configRepository.Save(updated);
        ConfigSaved?.Invoke(updated);
    }
}
