# extuml

A 3D UML diagram generator that transforms Mermaid-like DSL into glTF 2.0 format.

## Quick Start

### Using Makefile (Recommended)

```bash
# Show all available commands
make help

# Build the binary
make build

# Generate glTF from sample
make generate

# Start development server for viewer
make serve
# Then open http://localhost:8000 in your browser

# Development workflow (generate + serve)
make dev
```

### Manual Build

```bash
go build -o .bin/extuml ./cmd
```

### Generate glTF from extuml DSL

```bash
.bin/extuml generate --extuml etc/sample.extuml --output etc/output.gl
# or with short flags
.bin/extuml generate -e etc/sample.extuml -o etc/output.gl
```

### View Generated 3D Model

The viewer automatically loads and displays `etc/output.gl`:

```bash
# Start the development server
make serve
# or manually
cd .dist
python3 -m http.server 8000
# Then open http://localhost:8000 in your browser
```

The viewer features:
- � **Auto-reload**: Refreshes every 2 seconds during development
- � **Automatic loading**: Always displays .etc/output.gl
- 🎮 **Interactive controls**: Rotate, zoom, and pan with mouse
- 📊 **Metadata display**: Shows model info and glTF structure
- 📄 **JSON preview**: View complete glTF asset data
- Built with [model-viewer](https://modelviewer.dev/) by Google

## Project Structure

```
extuml/
├── .bin/                   # Built binaries
├── .dist/                  # Web viewer for glTF files
│   └── index.html         # model-viewer based 3D viewer
├── .etc/                   # Sample files
│   ├── sample.extuml      # Sample extuml DSL
│   └── output.gl          # Generated glTF output
├── cmd/                   # CLI entry point
│   └── main.go
├── pkg/                   # Application packages
│   ├── command/          # CLI commands (cobra-based)
│   ├── config/           # Dependency injection
│   ├── controller/       # Command handlers
│   ├── model/            # Data structures (extuml/, gltf/)
│   ├── repository/       # File I/O
│   └── usecase/          # Business logic
├── test/                  # Integration tests
├── go.mod
└── README.md
```

## Architecture

Following clean architecture pattern:

1. **Command Layer**: CLI command definitions
2. **Controller Layer**: Request validation and coordination
3. **UseCase Layer**: Core business logic
4. **Repository Layer**: Data access (file I/O)
5. **Model Layer**: Data structures
6. **Config Layer**: Dependency injection

## Development

### Run with sample

```bash
go run ./cmd generate --extuml .etc/sample.extuml --output output.gl
```

### Test

```bash
go test ./...
```

## License

MIT
# extuml
