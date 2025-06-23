package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

type Package struct {
	InstallCmd string
	PkgName    string
}

func fetchPackagesFromURL(url string) ([]Package, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error calling data: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to convert HTML data: %w", err)
	}

	packages := []Package{}
	doc.Find("div.SearchSnippet-headerContainer > h2 > a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		pkgName := strings.TrimSpace(s.Contents().Not("span").Text())
		installCmd := strings.TrimPrefix(href, "/")

		packages = append(packages, Package{
			InstallCmd: installCmd,
			PkgName:    pkgName,
		})
	})

	return packages, nil
}

func getPackage(pkgName string) (string, string, bool) {
	if strings.HasPrefix(pkgName, "https://pkg.go.dev") || strings.HasPrefix(pkgName, "https://github.com") {
		return pkgName, "", true
	}

	url := "https://pkg.go.dev/search?q=" + pkgName
	packages, err := fetchPackagesFromURL(url)
	if err != nil {
		log.Printf("failed to search package: %v", err)
		return "", "", false
	}

	if len(packages) == 0 {
		return "", "", false
	}

	return packages[0].InstallCmd, packages[0].PkgName, true
}

func getPackageList(pkgName string) []Package {
	url := "https://pkg.go.dev/search?q=" + pkgName
	packages, err := fetchPackagesFromURL(url)
	if err != nil {
		log.Printf("failed to search package: %v", err)
		return []Package{}
	}

	if len(packages) > 5 {
		return packages[:5]
	}
	return packages
}

func installPackage(installCmd string) error {
	cmd := exec.Command("go", "get", "-u", installCmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install package: %w\nOutput: %s", err, string(output))
	}
	return nil
}

var getAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"add", "a"},
	Short:   "Install package",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("❌ Please specify the package name to install")
			return
		}

		fmt.Println(" ")
		fmt.Println("⚡️ Checking package...")

		var wg sync.WaitGroup
		for i, arg := range args {
			wg.Add(1)
			go func(i int, arg string) {
				defer wg.Done()
				fmt.Printf("⏳ %d/%d Installing %s\n", 0, len(args), arg)

				installCmd, _, ok := getPackage(arg)
				if !ok {
					fmt.Printf("❌ Package %s not found\n", arg)
					return
				}

				err := installPackage(installCmd)
				if err != nil {
					fmt.Printf("❌ Failed to install %s: %v\n", installCmd, err)
					return
				}

				fmt.Printf("✅ %d/%d Installed %s\n\n", i+1, len(args), installCmd)
			}(i, arg)
		}
		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(getAddCmd)
}
