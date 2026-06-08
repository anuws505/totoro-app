// internal/totoro/workers/sheet_prize.go
package workers

import (
	"log"
	"time"
	"totoro-app/internal/core"
	"totoro-app/internal/totoro/models"

	"gorm.io/gorm"
)

func DoCalculatePrizes() {
  // 1. โหลด Location ไทย
  loc, _ := time.LoadLocation("Asia/Bangkok")
  log.Println("[Worker] Prize Calculator Worker started...")

  // 2. ตั้ง Ticker ให้ตื่นทุก 3 นาที
  ticker := time.NewTicker(3 * time.Minute)
  defer ticker.Stop()

  // รันเช็คทันที 1 ครั้งตอนเปิด Server (Optional)
  checkAndProcess(loc)

  for range ticker.C {
    checkAndProcess(loc)
  }
}

// แยก Logic การเช็คเวลาออกมาเพื่อให้โค้ดสะอาด
func checkAndProcess(loc *time.Location) {
  now := time.Now().In(loc)

  // เงื่อนไข: เฉพาะช่วงเวลา 15:00 ถึง 15:15 น.
  // 15:00:00 จนถึง 15:15:59
  if now.Hour() == 15 && (now.Minute() >= 0 && now.Minute() <= 15) {
    log.Printf("[Worker] It's prize time (%02d:%02d)! Processing...", now.Hour(), now.Minute())
    processPrizes()
  } else {
    // (Optional) Log เพื่อดูว่า Worker ยังตื่นอยู่แต่ไม่ทำอะไร
    // log.Printf("[Worker] Sleeping... current time: %02d:%02d", now.Hour(), now.Minute())
  }
}

func processPrizes() {
  var sheets []models.TotoroSheet
  now := time.Now().UTC()

  // 1. จองคิวสลับสถานะจาก PAID -> REWARD_PROCESSING แบบจำกัดทีละ 100 ใบ เพื่อความปลอดภัย (Concurrency Safe)
  subQuery := core.DB.Model(&models.TotoroSheet{}).
    Select("id").
    Where("status = ? AND reward_at <= ?", "PAID", now).
    Limit(100)

  lockResult := core.DB.Model(&models.TotoroSheet{}).
    Where("id IN (?)", subQuery).
    Update("status", models.StatusRewardProcessing)

  if lockResult.Error != nil {
    log.Printf("[Prize Worker] Error locking sheets: %v", lockResult.Error)
    return
  }

  if lockResult.RowsAffected == 0 {
    return
  }

  // 2. ดึงรายการที่จองสิทธิ์ประมวลผลสำเร็จมาไล่ตรวจรางวัล
  err := core.DB.Where("status = ?", models.StatusRewardProcessing).Find(&sheets).Error
  if err != nil {
    log.Printf("[Prize Worker] Error fetching processing sheets: %v", err)
    return
  }

  for _, sheet := range sheets {
    err := core.DB.Transaction(func(tx *gorm.DB) error {
      // 3. ค้นหาผลรางวัลจากตาราง Master โดยบังคับเปรียบเทียบเป็นค่าเวลามาตรฐาน UTC
      var currentResult models.MasterResult
      targetTimeUTC := sheet.RewardAt.UTC().Format("2006-01-02 15:04:05")

      err := tx.Where("is_active = ? AND date_reward = ?", true, targetTimeUTC).
        First(&currentResult).Error

      // หากรอบรางวัลแอดมินยังไม่ได้ประกาศผล ให้คืนสถานะกลับไปรอตรวจรอบถัดไป
      if err != nil || len(currentResult.Results) == 0 {
        tx.Model(&sheet).Update("status", "PAID")
        return nil
      }

      // 4. เริ่มตรวจรางวัลจากภายใน Items
      var totalReward float64 = 0
      hasPrize := false

      for i := range sheet.Items {
        itemCode := sheet.Items[i].Code // เช่น "DD" หรือ "TTT"

        // ตรวจสอบว่า Master ผลรางวัล มีการออกรางวัลรหัสนี้หรือไม่
        winningValue, exist := currentResult.Results[itemCode]

        // เทียบเลขที่ลูกค้าถือ (CodeValue) กับ เลขรางวัลที่ออก (winningValue)
        if exist && sheet.Items[i].CodeValue == winningValue {
          sheet.Items[i].Prize = true

          var multiplier float64
          if sheet.Items[i].Multiple == models.MultipleDirect {
              multiplier = sheet.Items[i].Type.Direct
          } else if sheet.Items[i].Multiple == models.MultipleMixed {
              multiplier = sheet.Items[i].Type.Mixed
          }

          // รางวัล = ราคาที่ซื้อ x ตัวคูณคูปอง (เช่น 10 * 25 = 250)
          calculatedReward := sheet.Items[i].Price * multiplier

          // อัปเดตยอดเงินรางวัลกลับเข้าไปเก็บใน Item ตัวนั้น
          sheet.Items[i].Reward = calculatedReward

          // สะสมยอดรางวัลรวมของสลิปใบนี้
          totalReward += calculatedReward
          hasPrize = true
        } else {
          sheet.Items[i].Prize = false
          sheet.Items[i].Reward = 0
        }
      }

      // 5. สรุปผลยอดบันทึกความเปลี่ยนแปลงลงตารางหลัก
      sheet.Totals.Reward = totalReward
      sheet.SheetCodeDD = currentResult.Results["DD"]
      sheet.SheetCodeTTT = currentResult.Results["TTT"]
      sheet.Status = models.StatusRewardCompleted
      sheet.RewardCompletedAt = &now

      // บันทึกความเปลี่ยนแปลงข้อมูลรวมลงก้อน JSON สลับกลับเข้าฐานข้อมูล
      if err := tx.Save(&sheet).Error; err != nil {
        return err
      }

      log.Printf("[Prize Worker] Sheet %s processing finalized. Won: %t, Total Reward: %.2f", sheet.SheetID, hasPrize, totalReward)
      return nil
    })

    if err != nil {
      log.Printf("[Prize Worker] Failed executing calculation transaction %s: %v", sheet.SheetID, err)
    }
  }
}
