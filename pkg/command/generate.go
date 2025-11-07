package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/extuml/extuml/pkg/config"
	"github.com/spf13/cobra"
)

// InitGenerateCmd creates the 'generate' subcommand which parses .extuml and
// generates .gl (glTF JSON).
func InitGenerateCmd() *cobra.Command {
	var (
		extumlPath string
		outputPath string
		htmlOutput string
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a .extuml diagram to .gl (glTF JSON)",
		Long:  "Generate a .extuml 3D UML diagram into a glTF 2.0 .gl JSON file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if extumlPath == "" && len(args) > 0 {
				extumlPath = args[0]
			}

			if strings.TrimSpace(extumlPath) == "" {
				return fmt.Errorf("extuml path is required (pass as arg or --extuml)")
			}

			if strings.TrimSpace(outputPath) == "" {
				return fmt.Errorf("output path is required (--output)")
			}

			if err := RunGenerate(extumlPath, outputPath, htmlOutput); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&extumlPath, "extuml", "e", "", "path to .extuml DSL file")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output .gl (glTF JSON) file path")
	cmd.Flags().StringVar(&htmlOutput, "html-output", "", "output HTML viewer file path (optional)")

	return cmd
}

// RunGenerate executes the generate command logic
func RunGenerate(extumlPath, outputPath, htmlOutput string) error {
	// Create config
	cfg := config.NewConfig()

	// Validate extuml file existence
	if _, statErr := os.Stat(extumlPath); statErr != nil {
		return fmt.Errorf("extuml file not found: %w", statErr)
	}

	// Execute generation via controller
	if err := cfg.GenerateCtrl.Generate(extumlPath, outputPath, htmlOutput); err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	fmt.Printf("Successfully generated: %s\n", outputPath)
	if htmlOutput != "" {
		fmt.Printf("Successfully generated: %s\n", htmlOutput)
	}
	return nil
}
