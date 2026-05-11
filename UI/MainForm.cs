using System.Drawing;
using System.Reflection;
using System.Windows.Forms;
using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Localization;
using CornealiusEyeworth.Parsing;
using CornealiusEyeworth.UI.Help;
using CornealiusEyeworth.UI.Options;

namespace CornealiusEyeworth.UI;

internal class MainForm : Form
{
    private readonly IConfigRepository _configRepository;
    private readonly ICronExpressionParser _cronParser;
    private readonly IScheduleDialogControlFactory _scheduleDialogControlFactory;
    private readonly ILanguageDialogControlFactory _languageDialogControlFactory;

    private readonly Label _scheduleLabel;
    private readonly Label _nextTriggerLabel;

    public event Action<Config>? ConfigSaved;

    public MainForm(
        MainFormViewModel viewModel,
        IMainFormControlFactory controlFactory,
        IConfigRepository configRepository,
        ICronExpressionParser cronParser,
        IScheduleDialogControlFactory scheduleDialogControlFactory,
        ILanguageDialogControlFactory languageDialogControlFactory)
    {
        _configRepository = configRepository;
        _cronParser = cronParser;
        _scheduleDialogControlFactory = scheduleDialogControlFactory;
        _languageDialogControlFactory = languageDialogControlFactory;

        Text = Strings.AppName;
        Icon = LoadEmbeddedIcon();
        Size = new Size(400, 220);
        FormBorderStyle = FormBorderStyle.FixedSingle;
        MaximizeBox = false;
        StartPosition = FormStartPosition.CenterScreen;
        ApplySurfaceColors();
        MainMenuStrip = controlFactory.CreateMenuStrip(OpenScheduleDialog, OpenLanguageDialog, OpenAboutDialog, OpenGitHub, OpenHelpDialog);

        _scheduleLabel    = controlFactory.CreateScheduleLabel(viewModel.ScheduleDescription);
        _nextTriggerLabel = controlFactory.CreateNextTriggerLabel(viewModel.NextTrigger);

        Controls.AddRange([
            MainMenuStrip,
            CreateLogoPictureBox(),
            controlFactory.CreateTitleLabel(),
            controlFactory.CreateStatusLabel(),
            _scheduleLabel,
            _nextTriggerLabel
        ]);
    }

    public void UpdateViewModel(MainFormViewModel viewModel)
    {
        _scheduleLabel.Text    = viewModel.ScheduleDescription;
        _nextTriggerLabel.Text = string.Format(Strings.NextTrigger, viewModel.NextTrigger);
    }

    public void NotifyFired() { }

    protected override void WndProc(ref Message m)
    {
        const int WM_SETTINGCHANGE = 0x001A;
        base.WndProc(ref m);
        if (m.Msg == WM_SETTINGCHANGE)
            ApplySurfaceColors();
    }

    private void ApplySurfaceColors()
    {
        bool dark = System.Windows.Forms.Application.IsDarkModeEnabled;
        BackColor = dark ? AppColors.Surface.Dark  : AppColors.Surface.Light;
        ForeColor = dark ? AppColors.Logo.GraphiteDark : AppColors.Logo.GraphiteLight;
    }

    private void OpenScheduleDialog()
    {
        var dialog = new ScheduleDialog(_configRepository, _cronParser, _scheduleDialogControlFactory);
        dialog.ConfigSaved += config => ConfigSaved?.Invoke(config);
        dialog.ShowDialog(this);
    }

    private void OpenLanguageDialog()
    {
        var dialog = new LanguageDialog(_configRepository, _languageDialogControlFactory);
        dialog.ShowDialog(this);
    }

    private void OpenAboutDialog() => new AboutDialog().ShowDialog(this);

    private void OpenHelpDialog() => new HelpDialog().ShowDialog(this);

    private static void OpenGitHub() =>
        System.Diagnostics.Process.Start(new System.Diagnostics.ProcessStartInfo
        {
            FileName = Strings.GitHubUrl,
            UseShellExecute = true
        });

    private static Icon? LoadEmbeddedIcon()
    {
        using var stream = Assembly.GetExecutingAssembly()
            .GetManifestResourceStream("CornealiusEyeworth.Properties.Resources.Logo.ico");
        return stream is null ? null : new Icon(stream);
    }

    private static PictureBox CreateLogoPictureBox()
    {
        var stream = Assembly.GetExecutingAssembly()
            .GetManifestResourceStream("CornealiusEyeworth.Properties.Resources.Logo.png");
        return new PictureBox
        {
            Image    = stream is null ? null : new Bitmap(stream),
            SizeMode = PictureBoxSizeMode.Zoom,
            Location = new Point(16, 38),
            Size     = new Size(30, 30),
            BackColor = Color.Transparent
        };
    }
}