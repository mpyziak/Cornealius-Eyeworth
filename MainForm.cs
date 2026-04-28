using System.Drawing;
using System.Windows.Forms;

namespace CornealiusEyeworth;

class MainForm : Form
{
    public MainForm(MainFormViewModel viewModel, MainFormControlFactory controlFactory)
    {
        Text = "Cornealius Eyeworth";
        Size = new Size(400, 220);
        FormBorderStyle = FormBorderStyle.FixedSingle;
        MaximizeBox = false;
        StartPosition = FormStartPosition.CenterScreen;
        BackColor = SystemColors.Window;
        ForeColor = SystemColors.WindowText;

        Controls.AddRange([
            controlFactory.CreateTitleLabel(),
            controlFactory.CreateStatusLabel(),
            controlFactory.CreateScheduleLabel(viewModel.ScheduleDescription),
            controlFactory.CreateNextTriggerLabel(viewModel.NextTrigger)
        ]);
    }

    public void NotifyFired() { }
}
