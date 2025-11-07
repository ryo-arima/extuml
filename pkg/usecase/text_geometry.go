package usecase

import (
	"encoding/base64"
	"encoding/binary"
	"math"

	"github.com/extuml/extuml/pkg/model/gltf"
)

// TextGeometryGenerator generates billboard text quads for labels
type TextGeometryGenerator struct{}

// NewTextGeometryGenerator creates a new text geometry generator
func NewTextGeometryGenerator() *TextGeometryGenerator {
	return &TextGeometryGenerator{}
}

// TextQuad represents a single text label as a billboard quad
type TextQuad struct {
	Text     string
	Position [3]float64 // Local position relative to parent
	Width    float64
	Height   float64
}

// GenerateTextQuad generates a billboard quad mesh for text
// The quad is centered at position with fixed dimensions
func (g *TextGeometryGenerator) GenerateTextQuad(text string, position [3]float64) (mesh gltf.Mesh, buffers []byte) {
	// Larger dimensions for multi-line text
	width := float32(2.4)  // Wider for better visibility
	height := float32(2.4) // Taller for multi-line text

	// Create quad vertices (2 triangles)
	// Quad is in XY plane, facing +Z
	w := width / 2
	h := height / 2

	// 4 vertices for quad with interleaved position (x,y,z) and UV (u,v)
	// Each vertex: 5 floats (3 for position, 2 for UV)
	vertices := []float32{
		// Vertex 0: bottom-left
		-w, -h, 0, // position
		0, 1, // UV
		// Vertex 1: bottom-right
		w, -h, 0, // position
		1, 1, // UV
		// Vertex 2: top-right
		w, h, 0, // position
		1, 0, // UV
		// Vertex 3: top-left
		-w, h, 0, // position
		0, 0, // UV
	}

	// 6 indices for 2 triangles
	indices := []uint16{
		0, 1, 2, // first triangle
		0, 2, 3, // second triangle
	}

	// Create buffer data
	buffers = g.createTextBufferData(vertices, indices)

	// Create mesh with TRIANGLES mode
	mesh = gltf.Mesh{
		Name: "text_" + text,
		Primitives: []gltf.Primitive{
			{
				Attributes: map[string]int{
					"POSITION":   0, // Will be set by caller
					"TEXCOORD_0": 0, // Will be set by caller (same buffer, different accessor)
				},
				Indices: intPtr(1), // Will be set by caller
				Mode:    intPtr(4), // TRIANGLES
			},
		},
	}

	return
}

// createTextBufferData creates binary buffer for text quad (interleaved position + UV)
func (g *TextGeometryGenerator) createTextBufferData(vertices []float32, indices []uint16) []byte {
	vertexBytes := len(vertices) * 4 // float32 = 4 bytes
	indexBytes := len(indices) * 2   // uint16 = 2 bytes

	// Align index offset to 4-byte boundary
	indexOffset := vertexBytes
	if indexOffset%4 != 0 {
		indexOffset += 4 - (indexOffset % 4)
	}

	totalSize := indexOffset + indexBytes
	buffer := make([]byte, totalSize)

	// Write interleaved vertices (position + UV)
	for i, v := range vertices {
		binary.LittleEndian.PutUint32(buffer[i*4:], math.Float32bits(v))
	}

	// Write indices
	for i, idx := range indices {
		binary.LittleEndian.PutUint16(buffer[indexOffset+i*2:], idx)
	}

	return buffer
}

// CreateTextBufferURI creates a data URI for the text buffer
func (g *TextGeometryGenerator) CreateTextBufferURI(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	return "data:application/octet-stream;base64," + encoded
}
