// internal/totoro/workers/sheet_expires.go
package workers

import (
	"log"
	"time"
	"totoro-app/internal/core"
	"totoro-app/internal/totoro/models"

	"gorm.io/gorm"
)

func DoCleanupExpiredDrafts() {
  // 1. สร้าง Ticker
  ticker := time.NewTicker(30 * time.Minute)

  // 2. ป้องกัน Memory Leak (แม้จะรันยาวก็ตาม)
  defer ticker.Stop()

  // 3. รันทันที 1 ครั้งตอนเปิด Server (เพราะ Ticker จะรอรอบแรก 30 นาที)
  CleanupExpiredDrafts()

  log.Println("[Worker] Expired Drafts Cleanup Worker started...")

  // 4. Infinite Loop ที่รอรับสัญญาณจาก Ticker
  for range ticker.C {
    CleanupExpiredDrafts()
  }
}

func CleanupExpiredDrafts() {
  now := time.Now().UTC()
  const batchSize = 100 // กำหนดขนาดที่ต้องการ

  // 1. "จอง" เฉพาะ 100 ใบแรกที่หมดอายุ
  // เราใช้ Subquery เพื่อหา ID ของ 100 ใบแรกก่อน แล้วค่อย Update
  subQuery := core.DB.Model(&models.TotoroSheet{}).
    Select("id").
    Where("status = ? AND expires_at < ?", models.StatusDraft, now).
    Limit(batchSize)

  result := core.DB.Model(&models.TotoroSheet{}).
    Where("id IN (?)", subQuery).
    Update("status", models.StatusExpiring)

  if result.Error != nil {
    log.Printf("[Worker] Error locking batch: %v", result.Error)
    return
  }

  if result.RowsAffected == 0 {
    return
  }

  // 2. ดึง 100 ใบที่จองสำเร็จมาเข้า Loop เดิม
  var expiringSheets []models.TotoroSheet
  core.DB.Where("status = ?", models.StatusExpiring).Find(&expiringSheets)

  log.Printf("[Worker] Processing batch of %d sheets", len(expiringSheets))

  for _, sheet := range expiringSheets {
    // ... (Logic Transaction คืนคูปอง) ...
    err := core.DB.Transaction(func(tx *gorm.DB) error {
      // คืนสิทธิ์ Promotion
      if sheet.Totals.Discount.Code != "" {
        tx.Model(&models.PromotionCode{}).
          Where("code = ?", sheet.Totals.Discount.Code).
          Update("used_count", gorm.Expr("GREATEST(used_count - 1, 0)"))
      }
      return tx.Model(&sheet).Update("status", models.StatusExpired).Error
    })

    if err != nil {
      log.Printf("[Worker] Failed on sheet %s: %v", sheet.SheetID, err)
    }
  }
}
