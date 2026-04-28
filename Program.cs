using System.Windows.Forms;
using CornealiusEyeworth;

try
{
    Application.EnableVisualStyles();
    Application.SetCompatibleTextRenderingDefault(false);

    var config = new ConfigLoader().Load();
    var notificationService = new NotificationService();
    var form = new MainForm(config);

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