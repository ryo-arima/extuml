package gltf

// glTF 2.0 structures
type GLTFAsset struct {
	Asset       Asset        `json:"asset"`
	Scenes      []Scene      `json:"scenes,omitempty"`
	Scene       int          `json:"scene,omitempty"`
	Nodes       []Node       `json:"nodes,omitempty"`
	Meshes      []Mesh       `json:"meshes,omitempty"`
	Materials   []Material   `json:"materials,omitempty"`
	Textures    []Texture    `json:"textures,omitempty"`
	Images      []Image      `json:"images,omitempty"`
	Samplers    []Sampler    `json:"samplers,omitempty"`
	Buffers     []Buffer     `json:"buffers,omitempty"`
	BufferViews []BufferView `json:"bufferViews,omitempty"`
	Accessors   []Accessor   `json:"accessors,omitempty"`
}

type Asset struct {
	Version   string `json:"version"`
	Generator string `json:"generator,omitempty"`
	Extras    any    `json:"extras,omitempty"`
}

type Scene struct {
	Nodes []int  `json:"nodes,omitempty"`
	Name  string `json:"name,omitempty"`
}

type Node struct {
	Name        string    `json:"name,omitempty"`
	Mesh        *int      `json:"mesh,omitempty"`
	Translation []float64 `json:"translation,omitempty"`
	Rotation    []float64 `json:"rotation,omitempty"`
	Scale       []float64 `json:"scale,omitempty"`
	Matrix      []float64 `json:"matrix,omitempty"`
	Children    []int     `json:"children,omitempty"`
	Extras      any       `json:"extras,omitempty"`
}

type Mesh struct {
	Name       string      `json:"name,omitempty"`
	Primitives []Primitive `json:"primitives"`
}

type Primitive struct {
	Attributes map[string]int `json:"attributes"`
	Indices    *int           `json:"indices,omitempty"`
	Material   *int           `json:"material,omitempty"`
	Mode       *int           `json:"mode,omitempty"`
}

type Material struct {
	Name                 string                `json:"name,omitempty"`
	PbrMetallicRoughness *PbrMetallicRoughness `json:"pbrMetallicRoughness,omitempty"`
	EmissiveFactor       []float64             `json:"emissiveFactor,omitempty"`
	DoubleSided          bool                  `json:"doubleSided,omitempty"`
}

type PbrMetallicRoughness struct {
	BaseColorFactor  []float64    `json:"baseColorFactor,omitempty"`
	BaseColorTexture *TextureInfo `json:"baseColorTexture,omitempty"`
	MetallicFactor   float64      `json:"metallicFactor,omitempty"`
	RoughnessFactor  float64      `json:"roughnessFactor,omitempty"`
}

type TextureInfo struct {
	Index    int `json:"index"`
	TexCoord int `json:"texCoord,omitempty"`
}

type Texture struct {
	Sampler *int `json:"sampler,omitempty"`
	Source  *int `json:"source,omitempty"`
}

type Image struct {
	URI      string `json:"uri,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	Name     string `json:"name,omitempty"`
}

type Sampler struct {
	MagFilter int `json:"magFilter,omitempty"`
	MinFilter int `json:"minFilter,omitempty"`
	WrapS     int `json:"wrapS,omitempty"`
	WrapT     int `json:"wrapT,omitempty"`
}

type Buffer struct {
	ByteLength int    `json:"byteLength"`
	URI        string `json:"uri,omitempty"`
}

type BufferView struct {
	Buffer     int  `json:"buffer"`
	ByteOffset int  `json:"byteOffset,omitempty"`
	ByteLength int  `json:"byteLength"`
	ByteStride *int `json:"byteStride,omitempty"`
	Target     *int `json:"target,omitempty"`
}

type Accessor struct {
	BufferView    *int      `json:"bufferView,omitempty"`
	ByteOffset    int       `json:"byteOffset,omitempty"`
	ComponentType int       `json:"componentType"`
	Count         int       `json:"count"`
	Type          string    `json:"type"`
	Max           []float64 `json:"max,omitempty"`
	Min           []float64 `json:"min,omitempty"`
}
