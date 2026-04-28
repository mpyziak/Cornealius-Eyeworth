using System.Drawing;
using System.Windows.Forms;

namespace CornealiusEyeworth;

class MainFormControlFactory
{
    public MenuStrip CreateMenuStrip(Action onOptionsClicked, Action onAboutClicked, Action onGitHubClicked)
    {
        var optionsItem = new ToolStripMenuItem("Options");
        var triggerTimesItem = new ToolStripMenuItem("Trigger times...");
        triggerTimesItem.Click += (_, _) => onOptionsClicked();
        optionsItem.DropDownItems.Add(triggerTimesItem);

        var helpItem = new ToolStripMenuItem("Help");
        var aboutItem = new ToolStripMenuItem("About...");
        aboutItem.Click += (_, _) => onAboutClicked();
        var gitHubItem = new ToolStripMenuItem("GitHub...");
        gitHubItem.Click += (_, _) => onGitHubClicked();
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
        Text = "\U0001F441\uFE0F  Cornealius Eyeworth",
        Font = new Font("Segoe UI", 13f, FontStyle.Bold),
        AutoSize = true,
        Location = new Point(20, 42),
        ForeColor = SystemColors.Highlight,
        BackColor = Color.Transparent
    };

    public Label CreateStatusLabel() => new()
    {
        Text = "\u25CF Serving",
        Font = new Font("Segoe UI", 10f),
        AutoSize = true,
        Location = new Point(22, 82),
        ForeColor = Color.LimeGreen,
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
        Text = $"Next trigger: {nextTrigger:HH:mm}",
        Font = new Font("Segoe UI", 9f),
        AutoSize = true,
        Location = new Point(22, 142),
        ForeColor = SystemColors.GrayText,
        BackColor = Color.Transparent
    };
}
