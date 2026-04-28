using System.Windows.Forms;
using CornealiusEyeworth;

try
{
    Application.EnableVisualStyles();
    Application.SetCompatibleTextRenderingDefault(false);
    Application.SetColorMode(SystemColorMode.System);

    var config = new ConfigLoader().Load();
    var notificationService = new NotificationService();

    INextTriggerProvider nextTriggerProvider = new NextTriggerProvider(config);
    var viewModel = new MainFormViewModel(
        ScheduleDescription: $"Speaks out at minutes: {string.Join(", ", config.MinutesOfHour)} of every hour",
        NextTrigger: nextTriggerProvider.GetNext()
    );
    var form = new MainForm(viewModel, new MainFormControlFactory());

    using var cts = new CancellationTokenSource();

    form.FormClosed += (_, _) => cts.Cancel();

    var schedulerTask = Task.Run(() =>
        new SchedulerService(config, notificationService, form.NotifyFired)
            .RunAsync(cts.Token));

    notificationService.SendStartupNotification();
    Application.Run(form);

    cts.Cancel();
    await schedulerTask;
}
catch (Exception ex)
{
    MessageBox.Show(
        $"Cornealius Eyeworth failed to start:\n\n{ex.Message}\n\n{ex.StackTrace}",
        "Fatal Error",
        MessageBoxButtons.OK,
        MessageBoxIcon.Error);
}
