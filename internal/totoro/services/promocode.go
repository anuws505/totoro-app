// internal/totoro/services/promocode.go
package services

import (
	"encoding/json"
	"net/http"
	"time"
	"totoro-app/internal/core"
	"totoro-app/internal/totoro/models"
)

func HandleGetPromotions(w http.ResponseWriter, r *http.Request) {
  var promos []models.PromotionCode

  // ตรวจสอบ Query Param ว่าอยากดูตัวที่ลบไปแล้วด้วยไหม (เช่น ?with_deleted=true)
  withDeleted := r.URL.Query().Get("with_deleted")

  query := core.DB
  if withDeleted == "true" {
    query = query.Unscoped() // ดึงมาหมดรวมตัวที่ถูก Soft Delete
  }

  // ดึงข้อมูลจากฐานข้อมูล
  if err := query.Order("created_at DESC").Find(&promos).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Failed to retrieve promotions", "50001")
    return
  }

  // ถ้าไม่พบข้อมูลเลย (เป็น Slice ว่าง)
  if len(promos) == 0 {
    core.WriteSuccess(w, http.StatusOK,
      "No promotions found", "20001", []models.PromotionCode{})
    return
  }

  core.WriteSuccess(w, http.StatusOK,
    "Promotions retrieved successfully", "20000", promos)
}

func HandleCreatePromotion(w http.ResponseWriter, r *http.Request) {
  var promo models.PromotionCode
  if err := json.NewDecoder(r.Body).Decode(&promo); err != nil {
    core.WriteError(w, http.StatusBadRequest, "Invalid request", "40000")
    return
  }

  if promo.Code == "" {
    core.WriteError(w, http.StatusBadRequest,
      "Promotion 'code' is required", "40001")
    return
  }

  if err := core.DB.Create(&promo).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Failed to create promotion (Code might exist)", "50002")
    return
  }

  core.WriteSuccess(w, http.StatusCreated,
    "Promotion created successfully", "20001", promo)
}

type UpdatePromotionRequest struct {
  Code       string     `json:"code"`
  Value      *float64   `json:"value"`
  MaxUsage   *int       `json:"max_usage"`
  ExpiryDate *time.Time `json:"expiry_date"`
  IsFeatured *bool      `json:"is_featured"`
}

func HandleUpdatePromotion(w http.ResponseWriter, r *http.Request) {
  var req UpdatePromotionRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    core.WriteError(w, http.StatusBadRequest, "Invalid request", "40000")
    return
  }

  // 1. ค้นหา promotion จาก code
  var promo models.PromotionCode
  if err := core.DB.Where("code = ?", req.Code).First(&promo).Error; err != nil {
    core.WriteError(w, http.StatusNotFound, "Promtion not found", "40402")
    return
  }

  // 2. เตรียมข้อมูลอัปเดต
  updateData := map[string]interface{}{}
  if req.Value != nil {
    updateData["value"] = *req.Value
  }
  if req.MaxUsage != nil {
    updateData["max_usage"] = *req.MaxUsage
  }
  if req.ExpiryDate != nil {
    updateData["expiry_date"] = req.ExpiryDate
  }
  if req.IsFeatured != nil {
    updateData["is_featured"] = *req.IsFeatured
  }

  // 3. บันทึกการเปลี่ยนแปลง
  if err := core.DB.Model(&promo).Updates(updateData).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Failed update promotion", "50000")
    return
  }

  // 4. Refresh ข้อมูลก่อนส่งกลับ
  core.DB.First(&promo, promo.ID)

  core.WriteSuccess(w, http.StatusOK,
    "Promotion updated successfully", "20002", promo)
}

func HandleDeletePromotion(w http.ResponseWriter, r *http.Request) {
  code := r.URL.Query().Get("code")

  if code == "" {
    core.WriteError(w, http.StatusBadRequest,
      "Promotion 'code' is required", "40001")
    return
  }

  var promo models.PromotionCode
  if err := core.DB.Where("code = ?", code).First(&promo).Error; err != nil {
    core.WriteError(w, http.StatusNotFound, "Promotion not found", "40004")
    return
  }

  if err := core.DB.Delete(&promo).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Failed to delete promotion", "50006")
    return
  }

  core.WriteSuccess(w, http.StatusOK,
    "Promotion deleted successfully (Soft delete)", "20005", nil)
}
