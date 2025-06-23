//go:build wireinject
// +build wireinject

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

package di

import (
	"github.com/burapha44/example/api/v1/controllers"
	"github.com/burapha44/example/api/v1/middleware"
	"github.com/burapha44/example/api/v1/services"

	"github.com/google/wire"
)

func NewAppContainer() (*AppContainer, error) {
	wire.Build(
		middleware.ProviderSet,
		services.ProviderSet,
		controllers.ProviderSet,
		wire.Struct(new(AppContainer), "*"),
	)
	return &AppContainer{}, nil
}

type AppContainer struct {
	AuthMiddleware    *middleware.BaseMiddleware
	BaseController    *controllers.BaseController
	ProductController *controllers.ProductController
}
