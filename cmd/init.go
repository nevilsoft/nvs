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
	"runtime"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed templates/*
var templatesFS embed.FS

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Create a new Golang project structure",
	Long: `Create a new Golang project structure with interactive prompts.

If project-name is provided, it will be used as the project name.
If not provided, you will be prompted to enter a project name.

Examples:
  nvs init                    # Interactive mode
  nvs init my-project         # Create project with name "my-project"
  nvs init my-project --force # Force overwrite existing directory`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("❌ Error getting working directory:", err)
			return
		}

		repo, _ := cmd.Flags().GetString("repo")
		force, _ := cmd.Flags().GetBool("force")

		// กำหนดชื่อโปรเจกต์
		var projectName string
		if len(args) > 0 {
			// ถ้ามี argument ให้ใช้ชื่อนั้น
			projectName = args[0]
			fmt.Printf("📦 Using project name: %s\n", projectName)
		} else {
			// ถ้าไม่มี argument ให้ถามชื่อ
			defaultProjectName := filepath.Base(wd)
			projectName = prompt("📦 Project name", defaultProjectName)
		}

		// ตรวจสอบ Git user.name
		gitUser, err := getGitUserName()
		if err != nil {
			fmt.Println("⚠️ Git user.name not found:", err)
			fmt.Println("👉 Please set Git user.name before using")
			return
		}

		// กำหนดโฟลเดอร์เป้าหมาย
		skipConfirmation := len(args) > 0 // ถ้ามี argument ให้ข้ามการยืนยัน
		targetDir, err := determineTargetDirectory(projectName, wd, skipConfirmation, force)
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
	// ตรวจสอบว่า Git มีอยู่จริง
	if _, err := exec.LookPath("git"); err != nil {
		return "", fmt.Errorf("git command not found in PATH")
	}

	cmd := exec.Command("git", "config", "--global", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git user.name: %w", err)
	}
	userName := strings.TrimSpace(string(output))
	if userName == "" {
		return "", fmt.Errorf("git user.name is empty")
	}
	return userName, nil
}

// determineTargetDirectory กำหนดโฟลเดอร์เป้าหมายสำหรับสร้างโปรเจกต์
func determineTargetDirectory(projectName, currentDir string, skipConfirmation, force bool) (string, error) {
	defaultProjectName := filepath.Base(currentDir)

	if projectName == defaultProjectName {
		return ".", nil
	}

	// ตรวจสอบว่าโฟลเดอร์มีอยู่แล้วหรือไม่
	if _, err := os.Stat(projectName); err == nil {
		// โฟลเดอร์มีอยู่แล้ว
		var overwrite bool
		if force {
			// ถ้าใช้ --force ให้ overwrite โดยไม่ถาม
			overwrite = true
			fmt.Printf("🗑️ Force removing existing directory: %s\n", projectName)
		} else {
			// ถ้าไม่ใช้ --force ให้ถามก่อน
			overwrite = confirmWithDefault(fmt.Sprintf("⚠️ Directory \"%s\" already exists. Overwrite?", projectName), false)
		}

		if !overwrite {
			return "", fmt.Errorf("user cancelled project creation")
		}

		// ลบโฟลเดอร์เก่า
		if err := os.RemoveAll(projectName); err != nil {
			return "", fmt.Errorf("failed to remove existing directory %s: %w", projectName, err)
		}
		if !force {
			fmt.Printf("🗑️ Removed existing directory: %s\n", projectName)
		}
	}

	// ถ้าไม่ต้องข้ามการยืนยัน ให้ถามยืนยัน
	if !skipConfirmation {
		confirmName := confirm(fmt.Sprintf("❓ Use \"%s\" to create this project?", projectName))
		if !confirmName {
			return "", fmt.Errorf("user cancelled project creation")
		}
	}

	// ใช้ permissions ที่เหมาะสมกับแต่ละ OS
	perm := os.FileMode(0755)
	if runtime.GOOS == "windows" {
		perm = os.FileMode(0666) // Windows ไม่สนใจ execute bit
	}

	if err := os.Mkdir(projectName, perm); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", projectName, err)
	}

	// ตรวจสอบว่าโฟลเดอร์สร้างสำเร็จและว่างเปล่า
	if entries, err := os.ReadDir(projectName); err != nil {
		return "", fmt.Errorf("failed to verify directory %s: %w", projectName, err)
	} else if len(entries) > 0 {
		return "", fmt.Errorf("directory %s is not empty after creation", projectName)
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

// confirmWithDefault ยืนยันคำตอบ y/n พร้อม default value
func confirmWithDefault(question string, defaultValue bool) bool {
	reader := bufio.NewReader(os.Stdin)
	defaultStr := "y"
	if !defaultValue {
		defaultStr = "n"
	}
	fmt.Printf("%s (y/n) [%s]: ", question, defaultStr)
	text, _ := reader.ReadString('\n')
	text = strings.ToLower(strings.TrimSpace(text))
	if text == "" {
		return defaultValue
	}
	return text == "y" || text == "yes"
}

// normalizePath แปลง path ให้เป็นรูปแบบที่ถูกต้องสำหรับ OS ปัจจุบัน
func normalizePath(path string) string {
	// แปลง forward slash เป็น OS-specific separator
	return filepath.FromSlash(path)
}

// copyEmbeddedTemplates คัดลอกไฟล์จาก embed และ render ด้วย template
func copyEmbeddedTemplates(srcDir, destDir string, data map[string]string) error {
	return fs.WalkDir(templatesFS, srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// ใช้ filepath.Rel เพื่อให้ได้ relative path ที่ถูกต้อง
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// แปลง path ให้เป็นรูปแบบที่ถูกต้องสำหรับ OS ปัจจุบัน
		relPath = normalizePath(relPath)
		targetPath := filepath.Join(destDir, relPath)

		if d.IsDir() {
			// ใช้ permissions ที่เหมาะสมกับแต่ละ OS
			perm := os.FileMode(0755)
			if runtime.GOOS == "windows" {
				perm = os.FileMode(0666)
			}
			return os.MkdirAll(targetPath, perm)
		}

		// ถ้าเป็นไฟล์ .tmpl ให้ parse เป็น template
		if strings.HasSuffix(d.Name(), ".tmpl") {
			content, err := templatesFS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read template file %s: %w", path, err)
			}

			// ใช้ filepath.Join แทน string concatenation
			tmplName := filepath.Join(strings.TrimPrefix(path, filepath.Join(srcDir, "")))
			tmpl, err := template.New(tmplName).Parse(string(content))
			if err != nil {
				return fmt.Errorf("template parse error in %s: %w", path, err)
			}

			// ตัด .tmpl ออกจากชื่อไฟล์
			targetPath = strings.TrimSuffix(targetPath, ".tmpl")

			// สร้าง directory parent ถ้ายังไม่มี
			targetDir := filepath.Dir(targetPath)
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
			}

			file, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}
			defer file.Close()

			return tmpl.Execute(file, data)
		}

		// ไม่ใช่ .tmpl → คัดลอกตรง ๆ
		rawContent, err := templatesFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// สร้าง directory parent ถ้ายังไม่มี
		targetDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
		}

		// ใช้ permissions ที่เหมาะสมกับแต่ละ OS
		perm := os.FileMode(0644)
		if runtime.GOOS == "windows" {
			perm = os.FileMode(0666)
		}

		return os.WriteFile(targetPath, rawContent, perm)
	})
}

func init() {
	initCmd.Flags().String("repo", "", "Specify the module name (e.g. github.com/user/project)")
	initCmd.Flags().BoolP("force", "f", false, "Force overwrite existing directory")
	RootCmd.AddCommand(initCmd)
}
