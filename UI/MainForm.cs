using System.Drawing;
using System.Windows.Forms;
using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Localization;
using CornealiusEyeworth.Parsing;

namespace CornealiusEyeworth.UI;

/// <summary>
/// The application main window. Provides access to Schedule and Language dialogs
/// via the Options menu, and About/GitHub via the Help menu.
/// </summary>
internal class MainForm : Form
{
    private readonly IConfigRepository _configRepository;
    private readonly IMinutesInputParser _minutesInputParser;
    private readonly IScheduleDialogControlFactory _scheduleDialogControlFactory;
    private readonly ILanguageDialogControlFactory _languageDialogControlFactory;

    private readonly Label _scheduleLabel;
    private readonly Label _nextTriggerLabel;

    public event Action<Config>? ConfigSaved;

    public MainForm(
        MainFormViewModel viewModel,
        IMainFormControlFactory controlFactory,
        IConfigRepository configRepository,
        IMinutesInputParser minutesInputParser,
        IScheduleDialogControlFactory scheduleDialogControlFactory,
        ILanguageDialogControlFactory languageDialogControlFactory)
    {
        _configRepository = configRepository;
        _minutesInputParser = minutesInputParser;
        _scheduleDialogControlFactory = scheduleDialogControlFactory;
        _languageDialogControlFactory = languageDialogControlFactory;

        Text = Strings.AppName;
        Size = new Size(400, 220);
        FormBorderStyle = FormBorderStyle.FixedSingle;
        MaximizeBox = false;
        StartPosition = FormStartPosition.CenterScreen;
        BackColor = SystemColors.Window;
        ForeColor = SystemColors.WindowText;
        MainMenuStrip = controlFactory.CreateMenuStrip(OpenScheduleDialog, OpenLanguageDialog, OpenAboutDialog, OpenGitHub);

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
        _nextTriggerLabel.Text = string.Format(Strings.NextTrigger, viewModel.NextTrigger);
    }

    public void NotifyFired() { }

    private void OpenScheduleDialog()
    {
        var dialog = new ScheduleDialog(_configRepository, _minutesInputParser, _scheduleDialogControlFactory);
        dialog.ConfigSaved += config => ConfigSaved?.Invoke(config);
        dialog.ShowDialog(this);
    }

    private void OpenLanguageDialog()
    {
        var dialog = new LanguageDialog(_configRepository, _languageDialogControlFactory);
        dialog.ShowDialog(this);
    }

    private void OpenAboutDialog() => new AboutDialog().ShowDialog(this);

    private static void OpenGitHub() =>
        System.Diagnostics.Process.Start(new System.Diagnostics.ProcessStartInfo
        {
            FileName = Strings.GitHubUrl,
            UseShellExecute = true
        });
}