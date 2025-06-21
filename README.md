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

## 📋 Prerequisites

- Go 1.24.0 or higher
- Git (for user configuration)

## 🏃‍♂️ Quick Start

### Installation

```bash
go install https://github.com/nevilsoft/nvscli@latest
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
│       ├── routes/                # Route definitions
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

© 2025 Nevilsoft Part., Ltd. All Rights Reserved.

This project contains confidential and proprietary information. Unauthorized copying, modification, distribution, or use is strictly prohibited.

## 📞 Support

For support and questions, please contact the development team at Nevilsoft Part., Ltd.

---

**Note**: This project contains confidential business information and is restricted to authorized personnel only. Violation of these terms may result in disciplinary action and legal proceedings under the Computer Crime Act B.E. 2560 (Sections 7, 9, 10) and other applicable laws. 