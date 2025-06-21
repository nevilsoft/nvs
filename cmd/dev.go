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
			fmt.Println("🚨 ไม่พบคำสั่ง 'swag' กำลังติดตั้ง...")

			installCmd := exec.Command("go", "install", "github.com/swaggo/swag/cmd/swag@latest")
			installCmd.Stdout = os.Stdout
			installCmd.Stderr = os.Stderr

			if err := installCmd.Run(); err != nil {
				fmt.Println("❌ ติดตั้ง swag ไม่สำเร็จ:", err)
				os.Exit(1)
			}
			os.Exit(1)
		}

		// ตรวจสอบว่า air ถูกติดตั้งหรือไม่
		_, err = exec.LookPath("air")
		if err != nil {
			fmt.Println("🚨 ไม่พบคำสั่ง 'air' กำลังติดตั้ง...")

			// ติดตั้ง air ด้วย go install
			installCmd := exec.Command("go", "install", "github.com/cosmtrek/air@latest")
			installCmd.Stdout = os.Stdout
			installCmd.Stderr = os.Stderr

			if err := installCmd.Run(); err != nil {
				fmt.Println("❌ ติดตั้ง air ไม่สำเร็จ:", err)
				os.Exit(1)
			}

			fmt.Println("✅ ติดตั้ง air เรียบร้อยแล้ว กรุณาเปิด Terminal ใหม่ หรือเพิ่ม Go bin ลงใน PATH.")
			return
		}

		// เรียกใช้งาน air
		os.Setenv("ENV", "dev")
		fmt.Println("🚀 กำลังรัน dev ด้วย Air...")
		runCmd := exec.Command("air", "-c", ".air.toml")
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		runCmd.Stdin = os.Stdin

		if err := runCmd.Run(); err != nil {
			fmt.Println("❌ รัน air ไม่สำเร็จ:", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(devCmd)
}
