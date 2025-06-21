package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	output string
	target string
	semver string
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "สร้างโปรเจกต์แบบ obfuscate (ด้วย garble)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔒 กำลัง build แบบ obfuscate ด้วย garble...")

		// ตรวจสอบว่า garble ติดตั้งอยู่หรือไม่
		if _, err := exec.LookPath("garble"); err != nil {
			fmt.Println("📦 กำลังติดตั้ง garble...")
			install := exec.Command("go", "install", "mvdan.cc/garble@latest")
			install.Stdout = os.Stdout
			install.Stderr = os.Stderr
			install.Env = os.Environ()
			if err := install.Run(); err != nil {
				fmt.Println("❌ ติดตั้ง garble ไม่สำเร็จ:", err)
				return
			}
		}

		// แยก GOOS และ GOARCH จาก target
		goos := runtime.GOOS
		goarch := runtime.GOARCH
		if target != "" {
			parts := strings.Split(target, "/")
			if len(parts) != 2 {
				fmt.Println("❌ รูปแบบ target ไม่ถูกต้อง. ใช้รูปแบบ เช่น linux/amd64 หรือ windows/amd64")
				return
			}
			goos, goarch = parts[0], parts[1]
		}

		// ตั้งชื่อไฟล์ output
		if output == "" {
			output = "main"
		}
		if goos == "windows" && !strings.HasSuffix(output, ".exe") {
			output += ".exe"
		}
		buildNumber := fmt.Sprintf("%06d", rand.Intn(1000000))
		versionString := semver

		if buildNumber != "" {
			versionString += "." + buildNumber
		}

		ldflags := fmt.Sprintf("-X 'main.Version=%s'", versionString)
		ldflags += fmt.Sprintf("-X 'main.RunnerID=%s'", "132456")

		// สั่ง garble build
		cmdGarble := exec.Command("garble", "build", "-ldflags", ldflags, "-o", output, "./main.go")
		cmdGarble.Stdout = os.Stdout
		cmdGarble.Stderr = os.Stderr
		cmdGarble.Env = append(os.Environ(),
			"GOOS="+goos,
			"GOARCH="+goarch,
		)

		fmt.Printf("🛠️  Target: %s/%s\n", goos, goarch)
		if err := cmdGarble.Run(); err != nil {
			fmt.Println("❌ ไม่สามารถ build ได้:", err)
			return
		}

		fmt.Printf("✅ สร้างโปรแกรม obfuscated เรียบร้อยที่: %s\n", output)
	},
}

func init() {
	buildCmd.Flags().StringVarP(&output, "output", "o", "", "ชื่อไฟล์ output (default: main)")
	buildCmd.Flags().StringVarP(&target, "target", "t", "", "ระบบเป้าหมาย (เช่น linux/amd64, windows/amd64, darwin/arm64)")
	buildCmd.Flags().StringVarP(&semver, "version", "v", "dev", "Build version (default: dev)")
	RootCmd.AddCommand(buildCmd)
}
