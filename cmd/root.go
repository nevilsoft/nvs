package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd คือคำสั่งหลักของ CLI
var RootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "CLI สำหรับสร้างโครงสร้างโปรเจกต์ Golang",
	Long:  `CLI ที่ช่วยสร้างโครงสร้างโปรเจกต์ Golang พร้อมตั้งค่าต่างๆ เช่น Fiber, gRPC, PostgreSQL`,
}

// Execute ทำให้ CLI ทำงาน
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
