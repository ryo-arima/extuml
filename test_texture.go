package extuml
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/extuml/extuml/pkg/usecase"
)

func main() {
	gen := usecase.NewTextTextureGenerator()
	
	textColor := color.RGBA{R: 25, G: 77, B: 102, A: 255} // Dark blue
	bgColor := color.RGBA{R: 255, G: 255, B: 255, A: 230} // White
	
	result, err := gen.GenerateTextTexture("Person", textColor, bgColor)
	if err != nil {
		panic(err)
	}
	
	// Extract base64 data
	parts := strings.Split(result.DataURI, ",")
	if len(parts) != 2 {
		panic("Invalid data URI")
	}
	
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		panic(err)
	}
	
	// Save to file
	err = os.WriteFile("test_person.png", data, 0644)
	if err != nil {
		panic(err)
	}
	
	fmt.Println("Generated test_person.png")
	fmt.Printf("Size: %dx%d\n", result.Width, result.Height)
	fmt.Printf("Data size: %d bytes\n", len(data))
}
