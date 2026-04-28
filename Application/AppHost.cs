using System.Globalization;
using System.Windows.Forms;
using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Localization;
using CornealiusEyeworth.Notifications;
using CornealiusEyeworth.Parsing;
using CornealiusEyeworth.Scheduling;
using CornealiusEyeworth.UI;

namespace CornealiusEyeworth.Application;

/// <summary>
/// Composes all application dependencies, wires up lifecycle events,
/// and runs the WinForms message loop.
/// </summary>
internal class AppHost
{
    private readonly IConfigRepository _configRepository;
    private readonly INotificationService _notificationService;
    private readonly IMinutesInputParser _minutesInputParser;
    private readonly IMainFormControlFactory _mainFormControlFactory;
    private readonly IScheduleDialogControlFactory _scheduleDialogControlFactory;
    private readonly ILanguageDialogControlFactory _languageDialogControlFactory;

    public AppHost(
        IConfigRepository configRepository,
        INotificationService notificationService,
        IMinutesInputParser minutesInputParser,
        IMainFormControlFactory mainFormControlFactory,
        IScheduleDialogControlFactory scheduleDialogControlFactory,
        ILanguageDialogControlFactory languageDialogControlFactory)
    {
        _configRepository = configRepository;
        _notificationService = notificationService;
        _minutesInputParser = minutesInputParser;
        _mainFormControlFactory = mainFormControlFactory;
        _scheduleDialogControlFactory = scheduleDialogControlFactory;
        _languageDialogControlFactory = languageDialogControlFactory;
    }

    public async Task RunAsync()
    {
        var config = _configRepository.Load();

        // Apply language override before any UI is created.
        if (!string.IsNullOrWhiteSpace(config.Language))
        {
            var culture = new CultureInfo(config.Language);
            CultureInfo.DefaultThreadCurrentUICulture = culture;
            CultureInfo.DefaultThreadCurrentCulture = culture;
            Thread.CurrentThread.CurrentUICulture = culture;
            Thread.CurrentThread.CurrentCulture = culture;
        }

        var schedulerCts = new CancellationTokenSource();
        Task schedulerTask = Task.CompletedTask;

        var form = new MainForm(
            BuildViewModel(config),
            _mainFormControlFactory,
            _configRepository,
            _minutesInputParser,
            _scheduleDialogControlFactory,
            _languageDialogControlFactory);

        async Task RestartScheduler(Config c)
        {
            await schedulerCts.CancelAsync();
            await schedulerTask;
            schedulerCts = new CancellationTokenSource();
            schedulerTask = new SchedulerService(c, _notificationService, form.NotifyFired)
                .RunAsync(schedulerCts.Token);
        }

        form.ConfigSaved += updatedConfig =>
            Task.Run(async () =>
            {
                await RestartScheduler(updatedConfig);
                form.Invoke(() => form.UpdateViewModel(BuildViewModel(updatedConfig)));
            });

        form.FormClosed += async (_, _) =>
        {
            await schedulerCts.CancelAsync();
            await schedulerTask;
            Environment.Exit(0);
        };

        _notificationService.SendStartupNotification();
        await RestartScheduler(config);

        System.Windows.Forms.Application.Run(form);
    }

    private static MainFormViewModel BuildViewModel(Config config) => new(
        ScheduleDescription: string.Format(Strings.ScheduleDescription, string.Join(", ", config.MinutesOfHour)),
        NextTrigger: new NextTriggerProvider(config).GetNext()
    );
}