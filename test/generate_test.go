package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/extuml/extuml/pkg/config"
	"github.com/extuml/extuml/pkg/model/gltf"
)

func TestGenerateCommand(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "test.extuml")
	outputPath := filepath.Join(tmpDir, "output.gl")

	// Create minimal extuml DSL file
	input := `extuml classDiagram3D

class Test {
  +string id
}
	`
	if err := os.WriteFile(inputPath, []byte(input), 0o644); err != nil {
		t.Fatalf("failed to write test input: %v", err)
	}

	// Use dependency injection to test
	cfg := config.NewConfig()
	err := cfg.GenerateCtrl.Generate(inputPath, outputPath, "") // No HTML output in test
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	// Read and validate output
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	var gltfAsset gltf.GLTFAsset
	if err := json.Unmarshal(data, &gltfAsset); err != nil {
		t.Fatalf("invalid glTF JSON: %v", err)
	}

	if gltfAsset.Asset.Version != "2.0" {
		t.Errorf("expected asset.version=2.0, got %s", gltfAsset.Asset.Version)
	}

	if gltfAsset.Asset.Generator != "extuml-cli v0.1" {
		t.Errorf("unexpected generator: %s", gltfAsset.Asset.Generator)
	}

	// Validate extras
	extras, ok := gltfAsset.Asset.Extras.(map[string]any)
	if !ok {
		t.Fatalf("expected extras to be map[string]any")
	}

	extumlExtras, ok := extras["extuml"].(map[string]any)
	if !ok {
		t.Fatalf("expected extuml extras")
	}

	version, ok := extumlExtras["version"].(string)
	if !ok || version != "0.1" {
		t.Errorf("expected version=0.1, got %v", version)
	}
}
