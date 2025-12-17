// Package cmd provides the CLI commands for Matrix.
package cmd

import (
	"github.com/spf13/cobra"
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
	}

	cmd.AddCommand(newVersionCmd())

	return cmd
}

// Execute runs the root command.
func Execute() error {
	return newRootCmd().Execute()
}
