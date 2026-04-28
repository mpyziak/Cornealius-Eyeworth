using System.Drawing;
using System.Windows.Forms;
using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Localization;
using CornealiusEyeworth.Parsing;

namespace CornealiusEyeworth.UI;

/// <summary>
/// The application's main window. Displays the current schedule and next trigger
/// time, and provides access to the Options and About dialogs via the menu strip.
/// </summary>
internal class MainForm : Form
{
    private readonly IConfigRepository _configRepository;
    private readonly IMinutesInputParser _minutesInputParser;
    private readonly IOptionsDialogControlFactory _optionsDialogControlFactory;

    private readonly Label _scheduleLabel;
    private readonly Label _nextTriggerLabel;

    /// <summary>
    /// Raised after the user saves new settings in the Options dialog.
    /// The updated <see cref="Config"/> is passed to subscribers.
    /// </summary>
    public event Action<Config>? ConfigSaved;

    /// <summary>
    /// Initialises the form with the given view-model and injected dependencies.
    /// </summary>
    public MainForm(
        MainFormViewModel viewModel,
        IMainFormControlFactory controlFactory,
        IConfigRepository configRepository,
        IMinutesInputParser minutesInputParser,
        IOptionsDialogControlFactory optionsDialogControlFactory)
    {
        _configRepository = configRepository;
        _minutesInputParser = minutesInputParser;
        _optionsDialogControlFactory = optionsDialogControlFactory;

        Text = Strings.AppName;
        Size = new Size(400, 220);
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

    /// <summary>Refreshes the displayed schedule description and next trigger time.</summary>
    public void UpdateViewModel(MainFormViewModel viewModel)
    {
        _scheduleLabel.Text = viewModel.ScheduleDescription;
        _nextTriggerLabel.Text = string.Format(Strings.NextTrigger, viewModel.NextTrigger);
    }

    /// <summary>Called by the scheduler each time a reminder fires.</summary>
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
            FileName = Strings.GitHubUrl,
            UseShellExecute = true
        });
}
