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
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed templates/*
var templatesFS embed.FS

var newCmd = &cobra.Command{
	Use:   "new [project_name]",
	Short: "สร้างโครงสร้างโปรเจกต์ Golang",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		dbType, _ := cmd.Flags().GetString("db")
		useGRPC, _ := cmd.Flags().GetBool("grpc")
		repo, _ := cmd.Flags().GetString("repo")

		// ถ้าไม่มี `--repo` ให้ใช้ projectName เป็น module
		if repo == "" {
			repo = projectName
		}

		fmt.Printf("🔧 กำลังสร้างโปรเจกต์ %s (module: %s)...\n", projectName, repo)

		// สร้างโฟลเดอร์เป้าหมาย
		if err := os.Mkdir(projectName, 0755); err != nil {
			fmt.Println("❌ ไม่สามารถสร้างโฟลเดอร์:", err)
			return
		}

		// คัดลอกไฟล์เทมเพลต
		err := copyEmbeddedTemplates("templates", projectName, map[string]string{
			"ProjectName": projectName,
			"ModuleName":  repo, // ใช้เป็น module
			"DBType":      dbType,
			"UseGRPC":     fmt.Sprintf("%t", useGRPC),
		})

		if err != nil {
			fmt.Println("❌ พบข้อผิดพลาด:", err)
			return
		}

		fmt.Println("✅ โปรเจกต์สร้างเสร็จแล้ว!")
		fmt.Printf("📂 หลังจากคุณอยู่ในโฟลเดอร์โปรเจกต์แล้ว, ให้รันคำสั่งต่อไปนี้:\n")
		fmt.Printf("  1. cd %s\n", projectName)
		fmt.Println("  2. go mod tidy")
		fmt.Println("  3. go run main.go")
	},
}

// ฟังก์ชันคัดลอกไฟล์จาก embed
func copyEmbeddedTemplates(srcDir, destDir string, data map[string]string) error {
	return fs.WalkDir(templatesFS, srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(srcDir, path)
		targetPath := filepath.Join(destDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		content, err := templatesFS.ReadFile(path)
		if err != nil {
			return err
		}

		tmpl, err := template.New("").Parse(string(content))
		if err != nil {
			return err
		}

		file, err := os.Create(strings.TrimSuffix(targetPath, ".tmpl"))
		if err != nil {
			return err
		}
		defer file.Close()

		return tmpl.Execute(file, data)
	})
}

func init() {
	newCmd.Flags().String("db", "", "เลือกฐานข้อมูล (เช่น postgres, mysql)")
	newCmd.Flags().Bool("grpc", false, "เปิดใช้งาน gRPC")
	newCmd.Flags().String("repo", "", "ระบุ module name (เช่น github.com/user/project)")

	RootCmd.AddCommand(newCmd)
}
