using System.Windows.Forms;

namespace CornealiusEyeworth.UI;

/// <summary>
/// Factory that creates all WinForms controls used by <see cref="MainForm"/>.
/// Implement this interface to provide an alternative visual theme or layout.
/// </summary>
internal interface IMainFormControlFactory
{
    /// <summary>Creates the main menu strip wired to the provided action callbacks.</summary>
    MenuStrip CreateMenuStrip(Action onOptionsClicked, Action onAboutClicked, Action onGitHubClicked);

    /// <summary>Creates the application title label.</summary>
    Label CreateTitleLabel();

    /// <summary>Creates the "Serving" status label.</summary>
    Label CreateStatusLabel();

    /// <summary>Creates the label displaying the current schedule description.</summary>
    Label CreateScheduleLabel(string scheduleDescription);

    /// <summary>Creates the label displaying the next trigger time.</summary>
    Label CreateNextTriggerLabel(DateTime nextTrigger);
}
