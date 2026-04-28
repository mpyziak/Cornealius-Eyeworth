using System.Windows.Forms;

namespace CornealiusEyeworth.UI;

/// <summary>
/// Factory that creates all WinForms controls used by <see cref="MainForm"/>.
/// </summary>
internal interface IMainFormControlFactory
{
    MenuStrip CreateMenuStrip(Action onScheduleClicked, Action onLanguageClicked, Action onAboutClicked, Action onGitHubClicked);
    Label CreateTitleLabel();
    Label CreateStatusLabel();
    Label CreateScheduleLabel(string scheduleDescription);
    Label CreateNextTriggerLabel(DateTime nextTrigger);
}