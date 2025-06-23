package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run project in development mode using Air",
	Run: func(cmd *cobra.Command, args []string) {

		// ตรวจสอบว่า swag ถูกติดตั้งหรือไม่
		_, err := exec.LookPath("swag")
		if err != nil {
			fmt.Println("🚨 'swag' not found, installing...")

			installCmd := exec.Command("go", "install", "github.com/swaggo/swag/cmd/swag@latest")
			installCmd.Stdout = os.Stdout
			installCmd.Stderr = os.Stderr

			if err := installCmd.Run(); err != nil {
				fmt.Println("❌ Failed to install swag:", err)
				os.Exit(1)
			}
			os.Exit(1)
		}

		// ตรวจสอบว่า air ถูกติดตั้งหรือไม่
		_, err = exec.LookPath("air")
		if err != nil {
			fmt.Println("🚨 'air' not found, installing...")

			// ติดตั้ง air ด้วย go install
			installCmd := exec.Command("go", "install", "github.com/cosmtrek/air@latest")
			installCmd.Stdout = os.Stdout
			installCmd.Stderr = os.Stderr

			if err := installCmd.Run(); err != nil {
				fmt.Println("❌ Failed to install air:", err)
				os.Exit(1)
			}

			fmt.Println("✅ air installed successfully, please open a new terminal or add Go bin to PATH.")
			return
		}

		// เรียกใช้งาน air
		os.Setenv("ENV", "dev")
		fmt.Println("🚀 Running dev with Air...")
		runCmd := exec.Command("air", "-c", ".air.toml")
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		runCmd.Stdin = os.Stdin

		if err := runCmd.Run(); err != nil {
			fmt.Println("❌ Failed to run air:", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(devCmd)
}
