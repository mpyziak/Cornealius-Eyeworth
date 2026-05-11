using System.Drawing;
using System.Windows.Forms;
using CornealiusEyeworth.Localization;

namespace CornealiusEyeworth.UI;

/// <summary>
/// Default WinForms control factory for <see cref="MainForm"/>.
/// </summary>
internal class MainFormControlFactory : IMainFormControlFactory
{
    public MenuStrip CreateMenuStrip(Action onScheduleClicked, Action onLanguageClicked, Action onAboutClicked, Action onGitHubClicked, Action onHelpClicked)
    {
        var scheduleItem = new ToolStripMenuItem(Strings.MenuTriggerTimes);
        scheduleItem.Click += (_, _) => onScheduleClicked();

        var languageItem = new ToolStripMenuItem(Strings.MenuLanguage);
        languageItem.Click += (_, _) => onLanguageClicked();

        var optionsItem = new ToolStripMenuItem(Strings.MenuOptions);
        optionsItem.DropDownItems.Add(scheduleItem);
        optionsItem.DropDownItems.Add(languageItem);

        var aboutItem = new ToolStripMenuItem(Strings.MenuAbout);
        aboutItem.Click += (_, _) => onAboutClicked();

        var howToUseItem = new ToolStripMenuItem(Strings.MenuHowToUse);
        howToUseItem.Click += (_, _) => onHelpClicked();

        var gitHubItem = new ToolStripMenuItem(Strings.MenuGitHub);
        gitHubItem.Click += (_, _) => onGitHubClicked();

        var helpItem = new ToolStripMenuItem(Strings.MenuHelp);
        helpItem.DropDownItems.Add(howToUseItem);
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

    public Label CreateTitleLabel() => new()
    {
        Text = Strings.AppTitle,
        Font = new Font("Segoe UI", 13f, FontStyle.Bold),
        AutoSize = true,
        Location = new Point(52, 42),
        ForeColor = AppColors.Accent.PrussianInkBlue,
        BackColor = Color.Transparent
    };

    public Label CreateStatusLabel() => new()
    {
        Text = Strings.StatusServing,
        Font = new Font("Segoe UI", 10f),
        AutoSize = true,
        Location = new Point(22, 82),
        ForeColor = AppColors.Accent.MutedSage,
        BackColor = Color.Transparent
    };

    public Label CreateScheduleLabel(string scheduleDescription) => new()
    {
        Text = scheduleDescription,
        Font = new Font("Segoe UI", 9.5f),
        AutoSize = true,
        Location = new Point(22, 114),
        ForeColor = SystemColors.WindowText,
        BackColor = Color.Transparent
    };

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