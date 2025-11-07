package usecase

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// TextTextureGenerator generates PNG texture images for text labels
type TextTextureGenerator struct{}

// NewTextTextureGenerator creates a new text texture generator
func NewTextTextureGenerator() *TextTextureGenerator {
	return &TextTextureGenerator{}
}

// TextureResult contains the generated texture and its dimensions
type TextureResult struct {
	DataURI string
	Width   int
	Height  int
}

// GenerateTextTexture generates a PNG texture with the given text
func (g *TextTextureGenerator) GenerateTextTexture(text string, textColor color.RGBA, bgColor color.RGBA) (*TextureResult, error) {
	// Calculate texture dimensions based on text length
	charWidth := 8   // basicfont width
	charHeight := 13 // basicfont height
	padding := 12    // Increased padding for better readability

	width := len(text)*charWidth + padding*2
	height := charHeight + padding*2

	// Ensure power-of-2 dimensions for better GPU compatibility
	width = nextPowerOf2(width)
	height = nextPowerOf2(height)

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Draw rounded rectangle background for text (centered)
	textWidth := len(text) * charWidth
	bgRect := image.Rect(
		padding-4,
		padding-4,
		padding+textWidth+4,
		padding+charHeight+4,
	)
	draw.Draw(img, bgRect, &image.Uniform{bgColor}, image.Point{}, draw.Over)

	// Draw text
	point := fixed.Point26_6{
		X: fixed.Int26_6(padding * 64),
		Y: fixed.Int26_6((padding + charHeight) * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{textColor},
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	// Create data URI
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	dataURI := "data:image/png;base64," + encoded

	return &TextureResult{
		DataURI: dataURI,
		Width:   width,
		Height:  height,
	}, nil
}

// nextPowerOf2 returns the next power of 2 greater than or equal to n
func nextPowerOf2(n int) int {
	if n <= 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}
