package config

import (
	"log"

	"github.com/extuml/extuml/pkg/controller"
	"github.com/extuml/extuml/pkg/repository"
	"github.com/extuml/extuml/pkg/usecase"
)

// Config holds all dependencies for the application
type Config struct {
	ExtumlRepo   repository.ExtumlRepository
	GLTFRepo     repository.GLTFRepository
	HTMLRepo     repository.HTMLRepository
	GenerateUC   usecase.GenerateUsecase
	GenerateCtrl controller.GenerateController
}

// NewConfig creates and wires all dependencies
func NewConfig() *Config {
	extumlRepo := repository.NewExtumlRepository()
	gltfRepo := repository.NewGLTFRepository()
	htmlRepo, err := repository.NewHTMLRepository()
	if err != nil {
		log.Fatalf("failed to create HTML repository: %v", err)
	}
	generateUC := usecase.NewGenerateUsecase(extumlRepo, gltfRepo, htmlRepo)
	generateCtrl := controller.NewGenerateController(generateUC)

	return &Config{
		ExtumlRepo:   extumlRepo,
		GLTFRepo:     gltfRepo,
		HTMLRepo:     htmlRepo,
		GenerateUC:   generateUC,
		GenerateCtrl: generateCtrl,
	}
}
