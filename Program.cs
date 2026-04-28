using System.Windows.Forms;
using CornealiusEyeworth.Application;
using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Localization;
using CornealiusEyeworth.Notifications;
using CornealiusEyeworth.Parsing;
using CornealiusEyeworth.UI;

Application.EnableVisualStyles();
Application.SetCompatibleTextRenderingDefault(false);
Application.SetColorMode(SystemColorMode.System);

try
{
    await new AppHost(
        configRepository:              new JsonConfigRepository(),
        notificationService:           new NotificationService(),
        minutesInputParser:            new MinutesInputParser(),
        mainFormControlFactory:        new MainFormControlFactory(),
        scheduleDialogControlFactory:  new ScheduleDialogControlFactory(),
        languageDialogControlFactory:  new LanguageDialogControlFactory()
    ).RunAsync();
}
catch (Exception ex)
{
    MessageBox.Show(
        string.Format(Strings.FatalErrorMessage, ex.Message, ex.StackTrace),
        Strings.FatalErrorTitle,
        MessageBoxButtons.OK,
        MessageBoxIcon.Error);
    Environment.Exit(1);
}