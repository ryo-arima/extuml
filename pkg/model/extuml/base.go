package extuml

// Minimal structures to parse header if needed in future extensions.
type Document struct {
	Version  string    `json:"version"`
	Meta     *Meta     `json:"meta,omitempty"`
	Elements *Elements `json:"elements,omitempty"`
}

type Meta struct {
	ID     string `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
}

type Elements struct {
	Classes    []Class     `json:"classes,omitempty"`
	Interfaces []Interface `json:"interfaces,omitempty"`
	Enums      []Enum      `json:"enums,omitempty"`
	Packages   []Package   `json:"packages,omitempty"`
	Notes      []Note      `json:"notes,omitempty"`
}

type Class struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Name       string      `json:"name"`
	URL        string      `json:"url,omitempty"`
	Attributes []Attribute `json:"attributes,omitempty"`
	Operations []Operation `json:"operations,omitempty"`
}

type Interface struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Name       string      `json:"name"`
	URL        string      `json:"url,omitempty"`
	Operations []Operation `json:"operations,omitempty"`
}

type Enum struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	URL      string   `json:"url,omitempty"`
	Literals []string `json:"literals,omitempty"`
}

type Package struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Children []string `json:"children,omitempty"`
}

type Note struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Text   string `json:"text"`
	Anchor string `json:"anchor,omitempty"`
}

type Attribute struct {
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
}

type Operation struct {
	Name       string `json:"name"`
	ReturnType string `json:"returnType,omitempty"`
}
