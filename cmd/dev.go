package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

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
			fmt.Println("✅ swag installed successfully")
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

			fmt.Println("✅ air installed successfully")

			// บน Windows อาจต้อง restart terminal หรือ add Go bin to PATH
			if runtime.GOOS == "windows" {
				fmt.Println("💡 On Windows, you may need to restart your terminal or add Go bin to PATH")
				fmt.Println("   Go bin path is usually: %GOPATH%\\bin or %GOROOT%\\bin")
			}
			return
		}

		// ตรวจสอบว่าไฟล์ .air.toml มีอยู่หรือไม่
		airConfigPath := ".air.toml"
		if _, err := os.Stat(airConfigPath); os.IsNotExist(err) {
			fmt.Println("⚠️ .air.toml not found, creating default configuration...")

			// สร้างไฟล์ .air.toml เริ่มต้น
			defaultConfig := `root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true`

			if err := os.WriteFile(airConfigPath, []byte(defaultConfig), 0644); err != nil {
				fmt.Printf("❌ Failed to create .air.toml: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Created default .air.toml configuration")
		}

		// เรียกใช้งาน air
		os.Setenv("ENV", "dev")
		fmt.Println("🚀 Running dev with Air...")
		runCmd := exec.Command("air", "-c", airConfigPath)
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
