using System.Windows.Forms;
using CornealiusEyeworth;

try
{
    Application.EnableVisualStyles();
    Application.SetCompatibleTextRenderingDefault(false);
    Application.SetColorMode(SystemColorMode.System);

    IConfigRepository configRepository = new JsonConfigRepository();
    var config = configRepository.Load();
    var notificationService = new NotificationService();
    var controlFactory = new MainFormControlFactory();
    IMinutesInputParser minutesInputParser = new MinutesInputParser();
    var optionsDialogControlFactory = new OptionsDialogControlFactory();

    MainFormViewModel BuildViewModel(Config c) => new(
        ScheduleDescription: $"Speaks out at minutes: {string.Join(", ", c.MinutesOfHour)} of every hour",
        NextTrigger: new NextTriggerProvider(c).GetNext()
    );

    var form = new MainForm(BuildViewModel(config), controlFactory, configRepository, minutesInputParser, optionsDialogControlFactory);

    var schedulerCts = new CancellationTokenSource();
    Task schedulerTask = Task.CompletedTask;

    async Task RestartScheduler(Config c)
    {
        await schedulerCts.CancelAsync();
        await schedulerTask;
        schedulerCts = new CancellationTokenSource();
        schedulerTask = new SchedulerService(c, notificationService, form.NotifyFired)
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

    notificationService.SendStartupNotification();
    await RestartScheduler(config);

    Application.Run(form);
}
catch (Exception ex)
{
    MessageBox.Show(
        $"Cornealius Eyeworth failed to start:\n\n{ex.Message}\n\n{ex.StackTrace}",
        "Fatal Error",
        MessageBoxButtons.OK,
        MessageBoxIcon.Error);
    Environment.Exit(1);
}
