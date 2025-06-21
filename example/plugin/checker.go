package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/burapha44/example/shared"
	"github.com/hashicorp/go-plugin"
)

// HashChecker - implementation ของ Checker interface
type HashChecker struct{}

func (h *HashChecker) Check() bool {
	// ตรวจสอบ hash ของ main binary
	hash, err := hashMainBinary()
	if err != nil {
		return false
	}

	expectedHash := os.Getenv("RUNNER_ID")
	return hash == expectedHash
}

func hashMainBinary() (string, error) {
	// หา path ของ main binary จาก parent process
	mainPath, err := getMainBinaryPath()
	if err != nil {
		return "", err
	}

	file, err := os.Open(mainPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func getMainBinaryPath() (string, error) {
	// วิธีง่าย: ส่ง path ผ่าน environment variable
	path := os.Getenv("MAIN_BINARY_PATH")
	if path == "" {
		return "", fmt.Errorf("MAIN_BINARY_PATH not set")
	}
	return path, nil
}

// Handshake configuration
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "CHECKER_PLUGIN",
	MagicCookieValue: "checker123",
}

// Plugin map
var pluginMap = map[string]plugin.Plugin{
	"checker": &shared.CheckerPlugin{Impl: &HashChecker{}},
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
