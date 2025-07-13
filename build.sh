#!/bin/bash

# NVS CLI Cross-Platform Build Script
# Builds the NVS CLI tool for Windows, Linux, and macOS

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Version and build info
VERSION=${VERSION:-"dev"}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S_UTC')
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${BLUE}🚀 Building NVS CLI v${VERSION}${NC}"
echo -e "${BLUE}📅 Build Time: ${BUILD_TIME}${NC}"
echo -e "${BLUE}🔗 Commit: ${COMMIT_HASH}${NC}"
echo ""

# Create build directory
BUILD_DIR="build"
mkdir -p $BUILD_DIR

# Build flags
LDFLAGS="-X 'main.Version=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.CommitHash=${COMMIT_HASH}'"

# Function to build for a specific platform
build_for_platform() {
    local GOOS=$1
    local GOARCH=$2
    local EXTENSION=$3
    local OUTPUT_NAME="nvs${EXTENSION}"
    
    echo -e "${YELLOW}🔨 Building for ${GOOS}/${GOARCH}...${NC}"
    
    # Set environment variables
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    export CGO_ENABLED=0
    
    # Build the binary
    go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/$OUTPUT_NAME" .
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Successfully built: $BUILD_DIR/$OUTPUT_NAME${NC}"
        
        # Get file size
        if command -v stat >/dev/null 2>&1; then
            SIZE=$(stat -f%z "$BUILD_DIR/$OUTPUT_NAME" 2>/dev/null || stat -c%s "$BUILD_DIR/$OUTPUT_NAME" 2>/dev/null || echo "unknown")
            echo -e "${BLUE}📦 File size: ${SIZE} bytes${NC}"
        fi
    else
        echo -e "${RED}❌ Failed to build for ${GOOS}/${GOARCH}${NC}"
        return 1
    fi
    
    echo ""
}

# Function to create checksums
create_checksums() {
    echo -e "${YELLOW}🔍 Creating checksums...${NC}"
    cd $BUILD_DIR
    
    # Create SHA256 checksums
    if command -v shasum >/dev/null 2>&1; then
        shasum -a 256 nvs* > checksums.sha256
    elif command -v sha256sum >/dev/null 2>&1; then
        sha256sum nvs* > checksums.sha256
    else
        echo -e "${RED}⚠️  No checksum tool found${NC}"
        return
    fi
    
    echo -e "${GREEN}✅ Checksums created: checksums.sha256${NC}"
    cd ..
    echo ""
}

# Function to create release archive
create_archive() {
    local GOOS=$1
    local GOARCH=$2
    local EXTENSION=$3
    local OUTPUT_NAME="nvs${EXTENSION}"
    
    echo -e "${YELLOW}📦 Creating archive for ${GOOS}/${GOARCH}...${NC}"
    
    cd $BUILD_DIR
    
    # Create archive based on platform
    if [ "$GOOS" = "windows" ]; then
        zip "nvs-${VERSION}-${GOOS}-${GOARCH}.zip" "$OUTPUT_NAME" checksums.sha256 >/dev/null 2>&1 || echo "zip not available"
    else
        tar -czf "nvs-${VERSION}-${GOOS}-${GOARCH}.tar.gz" "$OUTPUT_NAME" checksums.sha256 >/dev/null 2>&1 || echo "tar not available"
    fi
    
    cd ..
    echo -e "${GREEN}✅ Archive created${NC}"
    echo ""
}

# Main build process
echo -e "${BLUE}📋 Building for all supported platforms...${NC}"
echo ""

# Build for different platforms
build_for_platform "linux" "amd64" ""
build_for_platform "linux" "arm64" ""
build_for_platform "windows" "amd64" ".exe"
build_for_platform "windows" "arm64" ".exe"
build_for_platform "darwin" "amd64" ""
build_for_platform "darwin" "arm64" ""

# Create checksums
create_checksums

# Create archives
echo -e "${BLUE}📦 Creating release archives...${NC}"
create_archive "linux" "amd64" ""
create_archive "linux" "arm64" ""
create_archive "windows" "amd64" ".exe"
create_archive "windows" "arm64" ".exe"
create_archive "darwin" "amd64" ""
create_archive "darwin" "arm64" ""

# Summary
echo -e "${GREEN}🎉 Build completed successfully!${NC}"
echo ""
echo -e "${BLUE}📁 Build artifacts:${NC}"
ls -la $BUILD_DIR/
echo ""
echo -e "${BLUE}📋 Supported platforms:${NC}"
echo "  • Linux (amd64, arm64)"
echo "  • Windows (amd64, arm64)"
echo "  • macOS (amd64, arm64)"
echo ""
echo -e "${BLUE}💡 Usage:${NC}"
echo "  • Linux/macOS: ./nvs"
echo "  • Windows: nvs.exe"
echo ""
echo -e "${YELLOW}🔧 To install globally:${NC}"
echo "  • Copy the appropriate binary to a directory in your PATH"
echo "  • Or use: go install github.com/nevilsoft/nvs@latest" 