// Package cmd provides the CLI commands for Matrix.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/guilhermegouw/matrix-cli/internal/config"
	"github.com/guilhermegouw/matrix-cli/internal/tui"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "matrix",
		Short: "An AI-powered coding assistant CLI",
		Long: `Matrix is an AI-powered coding assistant that helps you write,
understand, and improve your code through conversation.

It supports multiple phases of development:
  - Matrix: Clarify requirements through dialogue
  - Planner: Design implementation strategy
  - Executor: Write and modify code`,
		RunE: runTUI,
	}

	cmd.AddCommand(newVersionCmd())

	return cmd
}

// runTUI launches the terminal user interface.
func runTUI(_ *cobra.Command, _ []string) error {
	// Check if this is first run.
	isFirstRun := config.IsFirstRun()

	// Load providers from catwalk (for the wizard).
	cfg := config.NewConfig()

	// Try to load providers even if config doesn't exist.
	providers, err := config.LoadProviders(cfg)
	if err != nil {
		// If we can't load providers, show an error.
		fmt.Fprintf(os.Stderr, "Warning: Failed to load providers: %v\n", err)
	}

	return tui.Run(providers, isFirstRun)
}

// Execute runs the root command.
func Execute() error {
	return newRootCmd().Execute()
}
