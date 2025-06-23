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

package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"

	"github.com/burapha44/example/constants"
	"github.com/burapha44/example/db"
)

// rollbackCtxTrx rollbacks active database transaction associated with the given Fiber context.
// If no transaction is associated with the context, it does nothing.
func rollbackCtxTrx(ctx *fiber.Ctx) {
	trx, _ := startNewPGTrx(ctx)

	if trx != nil {
		if err := trx.Rollback(ctx.UserContext()); err != nil {
			log.Fatalf("Error rollback transaction: %v", err)
		}
	}
}

// commitCtxTrx commits active database transaction associated with the given Fiber context.
// If no transaction is associated with the context, it does nothing.
// If commit fails, it returns an error response to the client with status code 500 (Internal Server Error).
func commitCtxTrx(ctx *fiber.Ctx) error {
	trx, err := startNewPGTrx(ctx)

	if err != nil {
		return BuildError(ctx, constants.UnableToGetTrxCode, fiber.StatusInternalServerError, err, true)
	}

	if trx != nil {
		if err := trx.Commit(ctx.UserContext()); err != nil {
			return BuildError(ctx, constants.UnableToCommitTrxCode, fiber.StatusInternalServerError, err, true)
		}
	}

	return nil
}

const (
	DbTrxKey = "db_trx_key"
)

// StartNewPGTrx returns a new Postgres transaction associated with the given Fiber context.
// If a transaction is already associated with the context, it is returned instead of creating
// a new one. The transaction is stored in the context under the key DbTrxKey.
func startNewPGTrx(ctx *fiber.Ctx) (pgx.Tx, error) {
	if trx := ctx.Locals(DbTrxKey); trx != nil {
		return trx.(pgx.Tx), nil
	}

	pgTrx, err := db.PGTransaction(ctx.UserContext())

	if err != nil {
		return nil, err
	}

	ctx.Locals(DbTrxKey, pgTrx)

	return pgTrx, nil
}
