.PHONY: build test clean run generate serve help

# Default target
help:
	@echo "Available targets:"
	@echo "  build      - Build the extuml binary"
	@echo "  test       - Run all tests"
	@echo "  clean      - Remove build artifacts"
	@echo "  run        - Build and run extuml"
	@echo "  generate   - Generate glTF from sample.extuml"
	@echo "  serve      - Start HTTP server for viewer (development)"
	@echo "  help       - Show this help message"

# Build the binary
build:
	@echo "Building extuml..."
	@go build -o .bin/extuml ./cmd
	@echo "✅ Build complete: .bin/extuml"

# Run tests
test:
	@echo "Running tests..."
	@go test ./test/... -v

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf .bin
	@rm -f etc/*.gl
	@echo "✅ Clean complete"

# Build and run
run: build
	@echo "Running extuml..."
	@./.bin/extuml --help

# Generate glTF from sample
generate: build
	@echo "Generating glTF and HTML from sample.extuml..."
	@./.bin/extuml generate -e etc/sample.extuml -o etc/output.gl --html-output etc/index.html
	@echo "✅ Generated: etc/output.gl"
	@echo "✅ Generated: etc/index.html"

# Start HTTP server for viewer (development mode)
serve:
	@echo "🚀 Starting HTTP server at http://localhost:8000"
	@echo "📂 Serving etc/ directory"
	@echo "📄 Open http://localhost:8000/index.html in your browser"
	@echo ""
	@echo "Press Ctrl+C to stop"
	@cd etc && python3 -m http.server 8000

# Development workflow: generate and serve
dev: generate serve
