package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/extuml/extuml/pkg/model/gltf"
)

// GLTFRepository defines interface for writing glTF output
type GLTFRepository interface {
	Write(path string, asset *gltf.GLTFAsset) error
}

type gltfRepositoryImpl struct{}

// NewGLTFRepository creates a new glTF repository
func NewGLTFRepository() GLTFRepository {
	return &gltfRepositoryImpl{}
}

func (r *gltfRepositoryImpl) Write(path string, asset *gltf.GLTFAsset) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	out, err := json.MarshalIndent(asset, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal glTF: %w", err)
	}

	if err := os.WriteFile(path, out, 0o644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
