using System.Windows.Forms;

namespace CornealiusEyeworth.UI;

/// <summary>
/// Factory that creates all WinForms controls used by <see cref="LanguageDialog"/>.
/// </summary>
internal interface ILanguageDialogControlFactory
{
    Label CreateInstructionLabel();
    ComboBox CreateLanguageDropdown(string? currentLanguage);
    Button CreateSaveButton();
    Button CreateCancelButton();
}