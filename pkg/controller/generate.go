package controller

import (
	"fmt"

	"github.com/extuml/extuml/pkg/usecase"
)

// GenerateController defines interface for generate command handling
type GenerateController interface {
	Generate(extumlPath, outputPath, htmlOutput string) error
}

type generateControllerImpl struct {
	usecase usecase.GenerateUsecase
}

// NewGenerateController creates a new generate controller
func NewGenerateController(uc usecase.GenerateUsecase) GenerateController {
	return &generateControllerImpl{
		usecase: uc,
	}
}

func (c *generateControllerImpl) Generate(extumlPath, outputPath, htmlOutput string) error {
	if extumlPath == "" || outputPath == "" {
		return fmt.Errorf("extuml and output paths are required")
	}

	if err := c.usecase.Execute(extumlPath, outputPath, htmlOutput); err != nil {
		return fmt.Errorf("generate failed: %w", err)
	}

	return nil
}
