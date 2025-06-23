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

package middleware

import (
	"time"

	"github.com/burapha44/example/constants"
	"github.com/burapha44/example/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type BaseMiddleware struct {
}

func NewBaseMiddleware() *BaseMiddleware {
	return &BaseMiddleware{}
}

func (mw *BaseMiddleware) RateLimit(count int, duration time.Duration) fiber.Handler {
	if duration == 0 {
		duration = time.Minute // Default to x requests per minute
	}
	return limiter.New(limiter.Config{
		Max:        count,
		Expiration: duration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() + "_" + c.Path() // Limit each IP to a unique request per path
		},
		LimitReached: func(ctx *fiber.Ctx) error {
			return handler.BuildError(ctx, constants.TooManyRequestsCode, fiber.ErrTooManyRequests.Code, nil, true)
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}
