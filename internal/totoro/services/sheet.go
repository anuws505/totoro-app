// internal/totoro/services/sheet.go
package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"totoro-app/internal/auth"
	"totoro-app/internal/core"
	"totoro-app/internal/totoro/models"

	"gorm.io/gorm"
)

func HandleCreateSheet(w http.ResponseWriter, r *http.Request) {
  userID, ok := r.Context().Value("userID").(string)
  if !ok {
    core.WriteError(w, http.StatusUnauthorized, "Unauthorized user", "40101")
    return
  }

  // 1. รับ Input
  var input struct {
    Items []struct {
      ItemSku  string  `json:"item_sku"`
      Price    float64 `json:"price"`
      Multiple string  `json:"multiple"`
    } `json:"items"`
    PromotionCode string `json:"promotion_code"`
  }

  if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
    core.WriteError(w, http.StatusBadRequest, "Invalid request", "40000")
    return
  }

  // 2. Validate & Snapshot Items
  var validatedItems []models.SheetItem
  for _, incoming := range input.Items {
    var master models.MasterItem
    if err := core.DB.Where("item_sku = ?", incoming.ItemSku).First(&master).Error; err != nil {
      core.WriteError(w, http.StatusNotFound,
        fmt.Sprintf("Item %s not found", incoming.ItemSku), "40402")
      return
    }

    sheetItem := models.SheetItem{
      ItemSku:    master.ItemSku,
      ItemName:   master.ItemName,
      Price:      incoming.Price,
      Type:       master.Type,
      Multiple:   incoming.Multiple,
      Code:       master.Code,
      CodeValue:  master.CodeValue,
      IsFeatured: master.IsFeatured,
      // หมายเหตุ: Reward/Prize จะถูกคำนวณใน BeforeSave ของ Sheet
    }
    validatedItems = append(validatedItems, sheetItem)
  }

  // 3. ตรวจสอบ Draft เดิมก่อนเพื่อจัดการเรื่องคูปอง
  var sheet models.TotoroSheet
  hasDraft := core.DB.Where("user_id = ? AND status = ?", userID, models.StatusDraft).First(&sheet).Error == nil

  // 4. ตรวจสอบและจอง Promotion Code (ถ้ามีการส่งมา)
  var promo models.PromotionCode
  if input.PromotionCode != "" {
    // ถ้ามี Draft และโค้ดที่ส่งมา "ไม่เหมือนเดิม"
    if hasDraft && sheet.Totals.Discount.Code != "" && sheet.Totals.Discount.Code != input.PromotionCode {
      // คืนสิทธิ์ให้โค้ดเก่าก่อน
      core.DB.Model(&models.PromotionCode{}).
        Where("code = ?", sheet.Totals.Discount.Code).
        Update("used_count", gorm.Expr("GREATEST(used_count - 1, 0)"))
    }

    // ทำการจองโค้ดใหม่ (Logic เดิมของคุณ)
    if !hasDraft || sheet.Totals.Discount.Code != input.PromotionCode {
      result := core.DB.Model(&models.PromotionCode{}).
        Where("code = ? AND (max_usage = 0 OR used_count < max_usage) AND (expiry_date IS NULL OR expiry_date > ?)",
          input.PromotionCode, time.Now()).
        Update("used_count", gorm.Expr("used_count + 1"))

      if result.RowsAffected > 0 {
        core.DB.Where("code = ?", input.PromotionCode).First(&promo)
      } else {
        core.WriteError(w, http.StatusBadRequest, "Promotion code invalid or full", "40005")
        return
      }
    } else {
      core.DB.Where("code = ?", input.PromotionCode).First(&promo)
    }
  }

  var user auth.User
  userName := "Guest"
  if err := core.DB.Where("user_id = ?", userID).First(&user).Error; err == nil {
    userName = user.UserName
  }

  // 5. จัดการบันทึกข้อมูล
  sheet.Items = validatedItems
  sheet.Totals.Discount = models.Discount{Code: promo.Code, Value: promo.Value}
  sheet.UserID = userID
  sheet.UserName = userName

  if hasDraft {
    // อัปเดตใบเดิม
    if err := core.DB.Save(&sheet).Error; err != nil {
      core.WriteError(w, http.StatusInternalServerError, "Update error", "50003")
      return
    }
  } else {
    // สร้างใบใหม่
    sheet.SheetID = GenerateSheetID()
    sheet.Status = models.StatusDraft
    expires := time.Now().AddDate(0, 0, 1)
    sheet.ExpiresAt = &expires

    // ค้นหารอบรางวัลถัดไปที่แอดมินสร้างเปิดระบบไว้ล่วงหน้า (Active และยังไม่ออกผล)
    var upcomingResult models.MasterResult
    err := core.DB.Where("is_active = ? AND date_reward > ?", true, time.Now().UTC()).
      Order("date_reward ASC"). // เอารอบที่ใกล้ที่สุด
      First(&upcomingResult).Error

    if err != nil {
      // เผื่อกรณีแอดมินลืมสร้างรอบทิ้งไว้ล่วงหน้า ให้ fallback ใช้ระบบบ่ายสามวันถัดไปแทน ป้องกันระบบพัง
      loc, _ := time.LoadLocation("Asia/Bangkok")
      now := time.Now().In(loc)
      cutoffToday := time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, loc)
      if now.Before(cutoffToday) {
        sheet.RewardAt = cutoffToday.UTC()
      } else {
        sheet.RewardAt = cutoffToday.AddDate(0, 0, 1).UTC()
      }
    } else {
      sheet.RewardAt = upcomingResult.DateReward
    }

    if err := core.DB.Create(&sheet).Error; err != nil {
      core.WriteError(w, http.StatusInternalServerError, "Create error", "50002")
      return
    }
  }

  // 6. Response
  core.WriteSuccess(w, http.StatusCreated, "Sheet processed successfully", "20100", sheet)
}

func GenerateSheetID() string {
  now := time.Now()
  dateStr := now.Format("060102")

  micro := now.UnixMicro()
  shortID := strconv.FormatInt(micro, 36)

  return fmt.Sprintf("SHEET-%s-%s", dateStr, shortID)
}
