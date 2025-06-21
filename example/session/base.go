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

package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"

	"github.com/burapha44/example/config"
)

type SessionManager struct {
	Store *session.Store
}

func NewSessionManager() *SessionManager {
	storage := redis.New(redis.Config{
		Host:     config.Conf.RedisHost,
		Port:     config.Conf.RedisPort,
		Password: config.Conf.RedisPassword,
		Database: 2,
		Reset:    false,
	})

	store := session.New(session.Config{
		Storage:   storage,
		KeyLookup: "cookie:sid",
	})

	return &SessionManager{Store: store}
}

func (sm *SessionManager) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := sm.Store.Get(c)
		if err != nil {
			return err
		}
		c.Locals("session", sess)
		return c.Next()
	}
}

// GetSession ดึง session ได้ง่ายขึ้น
func GetSession(c *fiber.Ctx) (*session.Session, error) {
	sess, ok := c.Locals("session").(*session.Session)
	if !ok {
		return nil, fiber.ErrInternalServerError
	}
	return sess, nil
}
