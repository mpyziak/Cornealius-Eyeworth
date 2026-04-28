using System.Windows.Forms;

namespace CornealiusEyeworth.UI;

internal interface IOptionsDialogControlFactory
{
    Label CreateInstructionLabel();
    TextBox CreateMinutesInput(string currentValue);
    Label CreateLanguageLabel();
    ComboBox CreateLanguageDropdown(string? currentLanguage);
    Button CreateSaveButton();
    Button CreateCancelButton();
}