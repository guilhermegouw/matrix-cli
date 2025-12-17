// Package logo renders the Matrix wordmark.
package logo

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/guilhermegouw/matrix-cli/internal/tui/styles"
)

// ASCII art for MATRIX logo.
// Using block characters for a clean, terminal-friendly look.
const matrixLogo = `
███╗   ███╗ █████╗ ████████╗██████╗ ██╗██╗  ██╗
████╗ ████║██╔══██╗╚══██╔══╝██╔══██╗██║╚██╗██╔╝
██╔████╔██║███████║   ██║   ██████╔╝██║ ╚███╔╝
██║╚██╔╝██║██╔══██║   ██║   ██╔══██╗██║ ██╔██╗
██║ ╚═╝ ██║██║  ██║   ██║   ██║  ██║██║██╔╝ ██╗
╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚═╝  ╚═╝
`

// Smaller logo for narrow spaces.
const matrixLogoSmall = `
╔╦╗╔═╗╔╦╗╦═╗╦═╗ ╦
║║║╠═╣ ║ ╠╦╝║╔╩╦╝
╩ ╩╩ ╩ ╩ ╩╚═╩╩ ╚═
`

// Render returns the Matrix logo with the current theme colors.
func Render() string {
	t := styles.CurrentTheme()
	logo := strings.TrimPrefix(matrixLogo, "\n")

	// Apply gradient from bright green to darker green.
	return styles.ApplyForegroundGrad(logo, t.Primary, t.Secondary)
}

// RenderSmall returns a smaller version of the logo.
func RenderSmall() string {
	t := styles.CurrentTheme()
	logo := strings.TrimPrefix(matrixLogoSmall, "\n")
	return styles.ApplyForegroundGrad(logo, t.Primary, t.Secondary)
}

// RenderWithTagline returns the logo with a Matrix-themed tagline.
func RenderWithTagline() string {
	t := styles.CurrentTheme()
	logo := Render()

	tagline := t.S().Muted.Render("Wake up, Neo...")

	return lipgloss.JoinVertical(lipgloss.Center, logo, "", tagline)
}

// Width returns the width of the full logo.
func Width() int {
	return lipgloss.Width(matrixLogo)
}

// Height returns the height of the full logo.
func Height() int {
	return lipgloss.Height(matrixLogo)
}
