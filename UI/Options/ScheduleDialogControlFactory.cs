using System.Drawing;
using System.Windows.Forms;
using CornealiusEyeworth.Localization;

namespace CornealiusEyeworth.UI.Options;

internal class ScheduleDialogControlFactory : IScheduleDialogControlFactory
{
    private const int Left = 16;
    private const int IndentLeft = 32;
    private const int InputWidth = 252;

    public Label CreateSimpleInstructionLabel() => new()
    {
        Text = Strings.OptionsInstruction,
        Font = new Font("Segoe UI", 9.5f),
        AutoSize = true,
        Location = new Point(Left, 16),
        ForeColor = SystemColors.WindowText,
        BackColor = Color.Transparent
    };

    public TextBox CreateMinutesInput(string currentValue) => new()
    {
        Text = currentValue,
        Font = new Font("Segoe UI", 10f),
        Location = new Point(Left, 54),
        Size = new Size(InputWidth, 26),
        BackColor = SystemColors.Window,
        ForeColor = SystemColors.WindowText,
        BorderStyle = BorderStyle.FixedSingle
    };

    public LinkLabel CreateStandardToggle() => new()
    {
        Text = Strings.ScheduleStandardToggle + " \u25bc",
        Font = new Font("Segoe UI", 9.5f),
        AutoSize = true,
        Location = new Point(Left, 0),
        LinkColor = SystemColors.HotTrack,
        BackColor = Color.Transparent,
        LinkBehavior = LinkBehavior.NeverUnderline
    };

    public LinkLabel CreateAdvancedToggle() => new()
    {
        Text = Strings.ScheduleAdvancedToggle + " \u25b6",
        Font = new Font("Segoe UI", 9.5f),
        AutoSize = true,
        Location = new Point(Left, 0),
        LinkColor = SystemColors.HotTrack,
        BackColor = Color.Transparent,
        LinkBehavior = LinkBehavior.NeverUnderline
    };

    public Label CreateAdvancedInstructionLabel() => new()
    {
        Text = Strings.ScheduleCronInstruction,
        Font = new Font("Segoe UI", 9.5f),
        AutoSize = true,
        Location = new Point(Left, 0),
        ForeColor = SystemColors.WindowText,
        BackColor = Color.Transparent
    };

    public TextBox CreateCronInput(string currentValue) => new()
    {
        Text = currentValue,
        Font = new Font("Segoe UI", 10f),
        Location = new Point(Left, 26),
        Size = new Size(InputWidth, 26),
        BackColor = SystemColors.Window,
        ForeColor = SystemColors.WindowText,
        BorderStyle = BorderStyle.FixedSingle
    };

    public Button CreateSaveButton()
    {
        var btn = new Button
        {
            Text = Strings.ButtonSave,
            Location = new Point(116, 0),
            Size = new Size(80, 28),
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
            Text = Strings.ButtonCancel,
            Location = new Point(204, 0),
            Size = new Size(80, 28),
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
