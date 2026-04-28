using System.Drawing;
using System.Windows.Forms;
using CornealiusEyeworth.Localization;

namespace CornealiusEyeworth.UI;

/// <summary>
/// Default WinForms control factory for <see cref="OptionsDialog"/>.
/// Creates controls with the standard Cornealius visual style.
/// </summary>
internal class OptionsDialogControlFactory : IOptionsDialogControlFactory
{
    /// <inheritdoc/>
    public Label CreateInstructionLabel() => new()
    {
        Text = Strings.OptionsInstruction,
        Font = new Font("Segoe UI", 9.5f),
        AutoSize = true,
        Location = new Point(16, 16),
        ForeColor = SystemColors.WindowText,
        BackColor = Color.Transparent
    };

    /// <inheritdoc/>
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

    /// <inheritdoc/>
    public Button CreateSaveButton()
    {
        var btn = new Button
        {
            Text = Strings.ButtonSave,
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

    /// <inheritdoc/>
    public Button CreateCancelButton()
    {
        var btn = new Button
        {
            Text = Strings.ButtonCancel,
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
