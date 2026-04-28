using System.Drawing;
using System.Windows.Forms;
using CornealiusEyeworth.Localization;

namespace CornealiusEyeworth.UI;

internal class OptionsDialogControlFactory : IOptionsDialogControlFactory
{
    public Label CreateInstructionLabel() => new()
    {
        Text = Strings.OptionsInstruction,
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

    public Label CreateLanguageLabel() => new()
    {
        Text = Strings.OptionsLanguageLabel,
        Font = new Font("Segoe UI", 9.5f),
        AutoSize = true,
        Location = new Point(16, 112),
        ForeColor = SystemColors.WindowText,
        BackColor = Color.Transparent
    };

    public ComboBox CreateLanguageDropdown(string? currentLanguage)
    {
        var combo = new ComboBox
        {
            Font = new Font("Segoe UI", 9.5f),
            Location = new Point(16, 136),
            Size = new Size(272, 28),
            DropDownStyle = ComboBoxStyle.DropDownList,
            BackColor = SystemColors.Window,
            ForeColor = SystemColors.WindowText
        };

        combo.Items.Add(new LanguageItem(Strings.OptionsLanguageDefault, null));
        combo.Items.Add(new LanguageItem("English", "en"));
        combo.Items.Add(new LanguageItem("Deutsch", "de"));
        combo.Items.Add(new LanguageItem("Polski", "pl"));

        var match = combo.Items.Cast<LanguageItem>()
            .FirstOrDefault(i => string.Equals(i.Code, currentLanguage, StringComparison.OrdinalIgnoreCase));
        combo.SelectedItem = match ?? combo.Items[0];

        return combo;
    }

    public Button CreateSaveButton()
    {
        var btn = new Button
        {
            Text = Strings.ButtonSave,
            Location = new Point(128, 180),
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
            Text = Strings.ButtonCancel,
            Location = new Point(216, 180),
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

/// <summary>Represents a language choice in the dropdown.</summary>
internal record LanguageItem(string DisplayName, string? Code)
{
    public override string ToString() => DisplayName;
}