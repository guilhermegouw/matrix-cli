// Package welcome provides the welcome/splash screen for Matrix CLI.
package welcome

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/guilhermegouw/matrix-cli/internal/tui/components/logo"
	"github.com/guilhermegouw/matrix-cli/internal/tui/styles"
	"github.com/guilhermegouw/matrix-cli/internal/tui/util"
)

// StartWizardMsg is sent when the user wants to start the wizard.
type StartWizardMsg struct{}

// Welcome displays the welcome screen with Matrix branding.
type Welcome struct {
	width  int
	height int
}

// New creates a new welcome screen.
func New() *Welcome {
	return &Welcome{}
}

// Init initializes the welcome screen.
func (w *Welcome) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (w *Welcome) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter", " ":
			return w, util.CmdHandler(StartWizardMsg{})
		case "q", "ctrl+c":
			return w, tea.Quit
		}
	}
	return w, nil
}

// View renders the welcome screen.
func (w *Welcome) View() string {
	t := styles.CurrentTheme()

	// Logo.
	logoView := logo.Render()

	// Matrix-themed messages.
	messages := []string{
		t.S().Text.Render("Wake up, Neo..."),
		"",
		t.S().Muted.Render("The Matrix has you..."),
		"",
		t.S().Text.Render("Follow the white rabbit."),
		"",
		t.S().Subtitle.Render("Let's configure your AI assistant."),
	}

	messageBlock := lipgloss.JoinVertical(lipgloss.Center, messages...)

	// Instructions.
	instructions := t.S().Muted.Render("Press Enter to begin setup • q to quit")

	// Combine everything.
	content := lipgloss.JoinVertical(lipgloss.Center,
		logoView,
		"",
		"",
		messageBlock,
		"",
		"",
		instructions,
	)

	// Center in available space.
	return lipgloss.Place(
		w.width, w.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

// SetSize sets the welcome screen size.
func (w *Welcome) SetSize(width, height int) {
	w.width = width
	w.height = height
}
