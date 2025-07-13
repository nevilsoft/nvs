# NVS CLI

A powerful command-line interface tool for creating and managing Go projects with modern development practices, built-in obfuscation, and comprehensive project scaffolding.

## 🚀 Features

- **Project Scaffolding**: Interactive project initialization with template system
- **Code Obfuscation**: Built-in Garble integration for secure builds
- **Development Tools**: Hot reload with Air, automatic dependency management
- **Multi-platform Support**: Cross-platform builds for various architectures
- **Template System**: Embedded templates with dynamic rendering
- **Security Features**: SHA256 verification and secure execution
- **Interactive CLI**: User-friendly prompts and confirmations
- **Route Generation**: Auto-generates route files with `_route.go` suffix (e.g., `user_route.go`, `product_route.go`)

## 📋 Prerequisites

- Go 1.24.0 or higher
- Git (for user configuration)

## 🌍 Cross-Platform Support

NVS CLI is fully compatible with Windows, Linux, and macOS. The tool automatically adapts to your operating system:

### Windows Support
- **Path Handling**: Uses Windows-style path separators (`\`)
- **File Permissions**: Handles Windows file permission system
- **Executable Detection**: Automatically adds `.exe` extension when needed
- **Git Integration**: Works with Git for Windows
- **Environment Variables**: Supports Windows environment variable syntax

### Linux Support
- **Unix Permissions**: Properly handles Unix file permissions (755, 644)
- **Executable Files**: Automatically makes files executable when needed
- **Path Handling**: Uses forward slashes (`/`) for paths
- **Shell Integration**: Compatible with bash, zsh, and other shells

### macOS Support
- **Unix-like System**: Inherits all Linux compatibility features
- **ARM64 Support**: Native support for Apple Silicon (M1/M2) processors
- **Homebrew Integration**: Works seamlessly with Homebrew-installed Go
- **macOS Permissions**: Handles macOS-specific permission requirements

### Platform-Specific Features

| Feature | Windows | Linux | macOS |
|---------|---------|-------|-------|
| Path Separators | `\` | `/` | `/` |
| File Permissions | 666 (files), 666 (dirs) | 644 (files), 755 (dirs) | 644 (files), 755 (dirs) |
| Executable Extension | `.exe` | None | None |
| Git Integration | ✅ | ✅ | ✅ |
| Hot Reload | ✅ | ✅ | ✅ |
| Code Obfuscation | ✅ | ✅ | ✅ |
| Cross-compilation | ✅ | ✅ | ✅ |

### Installation by Platform

#### Windows
```cmd
# Using Go install
go install github.com/nevilsoft/nvs@latest

# Add Go bin to PATH if needed
set PATH=%PATH%;%GOPATH%\bin
```

#### Linux
```bash
# Using Go install
go install github.com/nevilsoft/nvs@latest

# Add to PATH if needed
export PATH=$PATH:$GOPATH/bin
```

#### macOS
```bash
# Using Go install
go install github.com/nevilsoft/nvs@latest

# Add to PATH if needed
export PATH=$PATH:$GOPATH/bin
```

## 🏃‍♂️ Quick Start

### Installation

```bash
go install github.com/nevilsoft/nvs@latest
```

### Basic Usage

```bash
# Show help
nvs --help

# Initialize a new project
nvs init

# Run in development mode
nvs dev

# Build with obfuscation
nvs build

# Start the application
nvs start main
```

## 📖 Commands

### Project Initialization

```bash
# Initialize a new project with interactive prompts
nvs init

# Initialize with custom module name
nvs init --repo github.com/username/project-name
```

**Features:**
- Interactive project name input with defaults
- Automatic Git user detection
- Module name generation
- Template-based project structure
- Confirmation prompts

### Development Mode

```bash
# Start development server with hot reload
nvs dev
```

**Features:**
- Automatic Air installation
- Automatic Swag installation for API docs
- Hot reload with file watching
- Development environment setup

### Build System

```bash
# Build with default settings
nvs build

# Build with custom output name
nvs build -o myapp

# Build for specific platform
nvs build -t linux/amd64

# Build with version
nvs build -v 1.0.0
```

**Features:**
- Automatic Garble installation
- Code obfuscation for security
- Cross-platform builds
- Version embedding
- Build number generation

### Application Start

```bash
# Start in production mode
nvs start main

# Start in development mode
nvs start main -e dev
```

**Features:**
- SHA256 hash verification
- Environment configuration
- Secure execution
- Runner ID generation

### Route Generation

```bash
# Generate route files from controllers (auto-naming with _route.go suffix)
nvs generate routes
```

**Note:**
- Generated route files use the `_route.go` suffix (e.g., `user_route.go`, `product_route.go`).
- Controller names (e.g., `UserController`) are not included in the route file name.
- The CLI will not overwrite `base.go` and will only update the auto-generated section in `SetupRoutes`.

## 🏗️ Project Structure

```
example/
├── main.go                    # Main application entry point
├── go.mod                     # Go module definition
├── .air.toml                  # Air hot reload configuration
├── sqlc.yaml                  # SQLC database code generation config
├── api/                           # API layer templates
│   └── v1/
│       ├── controllers/           # HTTP controllers
│       ├── middleware/            # HTTP middleware
│       ├── routes/                # Route definitions (auto-generated: *_route.go)
│       └── services/              # Business logic services
├── cmd/                           # Command handlers
├── config/                        # Configuration management
│   └── base.go                    # Main configuration struct
├── constants/                     # Application constants
├── db/                            # Database layer
├── cache/                         # Caching layer
├── di/                            # Dependency injection
│   ├── wire.go                    # Wire DI configuration
│   └── wire_gen.go                # Generated Wire code
├── handler/                       # HTTP handlers
├── lang/                          # Internationalization
├── migrations/                    # Database migrations
├── plugin/                        # Plugin system
├── session/                       # Session management
├── shared/                        # Shared utilities
├── types/                         # Type definitions
└── utils/                         # Utility functions
```

## ⚙️ Configuration

### Environment Variables

```bash
# Development environment
ENV=dev

# Production environment  
ENV=prod

# Runner ID (auto-generated)
RUNNER_ID=<sha256-hash>
```

### Build Configuration

The build system supports various flags:

- `-o, --output`: Output binary name
- `-t, --target`: Target platform (e.g., linux/amd64, windows/amd64, darwin/arm64)
- `-v, --version`: Build version

### Development Tools

**Air Configuration** (`.air.toml`):
- Hot reload for Go files
- Template file watching
- Build directory management
- Excluded directories

**Garble Integration**:
- Automatic installation
- Code obfuscation
- Cross-platform builds
- Version embedding

## 🔧 Development

### Cross-Platform Development

NVS CLI supports development across Windows, Linux, and macOS:

#### Using Makefile (Linux/macOS)
```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run linter
make lint

# Clean build artifacts
make clean
```

#### Using Build Scripts
```bash
# Linux/macOS
./build.sh

# Windows
build.bat
```

#### Using Docker
```bash
# Development environment
docker-compose up nvs-dev

# Run tests
docker-compose up nvs-test

# Run linter
docker-compose up nvs-lint

# Production build
docker-compose up nvs-prod
```

### Adding New Commands

1. Create a new command file in `cmd/`
2. Define the command structure using Cobra
3. Add to the root command in `cmd/root.go`
4. Update this README with new command documentation

### Template System

Templates are embedded using Go's `embed` directive:

```go
//go:embed templates/*
var templatesFS embed.FS
```

**Template Variables:**
- `{{.ProjectName}}`: Project name
- `{{.ModuleName}}`: Go module name

### Build Process

1. **Dependency Check**: Verify required tools (Garble, Air, Swag)
2. **Auto-installation**: Install missing dependencies
3. **Template Rendering**: Process embedded templates
4. **Code Generation**: Generate project structure
5. **Build Execution**: Run build commands

### CI/CD Pipeline

The project includes comprehensive CI/CD support:

#### GitHub Actions
- **Automatic Testing**: Runs on every push and pull request
- **Cross-Platform Builds**: Builds for Windows, Linux, and macOS
- **Release Management**: Automatic releases on tag creation
- **Docker Integration**: Multi-architecture Docker images

#### Build Matrix
- **Linux**: amd64, arm64
- **Windows**: amd64
- **macOS**: amd64

#### Docker Support
- **Multi-stage Builds**: Optimized for size and security
- **Non-root User**: Secure runtime environment
- **Multi-architecture**: Support for different CPU architectures

### Route Generation

- Route files are generated with the `_route.go` suffix (e.g., `user_route.go`, `product_route.go`).
- Controller names are not included in the route file name.
- The CLI will not overwrite `base.go` and will only update the auto-generated section in `SetupRoutes`.

## 🛡️ Security Features

### Code Obfuscation

- **Garble Integration**: Automatic code obfuscation
- **Cross-platform**: Secure builds for multiple architectures
- **Version Embedding**: Build version and runner ID injection

### Execution Verification

- **SHA256 Hashing**: File integrity verification
- **Runner ID**: Unique execution identifier
- **Environment Isolation**: Separate dev/prod environments

## 📦 Dependencies

### Core Dependencies

- `github.com/spf13/cobra`: CLI framework
- `github.com/inconshreveable/mousetrap`: Windows compatibility
- `github.com/spf13/pflag`: Flag parsing

### Development Tools

- `github.com/cosmtrek/air`: Hot reload
- `github.com/swaggo/swag/cmd/swag`: API documentation
- `mvdan.cc/garble`: Code obfuscation

## 🚀 Example Workflow

```bash
# 1. Initialize new project
nvs init
# Follow interactive prompts

# 2. Navigate to project directory
cd my-project

# 3. Start development
nvs dev
# Server starts with hot reload

# 4. Build for production
nvs build -v 1.0.0 -t linux/amd64

# 5. Start production server
nvs start main
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📝 License

© 2025 Nevilsoft Ltd., Part. All Rights Reserved.

This project contains confidential and proprietary information. Unauthorized copying, modification, distribution, or use is strictly prohibited.

## 📞 Support

For support and questions, please contact the development team at Nevilsoft Ltd., Part.

---

**Note**: This project contains confidential business information and is restricted to authorized personnel only. Violation of these terms may result in disciplinary action and legal proceedings under the Computer Crime Act B.E. 2560 (Sections 7, 9, 10) and other applicable laws. 
