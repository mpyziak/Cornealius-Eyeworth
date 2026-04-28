using System.Windows.Forms;

namespace CornealiusEyeworth.UI;

/// <summary>
/// Factory that creates all WinForms controls used by <see cref="ScheduleDialog"/>.
/// </summary>
internal interface IScheduleDialogControlFactory
{
    Label CreateInstructionLabel();
    TextBox CreateMinutesInput(string currentValue);
    Button CreateSaveButton();
    Button CreateCancelButton();
}