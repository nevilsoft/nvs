package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		input := name
		hash, err := calculateSHA256(input)
		if err != nil {
			fmt.Println("❌ ไม่สามารถคำนวณ hash ได้:", err)
			return
		}
		env := cmd.Flag("env").Value.String()
		if env == "dev" {
			os.Setenv("ENV", "dev")
		} else {
			os.Setenv("ENV", "prod")
		}

		os.Setenv("RUNNER_ID", hash)
		cmdRun := exec.Command(input)
		cmdRun.Stdout = os.Stdout
		cmdRun.Stderr = os.Stderr
		if err := cmdRun.Run(); err != nil {
			fmt.Println("❌ ไม่สามารถรัน main ได้:", err)
		}
	},
}

func calculateSHA256(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func init() {
	startCmd.Flags().StringP("env", "e", "prod", "Environment")
	RootCmd.AddCommand(startCmd)
}
