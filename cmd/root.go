/*
 * Created on Tue Mar 04 2025
 *
 * © 2025 Nevilsoft Part., Ltd. All Rights Reserved.
 *
 * * ข้อมูลลับและสงวนสิทธิ์ *
 * ไฟล์นี้เป็นทรัพย์สินของ Nevilsoft Part., Ltd. และมีข้อมูลที่เป็นความลับทางธุรกิจ
 * อนุญาตให้เฉพาะพนักงานที่ได้รับสิทธิ์เข้าถึงเท่านั้น
 * ห้ามเผยแพร่ คัดลอก ดัดแปลง หรือใช้งานโดยไม่ได้รับอนุญาตจากฝ่ายบริหาร
 *
 * การละเมิดข้อตกลงนี้ อาจมีผลให้ถูกลงโทษทางวินัย รวมถึงการดำเนินคดีตามกฎหมาย
 * ตามพระราชบัญญัติว่าด้วยการกระทำความผิดเกี่ยวกับคอมพิวเตอร์ พ.ศ. 2560 (มาตรา 7, 9, 10)
 * และกฎหมายอื่นที่เกี่ยวข้อง
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd คือคำสั่งหลักของ CLI
var RootCmd = &cobra.Command{
	Use:   "nvscli",
	Short: "CLI สำหรับสร้างโครงสร้างโปรเจกต์ Golang",
	Long:  `CLI ที่ช่วยสร้างโครงสร้างโปรเจกต์ Golang พร้อมตั้งค่าต่างๆ เช่น Fiber, PostgreSQL`,
}

// Execute ทำให้ CLI ทำงาน
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
