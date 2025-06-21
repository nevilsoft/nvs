# NVS CLI

A powerful command-line interface tool built with Go, featuring dependency injection and modern development practices.

## 🚀 Features

- **Dependency Injection**: Built with Google Wire for clean architecture
- **Hot Reload**: Development with Air for automatic rebuilds
- **Modular Design**: Clean separation of concerns with controllers, services, and middleware
- **Template System**: Dynamic code generation with Go templates

## 📋 Prerequisites

- Go 1.21 or higher
- Air (for development hot reload)

## ��️ Installation

### Install Air (for development)
```bash
go install github.com/cosmtrek/air@latest
```

### Clone the repository
```bash
git clone <repository-url>
cd nvs-cli
```

## 🏃‍♂️ Quick Start

### Development Mode
```bash
# Start with hot reload
air
```

### Production Build
```bash
# Build the application
go build -o ./build/main .

# Run the built binary
./build/main
```

## 📖 Usage

### Basic Commands

```bash
# Show help information
nvs --help
nvs -h

# Show version
nvs --version
nvs -v

# Show available commands
nvs list
```

### Project Management

```bash
# Initialize a new project
nvs init my-project

# Create a new module
nvs create module user-service

# Generate API endpoints
nvs generate api user

# Create database migration
nvs migrate create add_users_table
```

### Development Commands

```bash
# Start development server with hot reload
nvs dev

# Run tests
nvs test

# Run tests with coverage
nvs test --coverage

# Build for production
nvs build

# Build for specific platform
nvs build --os=linux --arch=amd64
```

### Database Operations

```bash
# Run database migrations
nvs migrate up

# Rollback last migration
nvs migrate down

# Show migration status
nvs migrate status

# Seed database
nvs seed run
```

### Code Generation

```bash
# Generate new controller
nvs generate controller UserController

# Generate new service
nvs generate service UserService

# Generate new model
nvs generate model User

# Generate CRUD operations
nvs generate crud User
```

### Configuration

```bash
# Show current configuration
nvs config show

# Set configuration value
nvs config set database.host localhost

# Get configuration value
nvs config get database.host

# Edit configuration file
nvs config edit
```

### Utility Commands

```bash
# Format code
nvs fmt

# Lint code
nvs lint

# Clean build artifacts
nvs clean

# Install dependencies
nvs deps install

# Update dependencies
nvs deps update
```

### Advanced Usage

```bash
# Run with custom configuration
nvs --config=./custom-config.yaml

# Enable debug mode
nvs --debug

# Set log level
nvs --log-level=debug

# Run in background
nvs daemon start

# Stop background process
nvs daemon stop
```

## 📁 Project Structure

## 📞 Support

For support and questions, please contact the development team at Nevilsoft Part., Ltd.

---

**Note**: This project contains confidential business information and is restricted to authorized personnel only. Violation of these terms may result in disciplinary action and legal proceedings under the Computer Crime Act B.E. 2560 (Sections 7, 9, 10) and other applicable laws.