using System.Globalization;
using System.Windows.Forms;
using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Localization;
using CornealiusEyeworth.Notifications;
using CornealiusEyeworth.Parsing;
using CornealiusEyeworth.Scheduling;
using CornealiusEyeworth.UI;

namespace CornealiusEyeworth.Application;

internal class AppHost
{
    private readonly IConfigRepository _configRepository;
    private readonly INotificationService _notificationService;
    private readonly ICronExpressionParser _cronParser;
    private readonly IMainFormControlFactory _mainFormControlFactory;
    private readonly IScheduleDialogControlFactory _scheduleDialogControlFactory;
    private readonly ILanguageDialogControlFactory _languageDialogControlFactory;

    public AppHost(
        IConfigRepository configRepository,
        INotificationService notificationService,
        ICronExpressionParser cronParser,
        IMainFormControlFactory mainFormControlFactory,
        IScheduleDialogControlFactory scheduleDialogControlFactory,
        ILanguageDialogControlFactory languageDialogControlFactory)
    {
        _configRepository = configRepository;
        _notificationService = notificationService;
        _cronParser = cronParser;
        _mainFormControlFactory = mainFormControlFactory;
        _scheduleDialogControlFactory = scheduleDialogControlFactory;
        _languageDialogControlFactory = languageDialogControlFactory;
    }

    public async Task RunAsync()
    {
        var config = _configRepository.Load();

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
            _cronParser,
            _scheduleDialogControlFactory,
            _languageDialogControlFactory);

        async Task RestartScheduler(Config c)
        {
            await schedulerCts.CancelAsync();
            await schedulerTask;
            schedulerCts = new CancellationTokenSource();
            schedulerTask = new SchedulerService(c, _notificationService, () => form.Invoke(() => form.UpdateViewModel(BuildViewModel(c))))
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
        ScheduleDescription: CronDescriber.Describe(config.CronExpression),
        NextTrigger: new NextTriggerProvider(config).GetNext()
    );
}
