package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Code generator commands",
}

var generateRouteCmd = &cobra.Command{
	Use:   "route [name]",
	Short: "Generate a new route file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		routeName := args[0]
		if routeName == "" {
			fmt.Println("❌ Route name is required")
			return
		}

		// Capitalize first letter for function name
		routeNameTitle := strings.Title(routeName)

		tmplPath := filepath.Join("templates", "api", "v1", "routes", "route.go.tmpl")
		outputDir := filepath.Join("api", "v1", "routes")
		outputFile := filepath.Join(outputDir, routeName+".go")

		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Println("❌ Cannot create routes directory:", err)
			return
		}

		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			fmt.Println("❌ Cannot read template:", err)
			return
		}

		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("❌ Cannot create route file:", err)
			return
		}
		defer file.Close()

		data := map[string]string{
			"RouteName": routeNameTitle,
		}

		if err := tmpl.Execute(file, data); err != nil {
			fmt.Println("❌ Cannot render template:", err)
			return
		}

		fmt.Printf("✅ Created route: %s\n", outputFile)
	},
}

var generateRoutesCmd = &cobra.Command{
	Use:   "routes",
	Short: "Auto-generate route files from controller @Router/@Tags comments",
	Run: func(cmd *cobra.Command, args []string) {
		controllersDir := filepath.Join("api", "v1", "controllers")
		routesDir := filepath.Join("api", "v1", "routes")
		os.MkdirAll(routesDir, 0755)

		controllerFiles, _ := filepath.Glob(filepath.Join(controllersDir, "*.go"))
		if len(controllerFiles) == 0 {
			fmt.Println("❌ No controller files found.")
			return
		}

		// Regex for @Router, @Tags, and method signature
		reRouter := regexp.MustCompile(`@Router\s+([^\s]+) \[([a-zA-Z]+)\]`)
		reTag := regexp.MustCompile(`@Tags\s+([A-Za-z0-9_]+)`)
		reFunc := regexp.MustCompile(`func \(.*\*([A-Za-z0-9_]+)\) ([A-Za-z0-9_]+)\(.*\*fiber.Ctx.*\) error`)

		tagRoutes := map[string][]map[string]string{} // tag => []{method, path, handler, controller}

		for _, file := range controllerFiles {
			f, err := os.Open(file)
			if err != nil {
				fmt.Println("❌ Cannot open:", file)
				continue
			}
			scanner := bufio.NewScanner(f)
			var lastTag, lastRouter, lastMethod, lastHandler, lastController string
			for scanner.Scan() {
				line := scanner.Text()
				if m := reTag.FindStringSubmatch(line); m != nil {
					lastTag = m[1]
				}
				if m := reRouter.FindStringSubmatch(line); m != nil {
					lastRouter = m[1]
					lastMethod = strings.ToUpper(m[2])
				}
				if m := reFunc.FindStringSubmatch(line); m != nil {
					lastController = m[1]
					lastHandler = m[2]
					if lastTag != "" && lastRouter != "" && lastMethod != "" {
						tag := lastTag
						route := map[string]string{
							"method":     lastMethod,
							"path":       lastRouter,
							"handler":    lastHandler,
							"controller": lastController,
						}
						tagRoutes[tag] = append(tagRoutes[tag], route)
						// Reset for next
						lastRouter, lastMethod, lastHandler, lastController = "", "", "", ""
						lastTag = ""
					}
				}
			}
			f.Close()
		}

		// Read module name from go.mod
		moduleName := readModuleName()

		// For each tag, generate a route file (skip Base)
		for tag, routes := range tagRoutes {
			if strings.ToLower(tag) == "base" {
				continue // do not generate/overwrite base.go
			}
			// Remove 'controller' or 'Controller' suffix from tag for filename
			baseTag := tag
			if strings.HasSuffix(strings.ToLower(baseTag), "controller") {
				baseTag = baseTag[:len(baseTag)-10]
			}
			snakeTag := toSnakeCase(baseTag)
			fileName := filepath.Join(routesDir, snakeTag+"_route.go")
			file, err := os.Create(fileName)
			if err != nil {
				fmt.Println("❌ Cannot create:", fileName)
				continue
			}
			tmpl := template.New("route").Funcs(template.FuncMap{
				"title": func(s string) string {
					if len(s) == 0 {
						return s
					}
					return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
				},
			})
			tmpl = template.Must(tmpl.Parse(routeFileTemplate))
			data := map[string]interface{}{
				"Tag":        strings.Split(tag, "Controller")[0],
				"Routes":     routes,
				"ModuleName": moduleName,
			}
			tmpl.Execute(file, data)
			file.Close()
			fmt.Printf("✅ Generated: %s\n", fileName)

			// Remove controller name from route path in each route
			for i, r := range routes {
				if path, ok := r["path"]; ok {
					// Remove /Tag or /TagController from the path
					cleanPath := regexp.MustCompile(`/`+tag+`(Controller)?`).ReplaceAllString(path, "")
					if cleanPath == "" {
						cleanPath = "/"
					}
					routes[i]["path"] = cleanPath
				}
			}
		}

		// After generating route files, update SetupRoutes in base.go
		updateSetupRoutes(routesDir)
	},
}

var generateControllerCmd = &cobra.Command{
	Use:   "controller [name]",
	Short: "Generate a new controller and register it in ProviderSet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if name == "" {
			fmt.Println("❌ Controller name is required")
			return
		}
		controllerName := strings.Title(name) + "Controller"
		fileName := filepath.Join("api", "v1", "controllers", strings.ToLower(name)+".go")
		if _, err := os.Stat(fileName); err == nil {
			fmt.Println("❌ Controller file already exists:", fileName)
			return
		}
		// Generate controller file
		controllerTmpl := `package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type ` + controllerName + ` struct {}

func New` + controllerName + `() *` + controllerName + ` {
	return &` + controllerName + `{}
}

// Exemple Function
// @Summary      Get Server Info
// @Description  Get server info and dependencies status and uptime of server and more
// @Tags         ` + controllerName + `
// @Produce      json
// @Success      200				 {object}  services.ServerInfoResponse
// @Router       /api/v1/` + controllerName + `/info [get]
func (c *` + controllerName + `) Example(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello from ` + controllerName + `")
}
`
		os.WriteFile(fileName, []byte(controllerTmpl), 0644)
		fmt.Printf("✅ Created controller: %s\n", fileName)

		// Register in providers.go
		providersFile := filepath.Join("api", "v1", "controllers", "providers.go")
		providersContent, err := os.ReadFile(providersFile)
		if err != nil {
			fmt.Println("❌ Cannot read providers.go:", err)
			return
		}
		providersStr := string(providersContent)
		// Add NewXxxController to ProviderSet if not exists
		added := false
		if !strings.Contains(providersStr, "New"+controllerName) {
			reSet := regexp.MustCompile(`(?s)(var ProviderSet = wire.NewSet\()(.*?)(\))`)
			providersStr = reSet.ReplaceAllString(providersStr, "${1}\n\tNew"+controllerName+",${2}${3}")
			os.WriteFile(providersFile, []byte(providersStr), 0644)
			fmt.Println("✅ Registered New" + controllerName + " in ProviderSet")
			added = true
		} else {
			fmt.Println("ℹ️  New" + controllerName + " already registered in ProviderSet")
		}

		// Add XxxController to AppContainer struct in di/wire.go
		addedField := addControllerToAppContainer(controllerName)
		if added || addedField {
			// Run wire ./di only if something was added
			wireCmd := exec.Command("wire", "./di")
			wireCmd.Stdout = os.Stdout
			wireCmd.Stderr = os.Stderr
			if err := wireCmd.Run(); err != nil {
				fmt.Println("❌ Cannot run wire:", err)
			} else {
				fmt.Println("✅ wire ./di completed")
			}
		}
	},
}

const routeFileTemplate = `/*
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
 
 package routes

import (
	"github.com/gofiber/fiber/v2"
	"{{ .ModuleName }}/di"
)

func Register{{ .Tag }}Routes(app fiber.Router, c *di.AppContainer) {
    {{ $tag := .Tag }}
    {{ $ctrl := printf "%s" .Tag }}
    {{ $inst := printf "%s := c.%sController" $tag $tag }}
    {{ $inst }}
{{- range .Routes }}
	app.{{ .method | title }}("{{ .path }}", {{ $tag }}.{{ .handler }})
{{- end }}
}
`

func init() {
	generateCmd.AddCommand(generateRouteCmd)
	generateCmd.AddCommand(generateRoutesCmd)
	generateCmd.AddCommand(generateControllerCmd)
	RootCmd.AddCommand(generateCmd)
}

// Helper for title-case method
func title(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

// Helper: read module name from go.mod
func readModuleName() string {
	f, err := os.Open("go.mod")
	if err != nil {
		return ""
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}

// Update SetupRoutes in base.go to call all RegisterXxxRoutes
func updateSetupRoutes(routesDir string) {
	baseFile := filepath.Join(routesDir, "base.go")
	files, _ := filepath.Glob(filepath.Join(routesDir, "*.go"))
	var registerCalls []string
	for _, f := range files {
		if strings.HasSuffix(f, "base.go") {
			continue
		}
		content, _ := os.ReadFile(f)
		re := regexp.MustCompile(`func (Register([A-Za-z0-9_]+)Routes)\s*\(`)
		matches := re.FindAllStringSubmatch(string(content), -1)
		for _, m := range matches {
			if m[1] == "RegisterRoutes" {
				continue
			}
			registerCalls = append(registerCalls, fmt.Sprintf("\t%s(v1API, container)", m[1]))
		}
	}

	baseContent, _ := os.ReadFile(baseFile)
	// Read base.go and collect existing RegisterXxxRoutes calls
	baseContentStr := string(baseContent)
	existing := map[string]bool{}
	reCall := regexp.MustCompile(`Register([A-Za-z0-9_]+)Routes\(v1API, container\)`) // match function calls
	for _, m := range reCall.FindAllStringSubmatch(baseContentStr, -1) {
		existing[m[0]] = true
	}

	// Only add calls that do not already exist
	var uniqueCalls []string
	for _, call := range registerCalls {
		callTrim := strings.TrimSpace(call)
		if !existing[callTrim] {
			uniqueCalls = append(uniqueCalls, callTrim)
		}
	}

	// Replace only the auto-generated line
	reAuto := regexp.MustCompile(`(?m)^\s*// \(auto-generated: add more RegisterXxxRoutes here\).*$`)
	newBlock := strings.Join(uniqueCalls, "\n")
	if newBlock != "" {
		newBlock += "\n"
	}
	newBlock += "\t// (auto-generated: add more RegisterXxxRoutes here)"
	result := reAuto.ReplaceAll(baseContent, []byte(newBlock))
	_ = os.WriteFile(baseFile, result, 0644)
	fmt.Println("✅ Updated RegisterXxxRoutes in base.go")
}

// Add XxxController struct field to AppContainer in di/wire.go
// Returns true if a new field was added
func addControllerToAppContainer(controllerName string) bool {
	wireFile := filepath.Join("di", "wire.go")
	content, err := os.ReadFile(wireFile)
	if err != nil {
		fmt.Println("❌ Cannot read di/wire.go:", err)
		return false
	}
	structField := "\t" + controllerName + " *controllers." + controllerName + "\n"
	if strings.Contains(string(content), structField) {
		fmt.Println("ℹ️  " + controllerName + " already exists in AppContainer")
		return false
	}
	re := regexp.MustCompile(`(?s)(type AppContainer struct \{)(.*?)(\n\})`)
	newContent := re.ReplaceAllString(string(content), "${1}${2}\n"+structField+"${3}")
	os.WriteFile(wireFile, []byte(newContent), 0644)
	fmt.Println("✅ Added " + controllerName + " to AppContainer in di/wire.go")
	return true
}

// Helper: convert CamelCase or PascalCase to snake_case
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
