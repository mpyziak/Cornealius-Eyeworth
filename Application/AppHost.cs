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
/// <para>
/// Keeping this logic here rather than in <c>Program.cs</c> ensures the
/// entry point stays minimal and this class remains independently testable.
/// </para>
/// </summary>
internal class AppHost
{
    private readonly IConfigRepository _configRepository;
    private readonly INotificationService _notificationService;
    private readonly IMinutesInputParser _minutesInputParser;
    private readonly IMainFormControlFactory _mainFormControlFactory;
    private readonly IOptionsDialogControlFactory _optionsDialogControlFactory;

    /// <summary>
    /// Initialises the host with its required service dependencies.
    /// </summary>
    public AppHost(
        IConfigRepository configRepository,
        INotificationService notificationService,
        IMinutesInputParser minutesInputParser,
        IMainFormControlFactory mainFormControlFactory,
        IOptionsDialogControlFactory optionsDialogControlFactory)
    {
        _configRepository = configRepository;
        _notificationService = notificationService;
        _minutesInputParser = minutesInputParser;
        _mainFormControlFactory = mainFormControlFactory;
        _optionsDialogControlFactory = optionsDialogControlFactory;
    }

    /// <summary>
    /// Runs the application: starts the scheduler, shows the main form,
    /// and enters the WinForms message loop.
    /// </summary>
    public async Task RunAsync()
    {
        var config = _configRepository.Load();

        var schedulerCts = new CancellationTokenSource();
        Task schedulerTask = Task.CompletedTask;

        var form = new MainForm(
            BuildViewModel(config),
            _mainFormControlFactory,
            _configRepository,
            _minutesInputParser,
            _optionsDialogControlFactory);

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

    /// <summary>
    /// Builds a fresh <see cref="MainFormViewModel"/> from the given <paramref name="config"/>.
    /// </summary>
    private static MainFormViewModel BuildViewModel(Config config) => new(
        ScheduleDescription: string.Format(Strings.ScheduleDescription, string.Join(", ", config.MinutesOfHour)),
        NextTrigger: new NextTriggerProvider(config).GetNext()
    );
}
