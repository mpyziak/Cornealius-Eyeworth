using System.Windows.Forms;

namespace CornealiusEyeworth.UI.Options;

/// <summary>
/// Factory that creates all WinForms controls used by <see cref="ScheduleDialog"/>.
/// </summary>
internal interface IScheduleDialogControlFactory
{
    Label CreateSimpleInstructionLabel();
    TextBox CreateMinutesInput(string currentValue);
    LinkLabel CreateStandardToggle();
    LinkLabel CreateAdvancedToggle();
    Label CreateAdvancedInstructionLabel();
    TextBox CreateCronInput(string currentValue);
    Button CreateSaveButton();
    Button CreateCancelButton();
}
