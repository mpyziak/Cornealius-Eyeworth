using System.Drawing;
using System.Windows.Forms;

namespace CornealiusEyeworth;

class OptionsDialogControlFactory
{
    public Label CreateInstructionLabel() => new()
    {
        Text = "Minutes of each hour at which Cornealious shall\nremind you to rest your eyes (e.g. 20, 40, 55):",
        Font = new Font("Segoe UI", 9.5f),
        AutoSize = true,
        Location = new Point(16, 16),
        ForeColor = SystemColors.WindowText,
        BackColor = Color.Transparent
    };

    public TextBox CreateMinutesInput(string currentValue) => new()
    {
        Text = currentValue,
        Font = new Font("Segoe UI", 10f),
        Location = new Point(16, 72),
        Size = new Size(272, 28),
        BackColor = SystemColors.Window,
        ForeColor = SystemColors.WindowText,
        BorderStyle = BorderStyle.FixedSingle
    };

    public Button CreateSaveButton()
    {
        var btn = new Button
        {
            Text = "Save",
            Location = new Point(128, 116),
            Size = new Size(80, 30),
            Font = new Font("Segoe UI", 9.5f),
            BackColor = SystemColors.Highlight,
            ForeColor = SystemColors.HighlightText,
            FlatStyle = FlatStyle.Flat,
            DialogResult = DialogResult.OK
        };
        btn.FlatAppearance.BorderSize = 0;
        return btn;
    }

    public Button CreateCancelButton()
    {
        var btn = new Button
        {
            Text = "Cancel",
            Location = new Point(216, 116),
            Size = new Size(72, 30),
            Font = new Font("Segoe UI", 9.5f),
            BackColor = SystemColors.Control,
            ForeColor = SystemColors.ControlText,
            FlatStyle = FlatStyle.Flat,
            DialogResult = DialogResult.Cancel
        };
        btn.FlatAppearance.BorderSize = 0;
        return btn;
    }
}
