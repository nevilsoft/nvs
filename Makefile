# NVS CLI Makefile
# Cross-platform build and development tasks

# Variables
BINARY_NAME=nvs
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')
COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DIR=build
LDFLAGS=-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.CommitHash=$(COMMIT_HASH)'

# Detect OS
ifeq ($(OS),Windows_NT)
	BINARY_NAME := $(BINARY_NAME).exe
	RM := del /Q
	MKDIR := mkdir
	RMDIR := rmdir /S /Q
else
	RM := rm -f
	MKDIR := mkdir -p
	RMDIR := rm -rf
endif

# Default target
.PHONY: all
all: clean build

# Build for current platform
.PHONY: build
build:
	@echo "🔨 Building $(BINARY_NAME) for $(shell go env GOOS)/$(shell go env GOARCH)..."
	@go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .
	@echo "✅ Build completed: $(BINARY_NAME)"

# Build for all platforms
.PHONY: build-all
build-all: clean
	@echo "🚀 Building for all supported platforms..."
	@$(MKDIR) $(BUILD_DIR)
	@$(MAKE) build-linux-amd64
	@$(MAKE) build-linux-arm64
	@$(MAKE) build-windows-amd64
	@$(MAKE) build-windows-arm64
	@$(MAKE) build-darwin-amd64
	@$(MAKE) build-darwin-arm64
	@$(MAKE) checksums
	@echo "🎉 All builds completed!"

# Build for specific platforms
.PHONY: build-linux-amd64
build-linux-amd64:
	@echo "🔨 Building for linux/amd64..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/nvs .
	@echo "✅ Built: $(BUILD_DIR)/nvs"

.PHONY: build-linux-arm64
build-linux-arm64:
	@echo "🔨 Building for linux/arm64..."
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/nvs-linux-arm64 .
	@echo "✅ Built: $(BUILD_DIR)/nvs-linux-arm64"

.PHONY: build-windows-amd64
build-windows-amd64:
	@echo "🔨 Building for windows/amd64..."
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/nvs.exe .
	@echo "✅ Built: $(BUILD_DIR)/nvs.exe"

.PHONY: build-windows-arm64
build-windows-arm64:
	@echo "🔨 Building for windows/arm64..."
	@GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/nvs-windows-arm64.exe .
	@echo "✅ Built: $(BUILD_DIR)/nvs-windows-arm64.exe"

.PHONY: build-darwin-amd64
build-darwin-amd64:
	@echo "🔨 Building for darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/nvs-darwin-amd64 .
	@echo "✅ Built: $(BUILD_DIR)/nvs-darwin-amd64"

.PHONY: build-darwin-arm64
build-darwin-arm64:
	@echo "🔨 Building for darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/nvs-darwin-arm64 .
	@echo "✅ Built: $(BUILD_DIR)/nvs-darwin-arm64"

# Create checksums
.PHONY: checksums
checksums:
	@echo "🔍 Creating checksums..."
	@cd $(BUILD_DIR) && \
	if command -v shasum >/dev/null 2>&1; then \
		shasum -a 256 nvs* > checksums.sha256; \
	elif command -v sha256sum >/dev/null 2>&1; then \
		sha256sum nvs* > checksums.sha256; \
	else \
		echo "⚠️ No checksum tool found"; \
	fi
	@echo "✅ Checksums created"

# Create release archives
.PHONY: release
release: build-all
	@echo "📦 Creating release archives..."
	@cd $(BUILD_DIR) && \
	for file in nvs*; do \
		if [ -f "$$file" ]; then \
			platform=$$(echo $$file | sed 's/nvs//' | sed 's/\.exe//'); \
			if [ "$$platform" = "" ]; then \
				platform="linux-amd64"; \
			fi; \
			if echo "$$file" | grep -q "\.exe"; then \
				zip "nvs-$(VERSION)-$$platform.zip" "$$file" checksums.sha256 2>/dev/null || echo "zip not available"; \
			else \
				tar -czf "nvs-$(VERSION)-$$platform.tar.gz" "$$file" checksums.sha256 2>/dev/null || echo "tar not available"; \
			fi; \
			echo "✅ Created archive for $$platform"; \
		fi; \
	done
	@echo "🎉 Release archives created!"

# Development tasks
.PHONY: dev
dev:
	@echo "🚀 Starting development mode..."
	@go run . dev

.PHONY: test
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

.PHONY: test-coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: coverage.html"

# Code quality
.PHONY: lint
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️ golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

.PHONY: fmt
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...

.PHONY: vet
vet:
	@echo "🔍 Running go vet..."
	@go vet ./...

# Dependencies
.PHONY: deps
deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy

.PHONY: deps-update
deps-update:
	@echo "📦 Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# Cleanup
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	@$(RMDIR) $(BUILD_DIR) 2>/dev/null || true
	@$(RM) $(BINARY_NAME) 2>/dev/null || true
	@$(RM) coverage.out coverage.html 2>/dev/null || true
	@echo "✅ Clean completed"

# Install
.PHONY: install
install: build
	@echo "📦 Installing $(BINARY_NAME)..."
	@go install .
	@echo "✅ Installation completed"

# Uninstall
.PHONY: uninstall
uninstall:
	@echo "🗑️ Uninstalling $(BINARY_NAME)..."
	@go clean -i .
	@echo "✅ Uninstallation completed"

# Help
.PHONY: help
help:
	@echo "NVS CLI Makefile - Available targets:"
	@echo ""
	@echo "Build targets:"
	@echo "  build          - Build for current platform"
	@echo "  build-all      - Build for all supported platforms"
	@echo "  build-linux-*  - Build for specific Linux architecture"
	@echo "  build-windows-* - Build for specific Windows architecture"
	@echo "  build-darwin-* - Build for specific macOS architecture"
	@echo "  release        - Build all platforms and create archives"
	@echo ""
	@echo "Development targets:"
	@echo "  dev            - Start development mode"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo ""
	@echo "Code quality targets:"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo ""
	@echo "Dependency targets:"
	@echo "  deps           - Install dependencies"
	@echo "  deps-update    - Update dependencies"
	@echo ""
	@echo "Utility targets:"
	@echo "  clean          - Clean build artifacts"
	@echo "  install        - Install binary"
	@echo "  uninstall      - Uninstall binary"
	@echo "  help           - Show this help"
	@echo ""
	@echo "Environment variables:"
	@echo "  VERSION        - Set build version (default: dev)" 