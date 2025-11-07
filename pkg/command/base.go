package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd builds the root CLI command and wires subcommands.
func NewRootCmd() *cobra.Command {
	var root = &cobra.Command{
		Use:   "extuml",
		Short: "Generate glTF from .extuml 3D UML diagrams",
		Long:  "extuml is a CLI to render .extuml 3D UML diagrams into glTF 2.0 files.",
	}

	// Subcommands
	root.AddCommand(InitGenerateCmd())

	return root
}

// Execute runs the CLI.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
