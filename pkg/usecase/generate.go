package usecase

import (
	"fmt"
	"time"

	"github.com/extuml/extuml/pkg/model/extuml"
	"github.com/extuml/extuml/pkg/model/gltf"
	"github.com/extuml/extuml/pkg/repository"
)

// GenerateUsecase defines interface for generate business logic
type GenerateUsecase interface {
	Execute(extumlPath, outputPath, htmlOutput string) error
}

type generateUsecaseImpl struct {
	extumlRepo repository.ExtumlRepository
	gltfRepo   repository.GLTFRepository
	htmlRepo   repository.HTMLRepository
	geomGen    *GeometryGenerator
	textGen    *TextGeometryGenerator
}

// NewGenerateUsecase creates a new generate usecase
func NewGenerateUsecase(extumlRepo repository.ExtumlRepository, gltfRepo repository.GLTFRepository, htmlRepo repository.HTMLRepository) GenerateUsecase {
	return &generateUsecaseImpl{
		extumlRepo: extumlRepo,
		gltfRepo:   gltfRepo,
		htmlRepo:   htmlRepo,
		geomGen:    NewGeometryGenerator(),
		textGen:    NewTextGeometryGenerator(),
	}
}

func (u *generateUsecaseImpl) Execute(extumlPath, outputPath, htmlOutput string) error {
	// Load extuml DSL
	doc, err := u.extumlRepo.Load(extumlPath)
	if err != nil {
		return fmt.Errorf("load extuml: %w", err)
	}

	// Create glTF asset with geometry
	gltfAsset := &gltf.GLTFAsset{
		Asset: gltf.Asset{
			Version:   "2.0",
			Generator: "extuml-cli v0.1",
			Extras: map[string]any{
				"extuml": map[string]any{
					"version":     doc.Version,
					"generatedAt": time.Now().UTC().Format(time.RFC3339),
				},
			},
		},
		Scenes: []gltf.Scene{
			{Nodes: []int{}, Name: "Scene"},
		},
		Scene:       0,
		Nodes:       []gltf.Node{},
		Meshes:      []gltf.Mesh{},
		Materials:   []gltf.Material{},
		Textures:    []gltf.Texture{},
		Images:      []gltf.Image{},
		Samplers:    []gltf.Sampler{},
		Buffers:     []gltf.Buffer{},
		BufferViews: []gltf.BufferView{},
		Accessors:   []gltf.Accessor{},
	}

	// Generate geometry if elements exist
	if doc.Elements != nil {
		u.generateGeometry(doc, gltfAsset)

		// Calculate scene bounds and add camera hint to extras
		bounds := u.calculateSceneBounds(gltfAsset)
		if gltfAsset.Asset.Extras == nil {
			gltfAsset.Asset.Extras = make(map[string]any)
		}
		if extrasMap, ok := gltfAsset.Asset.Extras.(map[string]any); ok {
			extrasMap["camera"] = bounds
		}
	}

	// Write glTF output
	if err := u.gltfRepo.Write(outputPath, gltfAsset); err != nil {
		return fmt.Errorf("write glTF: %w", err)
	}

	// Write HTML viewer if requested
	if htmlOutput != "" {
		if err := u.htmlRepo.Write(htmlOutput, outputPath); err != nil {
			return fmt.Errorf("write HTML: %w", err)
		}
	}

	return nil
}

func (u *generateUsecaseImpl) generateGeometry(doc *extuml.Document, asset *gltf.GLTFAsset) {
	nodeIndex := 0
	spacing := 3.0 // Space between elements

	// Generate classes
	for i, class := range doc.Elements.Classes {
		position := [3]float64{float64(i) * spacing, 0, 0}
		u.addClassToScene(class, position, asset, &nodeIndex)
	}

	// Generate interfaces
	for i, iface := range doc.Elements.Interfaces {
		position := [3]float64{float64(i) * spacing, spacing, 0}
		u.addInterfaceToScene(iface, position, asset, &nodeIndex)
	}

	// Generate enums
	for i, enum := range doc.Elements.Enums {
		position := [3]float64{float64(i) * spacing, -spacing, 0}
		u.addEnumToScene(enum, position, asset, &nodeIndex)
	}

	// Update scene nodes
	if len(asset.Nodes) > 0 {
		nodeIndices := make([]int, len(asset.Nodes))
		for i := range asset.Nodes {
			nodeIndices[i] = i
		}
		asset.Scenes[0].Nodes = nodeIndices
	}
}

func (u *generateUsecaseImpl) addClassToScene(class extuml.Class, position [3]float64, asset *gltf.GLTFAsset, nodeIndex *int) {
	mesh, material, bufferData := u.geomGen.GenerateClassWireframe(class, position)

	meshIdx := len(asset.Meshes)
	materialIdx := len(asset.Materials)
	bufferIdx := len(asset.Buffers)

	// Add buffer
	asset.Buffers = append(asset.Buffers, gltf.Buffer{
		ByteLength: len(bufferData),
		URI:        u.geomGen.CreateBufferURI(bufferData),
	})

	// Calculate vertex and index count from buffer
	vertexCount, indexCount, vertexBytes, indexOffset := u.calculateBufferLayout(bufferData)

	// Add buffer views
	positionBufferView := len(asset.BufferViews)
	asset.BufferViews = append(asset.BufferViews, gltf.BufferView{
		Buffer:     bufferIdx,
		ByteOffset: 0,
		ByteLength: vertexBytes,
		Target:     intPtr(34962), // ARRAY_BUFFER
	})

	indicesBufferView := len(asset.BufferViews)
	asset.BufferViews = append(asset.BufferViews, gltf.BufferView{
		Buffer:     bufferIdx,
		ByteOffset: indexOffset,
		ByteLength: indexCount * 2,
		Target:     intPtr(34963), // ELEMENT_ARRAY_BUFFER
	})

	// Add accessors
	positionAccessor := len(asset.Accessors)
	boxHalf := 1.25 // boxSize / 2 = 2.5 / 2
	asset.Accessors = append(asset.Accessors, gltf.Accessor{
		BufferView:    &positionBufferView,
		ByteOffset:    0,
		ComponentType: 5126, // FLOAT
		Count:         vertexCount,
		Type:          "VEC3",
		Min:           []float64{-boxHalf, -boxHalf, -boxHalf},
		Max:           []float64{boxHalf, boxHalf, boxHalf},
	})

	indicesAccessor := len(asset.Accessors)
	asset.Accessors = append(asset.Accessors, gltf.Accessor{
		BufferView:    &indicesBufferView,
		ByteOffset:    0,
		ComponentType: 5123, // UNSIGNED_SHORT
		Count:         indexCount,
		Type:          "SCALAR",
	})

	// Update mesh primitive with correct accessor indices
	mesh.Primitives[0].Attributes["POSITION"] = positionAccessor
	mesh.Primitives[0].Indices = &indicesAccessor
	mesh.Primitives[0].Material = &materialIdx

	// Add mesh and material
	asset.Meshes = append(asset.Meshes, mesh)
	asset.Materials = append(asset.Materials, material)

	// Create combined text label for class (name + attributes + operations)

	// IMPORTANT: Text position must always be at the cube center (Z=0)
	// This ensures text is centered on the cube, not offset in front of it.
	// Text and cube positions are kept identical for proper alignment.
	textZ := 0.0

	// Build combined text: class name, attributes, and operations
	var combinedText string
	combinedText = class.Name

	// Add URL if present
	if class.URL != "" {
		combinedText += "\n" + class.URL
	}

	// Add separator line
	combinedText += "\n---"

	// Add attributes
	if len(class.Attributes) > 0 {
		for _, attr := range class.Attributes {
			combinedText += "\n" + attr.Name + ": " + attr.Type
		}
	}

	// Add separator line
	combinedText += "\n---"

	// Add operations
	if len(class.Operations) > 0 {
		for _, op := range class.Operations {
			combinedText += "\n" + op.Name + "(): " + op.ReturnType
		}
	}

	// Add class node (wireframe cube)
	classNodeIdx := len(asset.Nodes)
	asset.Nodes = append(asset.Nodes, gltf.Node{
		Name:        class.Name,
		Mesh:        &meshIdx,
		Translation: []float64{position[0], position[1], position[2]},
		Extras: map[string]any{
			"extuml": map[string]any{
				"type":       "class",
				"id":         class.ID,
				"attributes": len(class.Attributes),
				"operations": len(class.Operations),
			},
		},
	})

	// Add text label as independent node (not a child) with world coordinates
	// Text position = class position + Z offset
	textWorldPos := [3]float64{position[0], position[1], position[2] + textZ}
	u.addTextLabel(combinedText, textWorldPos, true, class.URL, asset)

	*nodeIndex = classNodeIdx + 2 // Class node + text node
}

func (u *generateUsecaseImpl) addInterfaceToScene(iface extuml.Interface, position [3]float64, asset *gltf.GLTFAsset, nodeIndex *int) {
	mesh, material, bufferData := u.geomGen.GenerateInterfaceWireframe(iface, position)

	meshIdx := len(asset.Meshes)
	materialIdx := len(asset.Materials)
	bufferIdx := len(asset.Buffers)

	asset.Buffers = append(asset.Buffers, gltf.Buffer{
		ByteLength: len(bufferData),
		URI:        u.geomGen.CreateBufferURI(bufferData),
	})

	// Calculate vertex and index count from buffer
	vertexCount, indexCount, vertexBytes, indexOffset := u.calculateBufferLayout(bufferData)

	positionBufferView := len(asset.BufferViews)
	asset.BufferViews = append(asset.BufferViews, gltf.BufferView{
		Buffer:     bufferIdx,
		ByteOffset: 0,
		ByteLength: vertexBytes,
		Target:     intPtr(34962),
	})

	indicesBufferView := len(asset.BufferViews)
	asset.BufferViews = append(asset.BufferViews, gltf.BufferView{
		Buffer:     bufferIdx,
		ByteOffset: indexOffset,
		ByteLength: indexCount * 2,
		Target:     intPtr(34963),
	})

	positionAccessor := len(asset.Accessors)
	asset.Accessors = append(asset.Accessors, gltf.Accessor{
		BufferView:    &positionBufferView,
		ByteOffset:    0,
		ComponentType: 5126,
		Count:         vertexCount,
		Type:          "VEC3",
		Min:           []float64{-1, -0.4, -0.25},
		Max:           []float64{1, 0.4, 0.25},
	})

	indicesAccessor := len(asset.Accessors)
	asset.Accessors = append(asset.Accessors, gltf.Accessor{
		BufferView:    &indicesBufferView,
		ByteOffset:    0,
		ComponentType: 5123,
		Count:         indexCount,
		Type:          "SCALAR",
	})

	mesh.Primitives[0].Attributes["POSITION"] = positionAccessor
	mesh.Primitives[0].Indices = &indicesAccessor
	mesh.Primitives[0].Material = &materialIdx

	asset.Meshes = append(asset.Meshes, mesh)
	asset.Materials = append(asset.Materials, material)

	asset.Nodes = append(asset.Nodes, gltf.Node{
		Name:        iface.Name,
		Mesh:        &meshIdx,
		Translation: []float64{position[0], position[1], position[2]},
		Extras: map[string]any{
			"extuml": map[string]any{
				"type":       "interface",
				"id":         iface.ID,
				"operations": len(iface.Operations),
			},
		},
	})

	*nodeIndex++
}

func (u *generateUsecaseImpl) addEnumToScene(enum extuml.Enum, position [3]float64, asset *gltf.GLTFAsset, nodeIndex *int) {
	mesh, material, bufferData := u.geomGen.GenerateEnumWireframe(enum, position)

	meshIdx := len(asset.Meshes)
	materialIdx := len(asset.Materials)
	bufferIdx := len(asset.Buffers)

	asset.Buffers = append(asset.Buffers, gltf.Buffer{
		ByteLength: len(bufferData),
		URI:        u.geomGen.CreateBufferURI(bufferData),
	})

	// Calculate vertex and index count from buffer
	vertexCount, indexCount, vertexBytes, indexOffset := u.calculateBufferLayout(bufferData)

	positionBufferView := len(asset.BufferViews)
	asset.BufferViews = append(asset.BufferViews, gltf.BufferView{
		Buffer:     bufferIdx,
		ByteOffset: 0,
		ByteLength: vertexBytes,
		Target:     intPtr(34962),
	})

	indicesBufferView := len(asset.BufferViews)
	asset.BufferViews = append(asset.BufferViews, gltf.BufferView{
		Buffer:     bufferIdx,
		ByteOffset: indexOffset,
		ByteLength: indexCount * 2,
		Target:     intPtr(34963),
	})

	positionAccessor := len(asset.Accessors)
	asset.Accessors = append(asset.Accessors, gltf.Accessor{
		BufferView:    &positionBufferView,
		ByteOffset:    0,
		ComponentType: 5126,
		Count:         vertexCount,
		Type:          "VEC3",
		Min:           []float64{-0.75, -0.3, -0.2},
		Max:           []float64{0.75, 0.3, 0.2},
	})

	indicesAccessor := len(asset.Accessors)
	asset.Accessors = append(asset.Accessors, gltf.Accessor{
		BufferView:    &indicesBufferView,
		ByteOffset:    0,
		ComponentType: 5123,
		Count:         indexCount,
		Type:          "SCALAR",
	})

	mesh.Primitives[0].Attributes["POSITION"] = positionAccessor
	mesh.Primitives[0].Indices = &indicesAccessor
	mesh.Primitives[0].Material = &materialIdx

	asset.Meshes = append(asset.Meshes, mesh)
	asset.Materials = append(asset.Materials, material)

	asset.Nodes = append(asset.Nodes, gltf.Node{
		Name:        enum.Name,
		Mesh:        &meshIdx,
		Translation: []float64{position[0], position[1], position[2]},
		Extras: map[string]any{
			"extuml": map[string]any{
				"type":     "enum",
				"id":       enum.ID,
				"literals": len(enum.Literals),
			},
		},
	})

	*nodeIndex++
}

// addTextLabel creates a billboard text label node and returns its index
func (u *generateUsecaseImpl) addTextLabel(text string, position [3]float64, billboard bool, url string, asset *gltf.GLTFAsset) int {
	mesh, bufferData := u.textGen.GenerateTextQuad(text, position)

	meshIdx := len(asset.Meshes)
	bufferIdx := len(asset.Buffers)

	// Add buffer
	asset.Buffers = append(asset.Buffers, gltf.Buffer{
		ByteLength: len(bufferData),
		URI:        u.textGen.CreateTextBufferURI(bufferData),
	})

	// Calculate buffer layout
	vertexCount := 4 // 4 vertices for quad
	indexCount := 6  // 6 indices for 2 triangles
	stride := 5 * 4  // 5 floats per vertex (3 for position, 2 for UV) * 4 bytes
	vertexBytes := vertexCount * stride
	indexOffset := vertexBytes
	if indexOffset%4 != 0 {
		indexOffset += 4 - (indexOffset % 4)
	}

	// Add buffer view for vertices (interleaved position + UV)
	vertexBufferView := len(asset.BufferViews)
	asset.BufferViews = append(asset.BufferViews, gltf.BufferView{
		Buffer:     bufferIdx,
		ByteOffset: 0,
		ByteLength: vertexBytes,
		ByteStride: intPtr(20),    // 5 floats * 4 bytes = 20 bytes per vertex
		Target:     intPtr(34962), // ARRAY_BUFFER
	})

	// Add buffer view for indices
	indicesBufferView := len(asset.BufferViews)
	asset.BufferViews = append(asset.BufferViews, gltf.BufferView{
		Buffer:     bufferIdx,
		ByteOffset: indexOffset,
		ByteLength: indexCount * 2,
		Target:     intPtr(34963), // ELEMENT_ARRAY_BUFFER
	})

	// Add accessor for positions
	positionAccessor := len(asset.Accessors)

	// Calculate min/max for the text quad (must match text_geometry.go dimensions)
	// Width and height are 2.4, so half is 1.2
	w := 1.2
	h := 1.2

	asset.Accessors = append(asset.Accessors, gltf.Accessor{
		BufferView:    &vertexBufferView,
		ByteOffset:    0,
		ComponentType: 5126, // FLOAT
		Count:         vertexCount,
		Type:          "VEC3",
		Min:           []float64{-w, -h, 0},
		Max:           []float64{w, h, 0},
	})

	// Add accessor for UV coordinates
	uvAccessor := len(asset.Accessors)
	asset.Accessors = append(asset.Accessors, gltf.Accessor{
		BufferView:    &vertexBufferView,
		ByteOffset:    12,   // 3 floats * 4 bytes offset for UV
		ComponentType: 5126, // FLOAT
		Count:         vertexCount,
		Type:          "VEC2",
	})

	// Add accessor for indices
	indicesAccessor := len(asset.Accessors)
	asset.Accessors = append(asset.Accessors, gltf.Accessor{
		BufferView:    &indicesBufferView,
		ByteOffset:    0,
		ComponentType: 5123, // UNSIGNED_SHORT
		Count:         indexCount,
		Type:          "SCALAR",
	})

	// Create simple placeholder material (texture will be dynamically generated in viewer)
	// Use material index as name to avoid issues with special characters in text
	materialIdx := len(asset.Materials)
	materialName := fmt.Sprintf("text_material_%d", materialIdx)
	asset.Materials = append(asset.Materials, gltf.Material{
		Name: materialName,
		PbrMetallicRoughness: &gltf.PbrMetallicRoughness{
			BaseColorFactor: []float64{1, 1, 1, 1}, // White placeholder
			MetallicFactor:  0,
			RoughnessFactor: 1,
		},
		DoubleSided: true,
	})

	// Update mesh
	mesh.Primitives[0].Attributes["POSITION"] = positionAccessor
	mesh.Primitives[0].Attributes["TEXCOORD_0"] = uvAccessor
	mesh.Primitives[0].Indices = &indicesAccessor
	mesh.Primitives[0].Material = &materialIdx

	// Add mesh
	asset.Meshes = append(asset.Meshes, mesh)

	// Create node with billboard extras
	nodeIdx := len(asset.Nodes)
	extras := map[string]any{
		"extuml": map[string]any{
			"type": "text",
			"text": text,
		},
	}
	if billboard {
		extras["billboard"] = true
	}
	if url != "" {
		extras["url"] = url
	}

	// Use simple node name based on index to avoid issues with special characters
	// The actual text content is stored in extras.extuml.text
	nodeName := fmt.Sprintf("text_node_%d", nodeIdx)

	asset.Nodes = append(asset.Nodes, gltf.Node{
		Name:        nodeName,
		Mesh:        &meshIdx,
		Translation: []float64{position[0], position[1], position[2]},
		Extras:      extras,
	})

	return nodeIdx
}

// calculateBufferLayout calculates vertex and index counts from buffer data
func (u *generateUsecaseImpl) calculateBufferLayout(bufferData []byte) (vertexCount, indexCount, vertexBytes, indexOffset int) {
	// Find configuration where index count is reasonable for wireframe
	// Prefer configurations where ratio is close to 2.5 (ideal for wireframe boxes with dividers)
	bestVC, bestIC, bestVB, bestIO := 0, 0, 0, 0
	bestRatio := 0.0
	targetRatio := 2.5

	for vc := 1; vc <= 100; vc++ {
		vb := vc * 3 * 4 // 3 floats per vertex * 4 bytes
		io := vb
		if io%4 != 0 {
			io += 4 - (io % 4)
		}
		if io >= len(bufferData) {
			break
		}
		remaining := len(bufferData) - io
		if remaining%2 == 0 && remaining > 0 {
			ic := remaining / 2
			ratio := float64(ic) / float64(vc)
			// For wireframes: index count should be 2-6x vertex count
			if ic >= vc*2 && ic <= vc*6 {
				// Prefer ratio closer to target (2.5)
				ratioDiff := ratio - targetRatio
				if ratioDiff < 0 {
					ratioDiff = -ratioDiff
				}
				currentBestDiff := bestRatio - targetRatio
				if currentBestDiff < 0 {
					currentBestDiff = -currentBestDiff
				}
				if bestVC == 0 || ratioDiff < currentBestDiff {
					bestVC, bestIC, bestVB, bestIO = vc, ic, vb, io
					bestRatio = ratio
				}
			}
		}
	}

	// If no good ratio found, use maximum valid vertex count
	if bestVC == 0 {
		for vc := 1; vc <= 100; vc++ {
			vb := vc * 3 * 4
			io := vb
			if io%4 != 0 {
				io += 4 - (io % 4)
			}
			if io >= len(bufferData) {
				break
			}
			remaining := len(bufferData) - io
			if remaining%2 == 0 && remaining > 0 {
				ic := remaining / 2
				if ic >= 2 && vc > bestVC {
					bestVC, bestIC, bestVB, bestIO = vc, ic, vb, io
				}
			}
		}
	}

	return bestVC, bestIC, bestVB, bestIO
}

// calculateSceneBounds calculates the bounding box of all nodes and recommends camera settings
func (u *generateUsecaseImpl) calculateSceneBounds(asset *gltf.GLTFAsset) map[string]any {
	minX, minY, minZ := 1e10, 1e10, 1e10
	maxX, maxY, maxZ := -1e10, -1e10, -1e10

	// Iterate through all nodes to find bounds
	for _, node := range asset.Nodes {
		if len(node.Translation) == 3 {
			x, y, z := node.Translation[0], node.Translation[1], node.Translation[2]

			// Estimate node size (box is 2x2x0.5)
			nodeWidth, nodeHeight, nodeDepth := 2.0, 2.0, 0.5

			if x-nodeWidth < minX {
				minX = x - nodeWidth
			}
			if x+nodeWidth > maxX {
				maxX = x + nodeWidth
			}
			if y-nodeHeight < minY {
				minY = y - nodeHeight
			}
			if y+nodeHeight > maxY {
				maxY = y + nodeHeight
			}
			if z-nodeDepth < minZ {
				minZ = z - nodeDepth
			}
			if z+nodeDepth > maxZ {
				maxZ = z + nodeDepth
			}
		}
	}

	// Calculate center and size
	centerX := (minX + maxX) / 2
	centerY := (minY + maxY) / 2
	centerZ := (minZ + maxZ) / 2

	sizeX := maxX - minX
	sizeY := maxY - minY
	sizeZ := maxZ - minZ

	// Calculate optimal camera distance (larger of width/height + depth)
	maxDimension := sizeX
	if sizeY > maxDimension {
		maxDimension = sizeY
	}

	// Camera distance should be ~2-3x the max dimension to fit everything
	cameraDistance := maxDimension * 2.5
	if cameraDistance < 3.0 {
		cameraDistance = 3.0 // Minimum distance
	}

	return map[string]any{
		"bounds": map[string]any{
			"min":    []float64{minX, minY, minZ},
			"max":    []float64{maxX, maxY, maxZ},
			"center": []float64{centerX, centerY, centerZ},
			"size":   []float64{sizeX, sizeY, sizeZ},
		},
		"recommended": map[string]any{
			"distance": cameraDistance,
			"target":   []float64{centerX, centerY, centerZ},
			"orbit":    []float64{45, 55, cameraDistance}, // theta, phi, radius
		},
	}
}
