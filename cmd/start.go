package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		// ตรวจสอบว่าไฟล์มีอยู่จริง
		if _, err := os.Stat(name); os.IsNotExist(err) {
			fmt.Printf("❌ File not found: %s\n", name)
			return
		}

		// ตรวจสอบว่าไฟล์เป็น executable หรือไม่
		if runtime.GOOS != "windows" {
			// บน Unix-like systems ตรวจสอบ execute permission
			info, err := os.Stat(name)
			if err != nil {
				fmt.Printf("❌ Unable to get file info: %v\n", err)
				return
			}

			// ตรวจสอบ execute permission
			if info.Mode()&0111 == 0 {
				fmt.Printf("⚠️ File %s is not executable. Attempting to make it executable...\n", name)
				if err := os.Chmod(name, info.Mode()|0111); err != nil {
					fmt.Printf("❌ Failed to make file executable: %v\n", err)
					return
				}
				fmt.Println("✅ Made file executable")
			}
		}

		// ใช้ absolute path สำหรับ executable
		absPath, err := filepath.Abs(name)
		if err != nil {
			fmt.Printf("❌ Unable to get absolute path: %v\n", err)
			return
		}

		hash, err := calculateSHA256(absPath)
		if err != nil {
			fmt.Println("❌ Unable to compute hash:", err)
			return
		}

		env := cmd.Flag("env").Value.String()
		if env == "dev" {
			os.Setenv("ENV", "dev")
		} else {
			os.Setenv("ENV", "prod")
		}

		os.Setenv("RUNNER_ID", hash)

		// ใช้ absolute path ในการรัน command
		cmdRun := exec.Command(absPath)
		cmdRun.Stdout = os.Stdout
		cmdRun.Stderr = os.Stderr
		cmdRun.Stdin = os.Stdin

		fmt.Printf("🚀 Starting %s with ENV=%s, RUNNER_ID=%s\n", filepath.Base(absPath), os.Getenv("ENV"), hash)

		if err := cmdRun.Run(); err != nil {
			fmt.Printf("❌ Unable to run %s: %v\n", filepath.Base(absPath), err)
		}
	},
}

func calculateSHA256(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func init() {
	startCmd.Flags().StringP("env", "e", "prod", "Environment")
	RootCmd.AddCommand(startCmd)
}
