using System.Windows.Forms;
using CornealiusEyeworth.Configuration;
using CornealiusEyeworth.Localization;
using CornealiusEyeworth.Parsing;

namespace CornealiusEyeworth.UI.Options;

/// <summary>
/// Modal dialog for editing the reminder schedule.
/// Uses an accordion: Standard (minutes list) and Advanced (CRON) sections.
/// Exactly one section is open at a time; both headers are always visible.
/// </summary>
internal class ScheduleDialog : Form
{
    private readonly IConfigRepository _configRepository;
    private readonly ICronExpressionParser _parser;

    private readonly Panel _standardHeader;
    private readonly LinkLabel _standardToggle;
    private readonly Panel _simplePanel;
    private readonly Panel _advancedHeader;
    private readonly LinkLabel _advancedToggle;
    private readonly Panel _cronPanel;
    private readonly Button _saveButton;
    private readonly Button _cancelButton;
    private readonly TextBox _minutesInput;
    private readonly TextBox _cronInput;

    private bool _advancedExpanded;

    // Layout constants
    private const int ClientW      = 400;
    private const int DialogMargin = 16;
    private const int ContentW     = ClientW - DialogMargin * 2;  // 368
    private const int HeaderH      = 26;
    private const int GapInner     = 6;
    private const int GapOuter     = 2;
    private const int ButtonGap    = 10;
    private const int ButtonH      = 28;
    private static readonly System.Drawing.Color HeaderBack     = System.Drawing.SystemColors.Control;
    private static readonly System.Drawing.Color SeparatorColor = System.Drawing.SystemColors.ControlDark;

    private const int InputH       = 26;
    private readonly int _simpleContentH;
    private readonly int _cronContentH;

    public event Action<Config>? ConfigSaved;

    public ScheduleDialog(
        IConfigRepository configRepository,
        ICronExpressionParser parser,
        IScheduleDialogControlFactory controlFactory)
    {
        _configRepository = configRepository;
        _parser = parser;

        Text = Strings.ScheduleDialogTitle;
        FormBorderStyle = FormBorderStyle.FixedDialog;
        MaximizeBox = false;
        MinimizeBox = false;
        StartPosition = FormStartPosition.CenterParent;
        BackColor = System.Drawing.SystemColors.Window;
        ForeColor = System.Drawing.SystemColors.WindowText;

        var currentConfig = _configRepository.Load();
        var currentCron = currentConfig.CronExpression;
        var simpleMinutes = TryExtractSimpleMinutes(currentCron);
        _advancedExpanded = simpleMinutes is null;

        // Build content using fixed width; labels word-wrap
        var simpleLabel = controlFactory.CreateSimpleInstructionLabel();
        var cronLabel   = controlFactory.CreateAdvancedInstructionLabel();
        simpleLabel.MaximumSize = new System.Drawing.Size(ContentW, 0);
        cronLabel.MaximumSize   = new System.Drawing.Size(ContentW, 0);

        // Standard section
        _minutesInput = controlFactory.CreateMinutesInput(simpleMinutes ?? string.Empty);
        simpleLabel.Location = new System.Drawing.Point(16, 0);
        int simpleLabelH = simpleLabel.GetPreferredSize(new System.Drawing.Size(ContentW, 0)).Height;
        _minutesInput.Location = new System.Drawing.Point(16, simpleLabelH + 6);
        _minutesInput.Size     = new System.Drawing.Size(ContentW, InputH);
        int simpleContentH = simpleLabelH + 6 + InputH;
        _simplePanel = new Panel
        {
            Size      = new System.Drawing.Size(ContentW + 16, simpleContentH),
            BackColor = System.Drawing.Color.Transparent
        };
        _simplePanel.Controls.AddRange([simpleLabel, _minutesInput]);

        // Standard header bar
        _standardToggle = controlFactory.CreateStandardToggle();
        _standardToggle.Location = new System.Drawing.Point(8, 4);
        _standardHeader = MakeHeaderBar(_standardToggle, ClientW);

        // Advanced section
        _cronInput = controlFactory.CreateCronInput(currentCron);
        cronLabel.Location = new System.Drawing.Point(16, 0);
        int cronLabelH = cronLabel.GetPreferredSize(new System.Drawing.Size(ContentW, 0)).Height;
        _cronInput.Location = new System.Drawing.Point(16, cronLabelH + 6);
        _cronInput.Size     = new System.Drawing.Size(ContentW, InputH);
        int cronContentH = cronLabelH + 6 + InputH;
        _cronPanel = new Panel
        {
            Size      = new System.Drawing.Size(ContentW + 16, cronContentH),
            BackColor = System.Drawing.Color.Transparent
        };
        _cronPanel.Controls.AddRange([cronLabel, _cronInput]);

        // Advanced header bar
        _advancedToggle = controlFactory.CreateAdvancedToggle();
        _advancedToggle.Location = new System.Drawing.Point(8, 4);
        _advancedHeader = MakeHeaderBar(_advancedToggle, ClientW);

        // Buttons — right-aligned
        _saveButton   = controlFactory.CreateSaveButton();
        _cancelButton = controlFactory.CreateCancelButton();
        _cancelButton.Left = ClientW - DialogMargin - 80;
        _saveButton.Left   = _cancelButton.Left - 8 - 80;

        _simpleContentH = simpleContentH;
        _cronContentH   = cronContentH;

        _standardToggle.LinkClicked += (_, _) => { if (_advancedExpanded)  { _advancedExpanded = false; UpdateLayout(); } };
        _advancedToggle.LinkClicked += (_, _) => { if (!_advancedExpanded) { _advancedExpanded = true;  UpdateLayout(); } };

        UpdateLayout();

        _saveButton.Click += OnSaveClicked;
        AcceptButton = _saveButton;
        CancelButton = _cancelButton;

        Controls.AddRange([
            _standardHeader, _simplePanel,
            _advancedHeader, _cronPanel,
            _saveButton, _cancelButton
        ]);
    }

    private static Panel MakeHeaderBar(LinkLabel toggle, int formWidth)
    {
        var bar = new Panel
        {
            Size      = new System.Drawing.Size(formWidth, HeaderH),
            Left      = 0,
            BackColor = HeaderBack
        };
        // 1px separator lines - top and bottom
        var topLine = new Panel    { Dock = DockStyle.Top,    Height = 1, BackColor = SeparatorColor };
        var bottomLine = new Panel { Dock = DockStyle.Bottom, Height = 1, BackColor = SeparatorColor };
        bar.Controls.AddRange([topLine, toggle, bottomLine]);
        return bar;
    }

    private void UpdateLayout()
    {
        int y = 8;

        // Standard header - always first
        _standardHeader.Top = y;
        y += HeaderH;

        if (!_advancedExpanded)
        {
            _simplePanel.Top     = y + GapInner;
            _simplePanel.Left    = 8;
            _simplePanel.Visible = true;
            y += GapInner + _simpleContentH + GapOuter;

            _standardToggle.Text = Strings.ScheduleStandardToggle + " \u25bc";
            _advancedToggle.Text = Strings.ScheduleAdvancedToggle + " \u25b6";
        }
        else
        {
            _simplePanel.Visible = false;
            y += GapOuter;

            _standardToggle.Text = Strings.ScheduleStandardToggle + " \u25b6";
            _advancedToggle.Text = Strings.ScheduleAdvancedToggle + " \u25bc";
        }

        // Advanced header
        _advancedHeader.Top = y;
        y += HeaderH;

        if (_advancedExpanded)
        {
            _cronPanel.Top     = y + GapInner;
            _cronPanel.Left    = 8;
            _cronPanel.Visible = true;
            y += GapInner + _cronContentH + GapOuter;
        }
        else
        {
            _cronPanel.Visible = false;
            y += GapOuter;
        }

        int buttonY = y + ButtonGap;
        _saveButton.Top   = buttonY;
        _cancelButton.Top = buttonY;
        ClientSize = new System.Drawing.Size(ClientW, buttonY + ButtonH + 12);
    }

    private void OnSaveClicked(object? sender, EventArgs e)
    {
        CronParseResult result;

        if (!_advancedExpanded)
            result = _parser.ParseMinutes(_minutesInput.Text);
        else
            result = _parser.Parse(_cronInput.Text);

        if (!result.IsSuccess)
        {
            MessageBox.Show(result.ErrorMessage, Strings.AppName, MessageBoxButtons.OK, MessageBoxIcon.Warning);
            DialogResult = DialogResult.None;
            return;
        }

        var existing = _configRepository.Load();
        var updated  = existing with { CronExpression = result.CronExpression! };
        _configRepository.Save(updated);
        ConfigSaved?.Invoke(updated);
    }

    /// <summary>
    /// Returns a comma-separated minutes string if the CRON matches the simple
    /// "0 m1,m2,... * * * ?" pattern, or null if it is a free-form expression.
    /// </summary>
    private static string? TryExtractSimpleMinutes(string cron)
    {
        var parts = cron.Trim().Split(' ');
        if (parts.Length != 6) return null;
        if (parts[0] != "0") return null;
        if (parts[2] != "*" || parts[3] != "*" || parts[4] != "*" || parts[5] != "?") return null;
        // Validate each token in the minutes field is a valid integer 0-59
        var tokens = parts[1].Split(',');
        if (tokens.Any(t => !int.TryParse(t, out var m) || m < 0 || m > 59)) return null;
        return string.Join(", ", tokens);
    }
}
