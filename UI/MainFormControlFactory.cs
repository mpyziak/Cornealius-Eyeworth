using System.Drawing;
using System.Windows.Forms;
using CornealiusEyeworth.Localization;

namespace CornealiusEyeworth.UI;

/// <summary>
/// Default WinForms control factory for <see cref="MainForm"/>.
/// Creates controls with the standard Cornealius visual style.
/// </summary>
internal class MainFormControlFactory : IMainFormControlFactory
{
    /// <inheritdoc/>
    public MenuStrip CreateMenuStrip(Action onOptionsClicked, Action onAboutClicked, Action onGitHubClicked)
    {
        var triggerTimesItem = new ToolStripMenuItem(Strings.MenuTriggerTimes);
        triggerTimesItem.Click += (_, _) => onOptionsClicked();

        var optionsItem = new ToolStripMenuItem(Strings.MenuOptions);
        optionsItem.DropDownItems.Add(triggerTimesItem);

        var aboutItem = new ToolStripMenuItem(Strings.MenuAbout);
        aboutItem.Click += (_, _) => onAboutClicked();

        var gitHubItem = new ToolStripMenuItem(Strings.MenuGitHub);
        gitHubItem.Click += (_, _) => onGitHubClicked();

        var helpItem = new ToolStripMenuItem(Strings.MenuHelp);
        helpItem.DropDownItems.Add(aboutItem);
        helpItem.DropDownItems.Add(gitHubItem);

        var menu = new MenuStrip
        {
            BackColor = SystemColors.MenuBar,
            ForeColor = SystemColors.MenuText,
            RenderMode = ToolStripRenderMode.System
        };
        menu.Items.Add(optionsItem);
        menu.Items.Add(helpItem);
        return menu;
    }

    /// <inheritdoc/>
    public Label CreateTitleLabel() => new()
    {
        Text = Strings.AppTitle,
        Font = new Font("Segoe UI", 13f, FontStyle.Bold),
        AutoSize = true,
        Location = new Point(20, 42),
        ForeColor = SystemColors.Highlight,
        BackColor = Color.Transparent
    };

    /// <inheritdoc/>
    public Label CreateStatusLabel() => new()
    {
        Text = Strings.StatusServing,
        Font = new Font("Segoe UI", 10f),
        AutoSize = true,
        Location = new Point(22, 82),
        ForeColor = Color.LimeGreen,
        BackColor = Color.Transparent
    };

    /// <inheritdoc/>
    public Label CreateScheduleLabel(string scheduleDescription) => new()
    {
        Text = scheduleDescription,
        Font = new Font("Segoe UI", 9.5f),
        AutoSize = true,
        Location = new Point(22, 114),
        ForeColor = SystemColors.WindowText,
        BackColor = Color.Transparent
    };

    /// <inheritdoc/>
    public Label CreateNextTriggerLabel(DateTime nextTrigger) => new()
    {
        Text = string.Format(Strings.NextTrigger, nextTrigger),
        Font = new Font("Segoe UI", 9f),
        AutoSize = true,
        Location = new Point(22, 142),
        ForeColor = SystemColors.GrayText,
        BackColor = Color.Transparent
    };
}
