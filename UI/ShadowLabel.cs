using System.Drawing;
using System.Drawing.Text;
using System.Windows.Forms;

namespace CornealiusEyeworth.UI;

/// <summary>
/// A <see cref="Label"/> that renders a soft glow/shadow behind its text.
/// </summary>
internal class ShadowLabel : Label
{
    public Color ShadowColor { get; init; } = Color.White;
    public int ShadowAlpha  { get; init; } = 30;
    public int ShadowRadius { get; init; } = 1;

    protected override void OnPaint(PaintEventArgs e)
    {
        e.Graphics.TextRenderingHint = TextRenderingHint.AntiAliasGridFit;

        using var sf          = new StringFormat(StringFormat.GenericTypographic);
        using var shadowBrush = new SolidBrush(Color.FromArgb(ShadowAlpha, ShadowColor));

        int r = ShadowRadius;
        for (int dx = -r; dx <= r; dx++)
        for (int dy = -r; dy <= r; dy++)
        {
            if (dx == 0 && dy == 0) continue;
            e.Graphics.DrawString(Text, Font, shadowBrush, new PointF(dx, dy), sf);
        }

        using var foreBrush = new SolidBrush(ForeColor);
        e.Graphics.DrawString(Text, Font, foreBrush, PointF.Empty, sf);
    }
}
