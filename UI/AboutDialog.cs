using System.Drawing;
using System.Windows.Forms;
using CornealiusEyeworth.Localization;

namespace CornealiusEyeworth.UI;

/// <summary>
/// Modal dialog that displays application information — name, description,
/// version, and copyright.
/// </summary>
internal class AboutDialog : Form
{
    /// <summary>Initialises and lays out all controls for the About dialog.</summary>
    public AboutDialog()
    {
        Text = Strings.AboutDialogTitle;
        Size = new Size(360, 240);
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
            Location = new Point(20, 20),
            ForeColor = SystemColors.Highlight,
            BackColor = Color.Transparent
        };

        var descLabel = new Label
        {
            Text = Strings.AboutDescription,
            Font = new Font("Segoe UI", 9.5f),
            AutoSize = true,
            Location = new Point(20, 66),
            ForeColor = SystemColors.WindowText,
            BackColor = Color.Transparent
        };

        var versionLabel = new Label
        {
            Text = Strings.AboutVersion,
            Font = new Font("Segoe UI", 8.5f),
            AutoSize = true,
            Location = new Point(20, 158),
            ForeColor = SystemColors.GrayText,
            BackColor = Color.Transparent
        };

        var okButton = new Button
        {
            Text = Strings.ButtonClose,
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
