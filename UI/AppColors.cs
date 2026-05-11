using System.Drawing;

namespace CornealiusEyeworth.UI;

internal static class AppColors
{
    // Logo — monochrome only; no accents or state colour
    public static class Logo
    {
        public static readonly Color GraphiteLight = ColorTranslator.FromHtml("#1C1C1C");
        public static readonly Color GraphiteDark  = ColorTranslator.FromHtml("#D6D6D6");
    }

    // Accent colours
    public static class Accent
    {
        /// <summary>Primary buttons, focus rings, active selections.</summary>
        public static readonly Color PrussianInkBlue = ColorTranslator.FromHtml("#233A5E");

        /// <summary>Secondary emphasis, hover states.</summary>
        public static readonly Color SmokedBrass = ColorTranslator.FromHtml("#9A7B3F");

        /// <summary>Calm / OK states, passive success.</summary>
        public static readonly Color MutedSage = ColorTranslator.FromHtml("#6B8B78");
    }

    // Surface colours
    public static class Surface
    {
        /// <summary>Light theme panels and backgrounds.</summary>
        public static readonly Color Light = ColorTranslator.FromHtml("#F5F5F5");

        /// <summary>Dark theme panels and backgrounds.</summary>
        public static readonly Color Dark = ColorTranslator.FromHtml("#181818");
    }
}
