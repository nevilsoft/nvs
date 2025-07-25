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

package controllers

import (
	"github.com/burapha44/example/api/v1/services"
	"github.com/burapha44/example/handler"

	"github.com/gofiber/fiber/v2"
)

type BaseController struct {
	Services *services.BaseService
}

func NewBaseController(bs *services.BaseService) *BaseController {
	return &BaseController{
		Services: bs,
	}
}

// Exemple Doc
// @Summary      Get Server Info
// @Description  Get server info and dependencies status and uptime of server and more
// @Tags         Base
// @Produce      json
// @Success      200				 {object}  services.ServerInfoResponse
// @Router       /api/v1/server/info [get]
func (base *BaseController) Health(c *fiber.Ctx) error {
	return handler.Success(c, base.Services.ServerInfo(c))
}
