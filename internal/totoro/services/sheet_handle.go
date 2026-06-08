// internal/totoro/services/sheet_handle.go
package services

import (
	"net/http"
	"totoro-app/internal/core"
	"totoro-app/internal/totoro/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func HandleCheckout(w http.ResponseWriter, r *http.Request) {
  userID := r.Context().Value("user_id").(string)
  sheetID := mux.Vars(r)["id"]

  var sheet models.TotoroSheet
  if err := core.DB.Where("sheet_id = ? AND user_id = ?", sheetID, userID).First(&sheet).Error; err != nil {
      core.WriteError(w, http.StatusNotFound, "Sheet not found", "40004")
      return
  }

  // 1. ล็อกสถานะก่อนไปจ่ายเงิน
  core.DB.Model(&sheet).Update("status", models.StatusWaitingPayment)

  // 2. เรียก Omise API (ตัวอย่างการสร้าง Charge)
  // ในที่นี้ต้องใช้ Omise Go SDK
  /*
  charge, _ := charge.Create(&omise.CreateCharge{
    Amount:   int64(sheet.Totals.FinalAmount * 100), // สตางค์
    Currency: "thb",
    ReturnURI: "https://your-app.com/orders/" + sheetID, // หน้าที่จะให้กลับมาหลังจากจ่ายเสร็จ
    Source:   input.SourceToken,
  })
  */

  // 3. ส่ง URL หน้าชำระเงินกลับไปให้ Frontend
  // core.WriteJSON(w, http.StatusOK, map[string]string{"authorize_url": charge.AuthorizeURI})
}

func HandleOmiseWebhook(w http.ResponseWriter, r *http.Request) {
  // 1. รับ JSON จาก Omise
  // 2. ตรวจสอบ Event ว่าเป็น "charge.complete" หรือไม่
  // 3. ดึง SheetID จาก metadata ที่เราส่งไปตอนสร้าง Charge

  sheetID := "..." // ดึงจาก payload

  err := core.DB.Transaction(func(tx *gorm.DB) error {
    var sheet models.TotoroSheet
    if err := tx.Where("sheet_id = ?", sheetID).First(&sheet).Error; err != nil {
      return err
    }

    // 4. อัปเดตสถานะเป็น PAID
    if err := tx.Model(&sheet).Update("status", models.StatusPaid).Error; err != nil {
      return err
    }

    return nil
  })

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusOK)
}

