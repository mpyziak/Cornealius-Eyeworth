using System.Drawing;
using System.Windows.Forms;

namespace CornealiusEyeworth;

class MainFormControlFactory
{
    public Label CreateTitleLabel() => new()
    {
        Text = "\U0001F441\uFE0F  Cornealius Eyeworth",
        Font = new Font("Segoe UI", 13f, FontStyle.Bold),
        AutoSize = true,
        Location = new Point(20, 18),
        ForeColor = SystemColors.Highlight,
        BackColor = Color.Transparent
    };

    public Label CreateStatusLabel() => new()
    {
        Text = "\u25CF Serving",
        Font = new Font("Segoe UI", 10f),
        AutoSize = true,
        Location = new Point(22, 58),
        ForeColor = Color.LimeGreen,
        BackColor = Color.Transparent
    };

    public Label CreateScheduleLabel(string scheduleDescription) => new()
    {
        Text = scheduleDescription,
        Font = new Font("Segoe UI", 9.5f),
        AutoSize = true,
        Location = new Point(22, 90),
        ForeColor = SystemColors.WindowText,
        BackColor = Color.Transparent
    };

    public Label CreateNextTriggerLabel(DateTime nextTrigger) => new()
    {
        Text = $"Next trigger: {nextTrigger:HH:mm}",
        Font = new Font("Segoe UI", 9f),
        AutoSize = true,
        Location = new Point(22, 118),
        ForeColor = SystemColors.GrayText,
        BackColor = Color.Transparent
    };
}
