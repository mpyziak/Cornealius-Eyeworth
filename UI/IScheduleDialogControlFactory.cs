using System.Windows.Forms;

namespace CornealiusEyeworth.UI;

/// <summary>
/// Factory that creates all WinForms controls used by <see cref="ScheduleDialog"/>.
/// </summary>
internal interface IScheduleDialogControlFactory
{
    RadioButton CreateSimpleRadio();
    Label CreateSimpleInstructionLabel();
    TextBox CreateMinutesInput(string currentValue);
    RadioButton CreateAdvancedRadio();
    Label CreateAdvancedInstructionLabel();
    TextBox CreateCronInput(string currentValue);
    Button CreateSaveButton();
    Button CreateCancelButton();
}