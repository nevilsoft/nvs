/*
 * Created on Tue Mar 04 2025
 *
 * © 2025 Nevilsoft Ltd., Part. All Rights Reserved.
 *
 * * ข้อมูลลับและสงวนสิทธิ์ *
 * ไฟล์นี้เป็นทรัพย์สินของ Nevilsoft Ltd., Part. และมีข้อมูลที่เป็นความลับทางธุรกิจ
 * อนุญาตให้เฉพาะพนักงานที่ได้รับสิทธิ์เข้าถึงเท่านั้น
 * ห้ามเผยแพร่ คัดลอก ดัดแปลง หรือใช้งานโดยไม่ได้รับอนุญาตจากฝ่ายบริหาร
 *
 * การละเมิดข้อตกลงนี้ อาจมีผลให้ถูกลงโทษทางวินัย รวมถึงการดำเนินคดีตามกฎหมาย
 * ตามพระราชบัญญัติว่าด้วยการกระทำความผิดเกี่ยวกับคอมพิวเตอร์ พ.ศ. 2560 (มาตรา 7, 9, 10)
 * และกฎหมายอื่นที่เกี่ยวข้อง
 */

package cmd

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed templates/*
var templatesFS embed.FS

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new Golang project structure",
	Run: func(cmd *cobra.Command, args []string) {
		wd, _ := os.Getwd()
		repo, _ := cmd.Flags().GetString("repo")

		defaultProjectName := filepath.Base(wd)
		projectName := prompt("📦 Project name", defaultProjectName)

		// ตรวจสอบ Git user.name
		gitUser, err := getGitUserName()
		if err != nil {
			fmt.Println("⚠️ Git user.name not found:", err)
			fmt.Println("👉 Please set Git user.name before using")
			return
		}

		// กำหนดโฟลเดอร์เป้าหมาย
		targetDir, err := determineTargetDirectory(projectName, wd)
		if err != nil {
			fmt.Println("❌ Unable to create directory:", err)
			return
		}

		// กำหนด module name
		if repo == "" {
			repo = prompt("🧱 Go module name", fmt.Sprintf("github.com/%s/%s", gitUser, projectName))
		}

		fmt.Printf("🔧 Creating project %s (module: %s)...\n", projectName, repo)

		// คัดลอกไฟล์เทมเพลต
		err = copyEmbeddedTemplates("templates", targetDir, map[string]string{
			"ProjectName": projectName,
			"ModuleName":  repo,
		})
		if err != nil {
			fmt.Println("❌ Error:", err)
			return
		}

		fmt.Println("✅ Project created successfully!")
		fmt.Println("📂 After entering the project directory, run the following commands:")
		fmt.Println("  1. go mod tidy")
		fmt.Println("  2. nvs dev")
		fmt.Println("  3. nvs build")
		fmt.Println("  4. nvs start")
	},
}

// getGitUserName ดึงชื่อผู้ใช้จาก Git config
func getGitUserName() (string, error) {
	cmd := exec.Command("git", "config", "--global", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	userName := strings.TrimSpace(string(output))
	if userName == "" {
		return "", fmt.Errorf("git user.name is empty")
	}
	return userName, nil
}

// determineTargetDirectory กำหนดโฟลเดอร์เป้าหมายสำหรับสร้างโปรเจกต์
func determineTargetDirectory(projectName, currentDir string) (string, error) {
	defaultProjectName := filepath.Base(currentDir)

	if projectName == defaultProjectName {
		return ".", nil
	}

	confirmName := confirm(fmt.Sprintf("❓ Use \"%s\" to create this project?", projectName))
	if !confirmName {
		return "", fmt.Errorf("user cancelled project creation")
	}

	if err := os.Mkdir(projectName, 0755); err != nil {
		return "", err
	}

	return projectName, nil
}

// prompt อ่านค่า input แบบมี default
func prompt(question, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (%s): ", question, defaultValue)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return defaultValue
	}
	return text
}

// confirm ยืนยันคำตอบ y/n default y
func confirm(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n) [y]: ", question)
	text, _ := reader.ReadString('\n')
	text = strings.ToLower(strings.TrimSpace(text))
	if text == "" {
		return true
	}
	return text == "y" || text == "yes"
}

// copyEmbeddedTemplates คัดลอกไฟล์จาก embed และ render ด้วย template
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

		// ถ้าเป็นไฟล์ .tmpl ให้ parse เป็น template
		if strings.HasSuffix(d.Name(), ".tmpl") {
			content, err := templatesFS.ReadFile(path)
			if err != nil {
				return err
			}

			tmplName := strings.TrimPrefix(path, srcDir+"/")
			tmpl, err := template.New(tmplName).Parse(string(content))
			if err != nil {
				return fmt.Errorf("template parse error in %s: %w", path, err)
			}

			// ตัด .tmpl ออกจากชื่อไฟล์
			targetPath = strings.TrimSuffix(targetPath, ".tmpl")

			file, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer file.Close()

			return tmpl.Execute(file, data)
		}

		// ไม่ใช่ .tmpl → คัดลอกตรง ๆ
		rawContent, err := templatesFS.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(targetPath, rawContent, 0755)
	})
}

func init() {
	initCmd.Flags().String("repo", "", "Specify the module name (e.g. github.com/user/project)")
	RootCmd.AddCommand(initCmd)
}
