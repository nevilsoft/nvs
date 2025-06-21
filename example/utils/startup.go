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

package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
)

// StartupStatus represents the status of startup operations
type StartupStatus struct {
	Message   string
	Status    string
	Timestamp time.Time
}

// StartupManager manages startup messages and status
type StartupManager struct {
	statuses []StartupStatus
}

// NewStartupManager creates a new startup manager
func NewStartupManager() *StartupManager {
	return &StartupManager{
		statuses: make([]StartupStatus, 0),
	}
}

// AddStatus adds a new status message
func (sm *StartupManager) AddStatus(message, status string) {
	startupStatus := StartupStatus{
		Message:   message,
		Status:    status,
		Timestamp: time.Now(),
	}
	sm.statuses = append(sm.statuses, startupStatus)

}

// printStatus prints a single status message
func (sm *StartupManager) PrintStatus(status StartupStatus) {
	var statusIcon string
	switch status.Status {
	case "SUCCESS":
		statusIcon = "✅"
	case "ERROR":
		statusIcon = "❌"
	case "WARNING":
		statusIcon = "⚠️"
	case "INFO":
		statusIcon = "ℹ️"
	case "LOADING":
		statusIcon = "🔄"
	default:
		statusIcon = "📝"
	}

	fmt.Printf("%s %s\n", statusIcon, status.Message)
}

const bannerWidth = 80

func ShowBanner(title string, lines ...string) {
	border := strings.Repeat("═", bannerWidth)
	fmt.Printf("╔%s╗\n", border)

	PrintBannerLine(title, "center")

	for _, line := range lines {
		parts := strings.SplitN(line, "|", 2)
		align := "center"
		text := line
		if len(parts) == 2 {
			align = parts[0]
			text = parts[1]
		}
		PrintBannerLine(text, align)
	}

	fmt.Printf("╚%s╝\n", border)
	fmt.Println()
}

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(text string) string {
	return ansiRegexp.ReplaceAllString(text, "")
}

func PrintBannerLine(text string, align string) {
	displayWidth := runewidth.StringWidth(stripANSI(text))
	if displayWidth > bannerWidth {
		text = runewidth.Truncate(text, bannerWidth, "…")
		displayWidth = runewidth.StringWidth(text)
	}

	switch align {
	case "start":
		leftPad := strings.Repeat(" ", 2)
		rightPad := strings.Repeat(" ", bannerWidth-displayWidth-2)
		fmt.Printf("║%s%s%s║\n", leftPad, text, rightPad)

	case "end":
		leftPad := strings.Repeat(" ", bannerWidth-displayWidth)
		fmt.Printf("║%s%s║\n", leftPad, text)

	case "center", "":
		padding := (bannerWidth - displayWidth) / 2
		leftPad := strings.Repeat(" ", padding)
		rightPad := strings.Repeat(" ", bannerWidth-displayWidth-padding)
		fmt.Printf("║%s%s%s║\n", leftPad, text, rightPad)

	default:
		// fallback: center
		padding := (bannerWidth - displayWidth) / 2
		leftPad := strings.Repeat(" ", padding)
		rightPad := strings.Repeat(" ", bannerWidth-displayWidth-padding)
		fmt.Printf("║%s%s%s║\n", leftPad, text, rightPad)
	}
}

// ShowStartupBanner displays the startup banner
func (sm *StartupManager) ShowStartupBanner() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                        NVS CLI v1.0.0                        ║")
	fmt.Println("║                  © 2025 Nevilsoft Part., Ltd.                ║")
	fmt.Println("║                                                              ║")
	fmt.Println("║  🚀 กำลังเริ่มต้นระบบ...                                         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

// ShowStartupSummary displays a summary of all startup operations
func (sm *StartupManager) ShowStartupSummary() {
	fmt.Println()
	fmt.Println("📊 สรุปการเริ่มต้นระบบ:")
	fmt.Println("┌──────────────────────────────────────────────────────────────┐")

	successCount := 0
	errorCount := 0
	warningCount := 0

	for _, status := range sm.statuses {
		switch status.Status {
		case "SUCCESS":
			successCount++
		case "ERROR":
			errorCount++
		case "WARNING":
			warningCount++
		}
	}

	fmt.Printf("│ ✅ สำเร็จ: %d  |  ❌ ผิดพลาด: %d  |  ⚠️ เตือน: %d           	       │\n",
		successCount, errorCount, warningCount)
	fmt.Println("└──────────────────────────────────────────────────────────────┘")

	if errorCount > 0 {
		fmt.Println("⚠️  พบข้อผิดพลาดในการเริ่มต้นระบบ กรุณาตรวจสอบข้อความด้านบน")
	} else {
		fmt.Println("🎉 ระบบเริ่มต้นสำเร็จแล้ว!")
	}
	fmt.Println()
}

// GetStatuses returns all status messages
func (sm *StartupManager) GetStatuses() []StartupStatus {
	return sm.statuses
}

// HasErrors checks if there are any errors in startup
func (sm *StartupManager) HasErrors() bool {
	for _, status := range sm.statuses {
		if status.Status == "ERROR" {
			return true
		}
	}
	return false
}

func StartupMessage() {

	// สร้าง startup manager
	sm := NewStartupManager()

	// แสดง banner
	sm.ShowStartupBanner()

	// จำลองการเริ่มต้นระบบ
	sm.AddStatus("ตรวจสอบการเชื่อมต่อฐานข้อมูล", "LOADING")
	time.Sleep(1 * time.Second)
	sm.AddStatus("ตรวจสอบการเชื่อมต่อฐานข้อมูล", "SUCCESS")

	sm.AddStatus("โหลดไฟล์การตั้งค่า", "LOADING")
	time.Sleep(500 * time.Millisecond)
	sm.AddStatus("โหลดไฟล์การตั้งค่า", "SUCCESS")

	sm.AddStatus("ตรวจสอบสิทธิ์การเข้าถึง", "LOADING")
	time.Sleep(800 * time.Millisecond)
	sm.AddStatus("ตรวจสอบสิทธิ์การเข้าถึง", "WARNING")

	sm.AddStatus("เริ่มต้น Web Server", "LOADING")
	time.Sleep(1 * time.Second)
	sm.AddStatus("เริ่มต้น Web Server", "SUCCESS")

	sm.AddStatus("โหลด Middleware", "LOADING")
	time.Sleep(600 * time.Millisecond)
	sm.AddStatus("โหลด Middleware", "SUCCESS")

	sm.AddStatus("ตรวจสอบการเชื่อมต่อ Redis", "LOADING")
	time.Sleep(700 * time.Millisecond)
	sm.AddStatus("ตรวจสอบการเชื่อมต่อ Redis", "ERROR")

	// แสดงสรุป
	sm.ShowStartupSummary()

	// ถ้ามี error ให้ exit ด้วย code 1
	if sm.HasErrors() {
		os.Exit(1)
	}
}

// InitializeStartup initializes the startup process
func InitializeStartup() *StartupManager {
	sm := NewStartupManager()
	sm.ShowStartupBanner()

	// ตรวจสอบการตั้งค่าพื้นฐาน
	sm.AddStatus("ตรวจสอบ Go version", "INFO")
	sm.AddStatus("ตรวจสอบการตั้งค่า Git", "INFO")
	sm.AddStatus("ตรวจสอบโฟลเดอร์โปรเจกต์", "INFO")

	return sm
}
