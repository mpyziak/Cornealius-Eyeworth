using System.Drawing;
using System.Windows.Forms;
using CornealiusEyeworth.Localization;

namespace CornealiusEyeworth.UI;

internal class LanguageDialogControlFactory : ILanguageDialogControlFactory
{
    private const int Left = 16;
    private const int Width = 272;

    public Label CreateInstructionLabel() => new()
    {
        Text = Strings.OptionsLanguageLabel,
        Font = new Font("Segoe UI", 9.5f),
        AutoSize = true,
        Location = new Point(Left, 16),
        ForeColor = SystemColors.WindowText,
        BackColor = Color.Transparent
    };

    public ComboBox CreateLanguageDropdown(string? currentLanguage)
    {
        var combo = new ComboBox
        {
            Font = new Font("Segoe UI", 9.5f),
            Location = new Point(Left, 40),
            Size = new Size(Width, 26),
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
            Location = new Point(116, 82),
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
            Location = new Point(204, 82),
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