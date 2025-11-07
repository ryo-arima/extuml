package repository

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/extuml/extuml/pkg/model/extuml"
)

// ExtumlRepository defines interface for loading extuml DSL
type ExtumlRepository interface {
	Load(path string) (*extuml.Document, error)
}

type extumlRepositoryImpl struct{}

// NewExtumlRepository creates a new extuml repository
func NewExtumlRepository() ExtumlRepository {
	return &extumlRepositoryImpl{}
}

func (r *extumlRepositoryImpl) Load(path string) (*extuml.Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read extuml: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Parse header
	var header string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "%%") {
			continue
		}
		header = line
		break
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read extuml: %w", err)
	}
	if header == "" || !strings.HasPrefix(header, "extuml") {
		return nil, fmt.Errorf("E100: DSL header not found (expected 'extuml classDiagram3D')")
	}

	// Parse body
	doc := &extuml.Document{
		Version: "0.1",
		Elements: &extuml.Elements{
			Classes:    []extuml.Class{},
			Interfaces: []extuml.Interface{},
			Enums:      []extuml.Enum{},
			Packages:   []extuml.Package{},
			Notes:      []extuml.Note{},
		},
	}

	var currentClass *extuml.Class
	var currentInterface *extuml.Interface
	var currentEnum *extuml.Enum

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "%%") {
			continue
		}

		// Parse class declaration
		if strings.HasPrefix(line, "class ") {
			className := strings.TrimSpace(strings.TrimPrefix(line, "class "))
			className = strings.TrimSuffix(className, " {")
			className = strings.TrimSpace(className)

			currentClass = &extuml.Class{
				ID:         className,
				Type:       "class",
				Name:       className,
				Attributes: []extuml.Attribute{},
				Operations: []extuml.Operation{},
			}
			currentInterface = nil
			currentEnum = nil
			continue
		}

		// Parse interface declaration
		if strings.HasPrefix(line, "interface ") {
			interfaceName := strings.TrimSpace(strings.TrimPrefix(line, "interface "))
			interfaceName = strings.TrimSuffix(interfaceName, " {")
			interfaceName = strings.TrimSpace(interfaceName)

			currentInterface = &extuml.Interface{
				ID:         interfaceName,
				Type:       "interface",
				Name:       interfaceName,
				Operations: []extuml.Operation{},
			}
			currentClass = nil
			currentEnum = nil
			continue
		}

		// Parse enum declaration
		if strings.HasPrefix(line, "enum ") {
			enumName := strings.TrimSpace(strings.TrimPrefix(line, "enum "))
			enumName = strings.TrimSuffix(enumName, " {")
			enumName = strings.TrimSpace(enumName)

			currentEnum = &extuml.Enum{
				ID:       enumName,
				Type:     "enum",
				Name:     enumName,
				Literals: []string{},
			}
			currentClass = nil
			currentInterface = nil
			continue
		}

		// Close block
		if line == "}" {
			if currentClass != nil {
				doc.Elements.Classes = append(doc.Elements.Classes, *currentClass)
				currentClass = nil
			} else if currentInterface != nil {
				doc.Elements.Interfaces = append(doc.Elements.Interfaces, *currentInterface)
				currentInterface = nil
			} else if currentEnum != nil {
				doc.Elements.Enums = append(doc.Elements.Enums, *currentEnum)
				currentEnum = nil
			}
			continue
		}

		// Parse class/interface members
		if currentClass != nil {
			r.parseClassMember(currentClass, line)
		} else if currentInterface != nil {
			r.parseInterfaceMember(currentInterface, line)
		} else if currentEnum != nil {
			r.parseEnumLiteral(currentEnum, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read extuml: %w", err)
	}

	return doc, nil
}

// parseClassMember parses a class member (attribute or operation)
func (r *extumlRepositoryImpl) parseClassMember(class *extuml.Class, line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// Check if it's a URL annotation
	if strings.HasPrefix(line, "@url:") {
		class.URL = strings.TrimSpace(strings.TrimPrefix(line, "@url:"))
		return
	}

	// Remove visibility modifiers
	if len(line) > 0 && (line[0] == '+' || line[0] == '-' || line[0] == '#' || line[0] == '~') {
		line = strings.TrimSpace(line[1:])
	}

	// Check if it's an operation (has parentheses)
	if strings.Contains(line, "(") {
		// Parse operation: name(params): returnType
		opName := line
		returnType := ""

		if idx := strings.Index(line, ":"); idx != -1 {
			opName = strings.TrimSpace(line[:idx])
			returnType = strings.TrimSpace(line[idx+1:])
		}

		// Extract just the method name (before parentheses)
		if idx := strings.Index(opName, "("); idx != -1 {
			opName = opName[:idx]
		}

		class.Operations = append(class.Operations, extuml.Operation{
			Name:       opName,
			ReturnType: returnType,
		})
	} else {
		// Parse attribute: type name
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			attrType := parts[0]
			attrName := parts[1]
			class.Attributes = append(class.Attributes, extuml.Attribute{
				Name: attrName,
				Type: attrType,
			})
		} else if len(parts) == 1 {
			// Just name, no type
			class.Attributes = append(class.Attributes, extuml.Attribute{
				Name: parts[0],
			})
		}
	}
}

// parseInterfaceMember parses an interface member (operation)
func (r *extumlRepositoryImpl) parseInterfaceMember(iface *extuml.Interface, line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// Check if it's a URL annotation
	if strings.HasPrefix(line, "@url:") {
		iface.URL = strings.TrimSpace(strings.TrimPrefix(line, "@url:"))
		return
	}

	// Remove visibility modifiers
	if len(line) > 0 && (line[0] == '+' || line[0] == '-' || line[0] == '#' || line[0] == '~') {
		line = strings.TrimSpace(line[1:])
	}

	// Parse operation
	opName := line
	returnType := ""

	if idx := strings.Index(line, ":"); idx != -1 {
		opName = strings.TrimSpace(line[:idx])
		returnType = strings.TrimSpace(line[idx+1:])
	}

	if idx := strings.Index(opName, "("); idx != -1 {
		opName = opName[:idx]
	}

	iface.Operations = append(iface.Operations, extuml.Operation{
		Name:       opName,
		ReturnType: returnType,
	})
}

// parseEnumLiteral parses an enum literal
func (r *extumlRepositoryImpl) parseEnumLiteral(enum *extuml.Enum, line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// Check if it's a URL annotation
	if strings.HasPrefix(line, "@url:") {
		enum.URL = strings.TrimSpace(strings.TrimPrefix(line, "@url:"))
		return
	}

	// Remove trailing comma if present
	line = strings.TrimSuffix(line, ",")
	enum.Literals = append(enum.Literals, line)
}
