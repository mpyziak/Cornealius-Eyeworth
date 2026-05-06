using System.Drawing;
using System.Windows.Forms;
using CornealiusEyeworth.Localization;

namespace CornealiusEyeworth.UI.Help;

/// <summary>
/// Modal dialog explaining how to use the application.
/// </summary>
internal class HelpDialog : Form
{
    public HelpDialog()
    {
        const int clientW = 440;
        const int margin  = 20;
        const int labelW  = clientW - margin * 2;

        Text = Strings.HelpDialogTitle;
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

        var bodyLabel = new Label
        {
            Text = Strings.HelpBody,
            Font = new Font("Segoe UI", 9.5f),
            AutoSize = true,
            MaximumSize = new Size(labelW, 0),
            Location = new Point(margin, 58),
            ForeColor = SystemColors.WindowText,
            BackColor = Color.Transparent
        };

        int bodyBottom = bodyLabel.GetPreferredSize(new Size(labelW, 0)).Height + 58 + 16;

        var okButton = new Button
        {
            Text = Strings.ButtonClose,
            Location = new Point(clientW - margin - 88, bodyBottom),
            Size = new Size(88, 30),
            Font = new Font("Segoe UI", 9.5f),
            BackColor = SystemColors.Highlight,
            ForeColor = SystemColors.HighlightText,
            FlatStyle = FlatStyle.Flat,
            DialogResult = DialogResult.OK
        };
        okButton.FlatAppearance.BorderSize = 0;

        ClientSize = new Size(clientW, bodyBottom + 30 + margin);

        AcceptButton = okButton;
        CancelButton = okButton;

        Controls.AddRange([titleLabel, bodyLabel, okButton]);
    }
}
