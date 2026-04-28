using System.Windows.Forms;

namespace CornealiusEyeworth.UI;

/// <summary>
/// Factory that creates all WinForms controls used by <see cref="OptionsDialog"/>.
/// Implement this interface to provide an alternative visual theme or layout.
/// </summary>
internal interface IOptionsDialogControlFactory
{
    /// <summary>Creates the instruction label shown above the minutes input.</summary>
    Label CreateInstructionLabel();

    /// <summary>Creates the text box pre-populated with <paramref name="currentValue"/>.</summary>
    TextBox CreateMinutesInput(string currentValue);

    /// <summary>Creates the Save button.</summary>
    Button CreateSaveButton();

    /// <summary>Creates the Cancel button.</summary>
    Button CreateCancelButton();
}
