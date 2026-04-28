using System.Windows.Forms;

namespace CornealiusEyeworth;

class MainForm : Form
{
    private readonly Label _statusLabel;
    private readonly Label _scheduleLabel;

    public MainForm(Config config)
    {
        Text = "Cornealius Eyeworth";
        Size = new System.Drawing.Size(400, 220);
        FormBorderStyle = FormBorderStyle.FixedSingle;
        MaximizeBox = false;
        StartPosition = FormStartPosition.CenterScreen;
        BackColor = System.Drawing.Color.FromArgb(30, 30, 30);
        ForeColor = System.Drawing.Color.WhiteSmoke;

        var titleLabel = new Label
        {
            Text = "👁  Cornealius Eyeworth",
            Font = new System.Drawing.Font("Segoe UI", 13f, System.Drawing.FontStyle.Bold),
            AutoSize = true,
            Location = new System.Drawing.Point(20, 18),
            ForeColor = System.Drawing.Color.LightSkyBlue,
            BackColor = System.Drawing.Color.Transparent
        };

        _statusLabel = new Label
        {
            Text = "● Serving",
            Font = new System.Drawing.Font("Segoe UI", 10f),
            AutoSize = true,
            Location = new System.Drawing.Point(22, 58),
            ForeColor = System.Drawing.Color.LimeGreen,
            BackColor = System.Drawing.Color.Transparent
        };

        var minuteList = string.Join(", ", config.MinutesOfHour);
        _scheduleLabel = new Label
        {
            Text = $"Speaks out at minutes: {minuteList} of every hour",
            Font = new System.Drawing.Font("Segoe UI", 9.5f),
            AutoSize = true,
            Location = new System.Drawing.Point(22, 90),
            ForeColor = System.Drawing.Color.WhiteSmoke,
            BackColor = System.Drawing.Color.Transparent
        };

        var nextLabel = BuildNextFiresLabel(config.MinutesOfHour);
        nextLabel.Location = new System.Drawing.Point(22, 118);

        Controls.AddRange([titleLabel, _statusLabel, _scheduleLabel, nextLabel]);
    }

    public void NotifyFired() { }

    private static Label BuildNextFiresLabel(int[] minutes)
    {
        var now = DateTime.Now;
        var upcoming = minutes
            .Select(m => new DateTime(now.Year, now.Month, now.Day, now.Hour, m, 0))
            .Select(t => t <= now ? t.AddHours(1) : t)
            .OrderBy(t => t)
            .First();

        return new Label
        {
            Text = $"Next trigger: {upcoming:HH:mm}",
            Font = new System.Drawing.Font("Segoe UI", 9f),
            AutoSize = true,
            ForeColor = System.Drawing.Color.Goldenrod,
            BackColor = System.Drawing.Color.Transparent
        };
    }
}
