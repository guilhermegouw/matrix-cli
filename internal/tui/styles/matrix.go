package styles

// NewMatrixTheme creates a Matrix-inspired dark green theme.
func NewMatrixTheme() *Theme {
	return &Theme{
		Name:   "matrix",
		IsDark: true,

		// Matrix signature green tones
		Primary:   ParseHex("#00ff41"), // Bright matrix green
		Secondary: ParseHex("#008f11"), // Darker green
		Tertiary:  ParseHex("#003b00"), // Very dark green
		Accent:    ParseHex("#00ff41"), // Bright green accent

		// Dark backgrounds
		BgBase:    ParseHex("#0d0d0d"), // Near black
		BgSubtle:  ParseHex("#1a1a1a"), // Slightly lighter
		BgOverlay: ParseHex("#262626"), // Overlay background

		// Green-tinted foregrounds
		FgBase:   ParseHex("#00ff41"), // Bright green text
		FgMuted:  ParseHex("#008f11"), // Muted green
		FgSubtle: ParseHex("#005500"), // Subtle green

		// Borders
		Border:      ParseHex("#003b00"),
		BorderFocus: ParseHex("#00ff41"),

		// Status colors
		Success: ParseHex("#00ff41"), // Green
		Error:   ParseHex("#ff0000"), // Red
		Warning: ParseHex("#ffcc00"), // Yellow
		Info:    ParseHex("#00bfff"), // Cyan
	}
}
