package repository

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

//go:embed template/*.tmpl
var templateFS embed.FS

// HTMLRepository defines interface for writing HTML viewer output
type HTMLRepository interface {
	Write(path string, gltfPath string) error
}

type htmlRepositoryImpl struct {
	tmpl *template.Template
}

// NewHTMLRepository creates a new HTML repository
func NewHTMLRepository() (HTMLRepository, error) {
	tmpl, err := template.ParseFS(templateFS, "template/viewer.html.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	return &htmlRepositoryImpl{
		tmpl: tmpl,
	}, nil
}

type viewerData struct {
	GLTFPath string
}

func (r *htmlRepositoryImpl) Write(path string, gltfPath string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	// Calculate relative path from HTML to glTF
	htmlDir := filepath.Dir(path)
	relPath, err := filepath.Rel(htmlDir, gltfPath)
	if err != nil {
		// Fallback to absolute path if relative calculation fails
		relPath = gltfPath
	}

	data := viewerData{
		GLTFPath: relPath,
	}

	if err := r.tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}
