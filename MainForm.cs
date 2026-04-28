using System.Drawing;
using System.Windows.Forms;

namespace CornealiusEyeworth;

class MainForm : Form
{
    private readonly IConfigRepository _configRepository;
    private readonly IMinutesInputParser _minutesInputParser;
    private readonly OptionsDialogControlFactory _optionsDialogControlFactory;

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

        Controls.AddRange([
            MainMenuStrip,
            controlFactory.CreateTitleLabel(),
            controlFactory.CreateStatusLabel(),
            controlFactory.CreateScheduleLabel(viewModel.ScheduleDescription),
            controlFactory.CreateNextTriggerLabel(viewModel.NextTrigger)
        ]);
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
