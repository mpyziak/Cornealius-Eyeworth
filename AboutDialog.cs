using System.Drawing;
using System.Windows.Forms;

namespace CornealiusEyeworth;

class AboutDialog : Form
{
    public AboutDialog()
    {
        Text = "About — Cornealius Eyeworth";
        Size = new Size(360, 240);
        FormBorderStyle = FormBorderStyle.FixedDialog;
        MaximizeBox = false;
        MinimizeBox = false;
        StartPosition = FormStartPosition.CenterParent;
        BackColor = SystemColors.Window;
        ForeColor = SystemColors.WindowText;

        var titleLabel = new Label
        {
            Text = "\U0001F441\uFE0F  Cornealius Eyeworth",
            Font = new Font("Segoe UI", 13f, FontStyle.Bold),
            AutoSize = true,
            Location = new Point(20, 20),
            ForeColor = SystemColors.Highlight,
            BackColor = Color.Transparent
        };

        var descLabel = new Label
        {
            Text = "A distinguished ocular butler who reminds you\nto rest your eyes at regular intervals.\n\nFollowing the 20-20-20 Rule, with decorum.",
            Font = new Font("Segoe UI", 9.5f),
            AutoSize = true,
            Location = new Point(20, 66),
            ForeColor = SystemColors.WindowText,
            BackColor = Color.Transparent
        };

        var versionLabel = new Label
        {
            Text = "Version 1.0  —  \u00A9 2026 mpyziak",
            Font = new Font("Segoe UI", 8.5f),
            AutoSize = true,
            Location = new Point(20, 158),
            ForeColor = SystemColors.GrayText,
            BackColor = Color.Transparent
        };

        var okButton = new Button
        {
            Text = "Close",
            Location = new Point(256, 150),
            Size = new Size(72, 30),
            Font = new Font("Segoe UI", 9.5f),
            BackColor = SystemColors.Highlight,
            ForeColor = SystemColors.HighlightText,
            FlatStyle = FlatStyle.Flat,
            DialogResult = DialogResult.OK
        };
        okButton.FlatAppearance.BorderSize = 0;

        AcceptButton = okButton;
        CancelButton = okButton;

        Controls.AddRange([titleLabel, descLabel, versionLabel, okButton]);
    }
}
