using System.Drawing;
using System.Windows.Forms;

namespace CornealiusEyeworth;

class MainForm : Form
{
    private readonly IConfigRepository _configRepository;
    private readonly IMinutesInputParser _minutesInputParser;
    private readonly OptionsDialogControlFactory _optionsDialogControlFactory;

    private readonly Label _scheduleLabel;
    private readonly Label _nextTriggerLabel;

    public event Action<Config>? ConfigSaved;

    public MainForm(
        MainFormViewModel viewModel,
        MainFormControlFactory controlFactory,
        IConfigRepository configRepository,
        IMinutesInputParser minutesInputParser,
        OptionsDialogControlFactory optionsDialogControlFactory)
    {
        _configRepository = configRepository;
        _minutesInputParser = minutesInputParser;
        _optionsDialogControlFactory = optionsDialogControlFactory;

        Text = "Cornealius Eyeworth";
        Size = new Size(400, 248);
        FormBorderStyle = FormBorderStyle.FixedSingle;
        MaximizeBox = false;
        StartPosition = FormStartPosition.CenterScreen;
        BackColor = SystemColors.Window;
        ForeColor = SystemColors.WindowText;
        MainMenuStrip = controlFactory.CreateMenuStrip(OpenOptionsDialog, OpenAboutDialog, OpenGitHub);

        _scheduleLabel = controlFactory.CreateScheduleLabel(viewModel.ScheduleDescription);
        _nextTriggerLabel = controlFactory.CreateNextTriggerLabel(viewModel.NextTrigger);

        Controls.AddRange([
            MainMenuStrip,
            controlFactory.CreateTitleLabel(),
            controlFactory.CreateStatusLabel(),
            _scheduleLabel,
            _nextTriggerLabel
        ]);
    }

    public void UpdateViewModel(MainFormViewModel viewModel)
    {
        _scheduleLabel.Text = viewModel.ScheduleDescription;
        _nextTriggerLabel.Text = $"Next trigger: {viewModel.NextTrigger:HH:mm}";
    }

    public void NotifyFired() { }

    private void OpenOptionsDialog()
    {
        var dialog = new OptionsDialog(_configRepository, _minutesInputParser, _optionsDialogControlFactory);
        dialog.ConfigSaved += config => ConfigSaved?.Invoke(config);
        dialog.ShowDialog(this);
    }

    private void OpenAboutDialog() => new AboutDialog().ShowDialog(this);

    private static void OpenGitHub() =>
        System.Diagnostics.Process.Start(new System.Diagnostics.ProcessStartInfo
        {
            FileName = "https://github.com/mpyziak/Cornealius-Eyeworth",
            UseShellExecute = true
        });
}
