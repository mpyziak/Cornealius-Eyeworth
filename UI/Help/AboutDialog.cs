using System.Drawing;
using System.Windows.Forms;
using CornealiusEyeworth.Localization;

namespace CornealiusEyeworth.UI.Help;

/// <summary>
/// Modal dialog that displays application information - name, description,
/// version, and copyright.
/// </summary>
internal class AboutDialog : Form
{
    /// <summary>Initialises and lays out all controls for the About dialog.</summary>
    public AboutDialog()
    {
        const int clientW   = 400;
        const int margin    = 20;
        const int labelW    = clientW - margin * 2;

        Text = Strings.AboutDialogTitle;
        FormBorderStyle = FormBorderStyle.FixedDialog;
        MaximizeBox = false;
        MinimizeBox = false;
        StartPosition = FormStartPosition.CenterParent;
        BackColor = SystemColors.Window;
        ForeColor = SystemColors.WindowText;

        var titleLabel = new Label
        {
            Text = Strings.AppTitle,
            Font = new Font("Segoe UI", 13f, FontStyle.Bold),
            AutoSize = true,
            Location = new Point(margin, 20),
            ForeColor = SystemColors.Highlight,
            BackColor = Color.Transparent
        };

        var descLabel = new Label
        {
            Text = Strings.AboutDescription,
            Font = new Font("Segoe UI", 9.5f),
            AutoSize = true,
            MaximumSize = new Size(labelW, 0),
            Location = new Point(margin, 60),
            ForeColor = SystemColors.WindowText,
            BackColor = Color.Transparent
        };

        // Position version and button below the auto-sized description
        int descBottom = descLabel.GetPreferredSize(new Size(labelW, 0)).Height + 60 + 10;

        var versionLabel = new Label
        {
            Text = Strings.AboutVersion,
            Font = new Font("Segoe UI", 8.5f),
            AutoSize = true,
            Location = new Point(margin, descBottom),
            ForeColor = SystemColors.GrayText,
            BackColor = Color.Transparent
        };

        int buttonY = descBottom + 24;
        var okButton = new Button
        {
            Text = Strings.ButtonClose,
            Location = new Point(clientW - margin - 88, buttonY),
            Size = new Size(88, 30),
            Font = new Font("Segoe UI", 9.5f),
            BackColor = SystemColors.Highlight,
            ForeColor = SystemColors.HighlightText,
            FlatStyle = FlatStyle.Flat,
            DialogResult = DialogResult.OK
        };
        okButton.FlatAppearance.BorderSize = 0;

        ClientSize = new Size(clientW, buttonY + 30 + margin);

        AcceptButton = okButton;
        CancelButton = okButton;

        Controls.AddRange([titleLabel, descLabel, versionLabel, okButton]);
    }
}
