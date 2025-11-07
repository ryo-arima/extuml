package usecase

import (
	"encoding/base64"
	"encoding/binary"
	"math"

	"github.com/extuml/extuml/pkg/model/extuml"
	"github.com/extuml/extuml/pkg/model/gltf"
)

// GeometryGenerator generates 3D geometry for extuml elements
type GeometryGenerator struct{}

// NewGeometryGenerator creates a new geometry generator
func NewGeometryGenerator() *GeometryGenerator {
	return &GeometryGenerator{}
}

// GenerateClassWireframe generates wireframe lines for a class with compartments
func (g *GeometryGenerator) GenerateClassWireframe(class extuml.Class, position [3]float64) (mesh gltf.Mesh, material gltf.Material, buffers []byte) {
	// Fixed cube dimensions for all classes
	size := 2.5
	width := size
	height := size
	depth := size

	// Generate simple wireframe cube (no compartment dividers)
	vertices, indices := g.createSimpleWireframeBox(float32(width), float32(height), float32(depth))

	// Create buffer data
	buffers = g.createBufferData(vertices, indices)

	// Create mesh with LINES mode
	mesh = gltf.Mesh{
		Name: class.Name + "_wireframe",
		Primitives: []gltf.Primitive{
			{
				Attributes: map[string]int{
					"POSITION": 0,
				},
				Indices: intPtr(1),
				Mode:    intPtr(1), // LINES mode
			},
		},
	}

	// Create emissive material for wireframe
	material = gltf.Material{
		Name: class.Name + "_material",
		PbrMetallicRoughness: &gltf.PbrMetallicRoughness{
			BaseColorFactor: []float64{0.2, 0.6, 0.8, 1.0}, // Blue color for classes
			MetallicFactor:  0.0,
			RoughnessFactor: 1.0,
		},
		EmissiveFactor: []float64{0.2, 0.6, 0.8}, // Make wireframe glow
		DoubleSided:    true,
	}

	return
}

// GenerateInterfaceWireframe generates wireframe lines for an interface with compartments
func (g *GeometryGenerator) GenerateInterfaceWireframe(iface extuml.Interface, position [3]float64) (mesh gltf.Mesh, material gltf.Material, buffers []byte) {
	width := 2.0
	// 2 compartments: name, operations
	nameHeight := float64(0.4)
	opHeight := float64(len(iface.Operations)) * 0.2
	if opHeight == 0 {
		opHeight = 0.2
	}
	height := nameHeight + opHeight
	depth := 0.5

	vertices, indices := g.createWireframeBox(float32(width), float32(height), float32(depth), 2, []float32{float32(nameHeight)})
	buffers = g.createBufferData(vertices, indices)

	mesh = gltf.Mesh{
		Name: iface.Name + "_wireframe",
		Primitives: []gltf.Primitive{
			{
				Attributes: map[string]int{
					"POSITION": 0,
				},
				Indices: intPtr(1),
				Mode:    intPtr(1), // LINES mode
			},
		},
	}

	material = gltf.Material{
		Name: iface.Name + "_material",
		PbrMetallicRoughness: &gltf.PbrMetallicRoughness{
			BaseColorFactor: []float64{0.8, 0.6, 0.2, 1.0}, // Orange color for interfaces
			MetallicFactor:  0.0,
			RoughnessFactor: 1.0,
		},
		EmissiveFactor: []float64{0.8, 0.6, 0.2},
		DoubleSided:    true,
	}

	return
}

// GenerateEnumWireframe generates wireframe lines for an enum
func (g *GeometryGenerator) GenerateEnumWireframe(enum extuml.Enum, position [3]float64) (mesh gltf.Mesh, material gltf.Material, buffers []byte) {
	width := 1.5
	// 2 compartments: name, literals
	nameHeight := float64(0.3)
	litHeight := float64(len(enum.Literals)) * 0.15
	if litHeight == 0 {
		litHeight = 0.15
	}
	height := nameHeight + litHeight
	depth := 0.4

	vertices, indices := g.createWireframeBox(float32(width), float32(height), float32(depth), 2, []float32{float32(nameHeight)})
	buffers = g.createBufferData(vertices, indices)

	mesh = gltf.Mesh{
		Name: enum.Name + "_wireframe",
		Primitives: []gltf.Primitive{
			{
				Attributes: map[string]int{
					"POSITION": 0,
				},
				Indices: intPtr(1),
				Mode:    intPtr(1), // LINES mode
			},
		},
	}

	material = gltf.Material{
		Name: enum.Name + "_material",
		PbrMetallicRoughness: &gltf.PbrMetallicRoughness{
			BaseColorFactor: []float64{0.6, 0.8, 0.6, 1.0}, // Green color for enums
			MetallicFactor:  0.0,
			RoughnessFactor: 1.0,
		},
		EmissiveFactor: []float64{0.6, 0.8, 0.6},
		DoubleSided:    true,
	}

	return
}

// GenerateInterfaceBox generates a box mesh for an interface (deprecated, use wireframe)
func (g *GeometryGenerator) GenerateInterfaceBox(iface extuml.Interface, position [3]float64) (mesh gltf.Mesh, material gltf.Material, buffers []byte) {
	width := 2.0
	height := 0.8 + float64(len(iface.Operations))*0.2
	depth := 0.5

	vertices, indices := g.createBoxGeometry(float32(width), float32(height), float32(depth))
	buffers = g.createBufferData(vertices, indices)

	mesh = gltf.Mesh{
		Name: iface.Name,
		Primitives: []gltf.Primitive{
			{
				Attributes: map[string]int{
					"POSITION": 0,
				},
				Indices: intPtr(1),
			},
		},
	}

	material = gltf.Material{
		Name: iface.Name + "_material",
		PbrMetallicRoughness: &gltf.PbrMetallicRoughness{
			BaseColorFactor: []float64{0.8, 0.6, 0.2, 1.0}, // Orange color for interfaces
			MetallicFactor:  0.1,
			RoughnessFactor: 0.9,
		},
		DoubleSided: false,
	}

	return
}

// GenerateEnumBox generates a box mesh for an enum
func (g *GeometryGenerator) GenerateEnumBox(enum extuml.Enum, position [3]float64) (mesh gltf.Mesh, material gltf.Material, buffers []byte) {
	width := 1.5
	height := 0.6 + float64(len(enum.Literals))*0.15
	depth := 0.4

	vertices, indices := g.createBoxGeometry(float32(width), float32(height), float32(depth))
	buffers = g.createBufferData(vertices, indices)

	mesh = gltf.Mesh{
		Name: enum.Name,
		Primitives: []gltf.Primitive{
			{
				Attributes: map[string]int{
					"POSITION": 0,
				},
				Indices: intPtr(1),
			},
		},
	}

	material = gltf.Material{
		Name: enum.Name + "_material",
		PbrMetallicRoughness: &gltf.PbrMetallicRoughness{
			BaseColorFactor: []float64{0.6, 0.8, 0.6, 1.0}, // Green color for enums
			MetallicFactor:  0.1,
			RoughnessFactor: 0.9,
		},
		DoubleSided: false,
	}

	return
}

// createBoxGeometry creates vertices and indices for a box
func (g *GeometryGenerator) createBoxGeometry(width, height, depth float32) ([]float32, []uint16) {
	w := width / 2
	h := height / 2
	d := depth / 2

	// 24 vertices (4 per face for proper normals)
	vertices := []float32{
		// Front face
		-w, -h, d, w, -h, d, w, h, d, -w, h, d,
		// Back face
		w, -h, -d, -w, -h, -d, -w, h, -d, w, h, -d,
		// Top face
		-w, h, d, w, h, d, w, h, -d, -w, h, -d,
		// Bottom face
		-w, -h, -d, w, -h, -d, w, -h, d, -w, -h, d,
		// Right face
		w, -h, d, w, -h, -d, w, h, -d, w, h, d,
		// Left face
		-w, -h, -d, -w, -h, d, -w, h, d, -w, h, -d,
	}

	// 36 indices (2 triangles per face * 6 faces)
	indices := []uint16{
		0, 1, 2, 0, 2, 3, // Front
		4, 5, 6, 4, 6, 7, // Back
		8, 9, 10, 8, 10, 11, // Top
		12, 13, 14, 12, 14, 15, // Bottom
		16, 17, 18, 16, 18, 19, // Right
		20, 21, 22, 20, 22, 23, // Left
	}

	return vertices, indices
}

// createWireframeBox creates wireframe edges for a box with horizontal compartment dividers
// compartments: number of compartments (e.g., 3 for class: name, attrs, ops)
// dividerHeights: heights of each compartment from top (length = compartments-1)
func (g *GeometryGenerator) createWireframeBox(width, height, depth float32, compartments int, dividerHeights []float32) ([]float32, []uint16) {
	w := width / 2
	h := height / 2
	d := depth / 2

	// 8 corner vertices for the box
	vertices := []float32{
		// Bottom 4 corners
		-w, -h, d, // 0: front-left-bottom
		w, -h, d, // 1: front-right-bottom
		w, -h, -d, // 2: back-right-bottom
		-w, -h, -d, // 3: back-left-bottom
		// Top 4 corners
		-w, h, d, // 4: front-left-top
		w, h, d, // 5: front-right-top
		w, h, -d, // 6: back-right-top
		-w, h, -d, // 7: back-left-top
	}

	// Add vertices for compartment dividers
	var yOffset float32 = h
	dividerVertices := make([]float32, 0)
	for _, divHeight := range dividerHeights {
		yOffset -= divHeight
		// 4 vertices per horizontal divider (front-left, front-right, back-right, back-left)
		dividerVertices = append(dividerVertices,
			-w, yOffset, d, // front-left
			w, yOffset, d, // front-right
			w, yOffset, -d, // back-right
			-w, yOffset, -d, // back-left
		)
	}
	vertices = append(vertices, dividerVertices...)

	// Create line indices
	indices := []uint16{
		// 12 edges of the box
		// Bottom face edges
		0, 1, 1, 2, 2, 3, 3, 0,
		// Top face edges
		4, 5, 5, 6, 6, 7, 7, 4,
		// Vertical edges
		0, 4, 1, 5, 2, 6, 3, 7,
	}

	// Add divider line indices
	baseIndex := uint16(8) // 8 corner vertices
	for i := 0; i < len(dividerHeights); i++ {
		offset := baseIndex + uint16(i*4)
		// Front edge
		indices = append(indices, offset, offset+1)
		// Right edge
		indices = append(indices, offset+1, offset+2)
		// Back edge
		indices = append(indices, offset+2, offset+3)
		// Left edge
		indices = append(indices, offset+3, offset)
	}

	return vertices, indices
}

// createSimpleWireframeBox creates wireframe edges for a simple box without compartment dividers
func (g *GeometryGenerator) createSimpleWireframeBox(width, height, depth float32) ([]float32, []uint16) {
	w := width / 2
	h := height / 2
	d := depth / 2

	// 8 corner vertices for the box
	vertices := []float32{
		// Bottom 4 corners
		-w, -h, d, // 0: front-left-bottom
		w, -h, d, // 1: front-right-bottom
		w, -h, -d, // 2: back-right-bottom
		-w, -h, -d, // 3: back-left-bottom
		// Top 4 corners
		-w, h, d, // 4: front-left-top
		w, h, d, // 5: front-right-top
		w, h, -d, // 6: back-right-top
		-w, h, -d, // 7: back-left-top
	}

	// Create line indices (12 edges of the box)
	indices := []uint16{
		// Bottom face edges
		0, 1, 1, 2, 2, 3, 3, 0,
		// Top face edges
		4, 5, 5, 6, 6, 7, 7, 4,
		// Vertical edges
		0, 4, 1, 5, 2, 6, 3, 7,
	}

	return vertices, indices
}

// createBufferData creates binary buffer data for vertices and indices
func (g *GeometryGenerator) createBufferData(vertices []float32, indices []uint16) []byte {
	// Calculate buffer size
	vertexBytes := len(vertices) * 4 // float32 = 4 bytes
	indexBytes := len(indices) * 2   // uint16 = 2 bytes

	// Align index offset to 4-byte boundary
	indexOffset := vertexBytes
	if indexOffset%4 != 0 {
		indexOffset += 4 - (indexOffset % 4)
	}

	totalSize := indexOffset + indexBytes
	buffer := make([]byte, totalSize)

	// Write vertices
	for i, v := range vertices {
		binary.LittleEndian.PutUint32(buffer[i*4:], math.Float32bits(v))
	}

	// Write indices
	for i, idx := range indices {
		binary.LittleEndian.PutUint16(buffer[indexOffset+i*2:], idx)
	}

	return buffer
}

// CreateBufferURI creates a data URI for the buffer
func (g *GeometryGenerator) CreateBufferURI(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	return "data:application/octet-stream;base64," + encoded
}

// intPtr returns a pointer to an int
func intPtr(i int) *int {
	return &i
}
