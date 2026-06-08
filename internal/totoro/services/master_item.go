// internal/totoro/services/master_item.go
package services

import (
	"encoding/json"
	"net/http"
	"totoro-app/internal/core"
	"totoro-app/internal/totoro/models"
)

func HandleGetMasterItems(w http.ResponseWriter, r *http.Request) {
  var items []models.MasterItem

  // ตรวจสอบ Query Param ว่าอยากดูตัวที่ลบไปแล้วด้วยไหม (เช่น ?with_deleted=true)
  withDeleted := r.URL.Query().Get("with_deleted")

  query := core.DB
  if withDeleted == "true" {
    query = query.Unscoped() // ดึงมาหมดรวมตัวที่ถูก Soft Delete
  }

  // ดึงข้อมูลจากฐานข้อมูล
  if err := query.Order("created_at DESC").Find(&items).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Failed to retrieve items", "50001")
    return
  }

  // ถ้าไม่พบข้อมูลเลย (เป็น Slice ว่าง)
  if len(items) == 0 {
    core.WriteSuccess(w, http.StatusOK,
      "No items found", "20001", []models.MasterItemResponse{})
    return
  }

  // แปลง Slice ของ MasterItem เป็น MasterItemResponse โดยใช้ ToResponse() ที่เราทำไว้
  responses := make([]models.MasterItemResponse, 0, len(items))
  for _, item := range items {
    responses = append(responses, item.ToResponse())
  }

  // ส่ง Response กลับ
  core.WriteSuccess(w, http.StatusOK,
    "Items retrieved successfully", "20000", responses)
}

func HandleCreateMasterItem(w http.ResponseWriter, r *http.Request) {
  // รับ Input ที่ตรงกับ MasterItemResponse เพราะเราต้องการฟิลด์เดียวกัน
  var input models.CreateOrPatchItemInput
  if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
    core.WriteError(w, http.StatusBadRequest, "Invalid request", "40000")
    return
  }

  if input.ItemSku == "" {
    core.WriteError(w, http.StatusBadRequest, "item_sku is required", "40001")
    return
  }

  // สร้าง Object สำหรับบันทึกลง DB
  newItem := models.MasterItem{
    ItemSku:    input.ItemSku,
    ItemName:   input.ItemName,
    Type:       input.Type,
    Code:       input.Code,
    CodeValue:  input.CodeValue,
    IsFeatured: input.IsFeatured,
  }

  // บันทึกลง Database
  if err := core.DB.Create(&newItem).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Failed to create item (ItemID might already exist)", "50002")
    return
  }

  // Return response
  core.WriteSuccess(w, http.StatusCreated,
    "Master Item created successfully", "20100", newItem.ToResponse())
}

type UpdateItemMasterRequest struct {
  ItemSku    string           `json:"item_sku"`
  ItemName   *string          `json:"item_name"`
  Type       *models.ItemType `json:"type"`
  Code       *string          `json:"code"`
  CodeValue  *string          `json:"code_value"`
  IsFeatured *bool            `json:"is_featured"`
}

func HandleUpdateMasterItem(w http.ResponseWriter, r *http.Request) {
  var req UpdateItemMasterRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    core.WriteError(w, http.StatusBadRequest, "Invalid request", "40000")
    return
  }

  // 1. ค้นหาไอเทมจาก item_sku
  var item models.MasterItem
  if err := core.DB.Where("item_sku = ?", req.ItemSku).First(&item).Error; err != nil {
    core.WriteError(w, http.StatusNotFound, "Item not found", "40402")
    return
  }

  // 2. เตรียมข้อมูลอัปเดต
  updateData := map[string]interface{}{}
  if req.ItemName != nil {
    updateData["item_name"] = *req.ItemName
  }

  if req.Type != nil {
    typeJson, err := json.Marshal(req.Type)
    if err == nil {
      updateData["type"] = string(typeJson)
    }
  }

  if req.Code != nil {
    updateData["code"] = *req.Code
  }

  if req.CodeValue != nil {
    updateData["code_value"] = *req.CodeValue
  }

  if req.IsFeatured != nil {
    updateData["is_featured"] = *req.IsFeatured
  }

  if err := core.DB.Model(&item).Updates(updateData).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Failed to update master item", "50005")
    return
  }

  // ดึงข้อมูลล่าสุดจาก DB (ToResponse) คือค่าล่าสุดจริงๆ
  core.DB.First(&item, item.ID)

  // 3. ส่งข้อมูลที่อัปเดตแล้วกลับไป (ใช้ ToResponse เพื่อความสวยงาม)
  core.WriteSuccess(w, http.StatusOK,
    "Master Item updated successfully", "20004", item.ToResponse())
}

func HandleDeleteMasterItem(w http.ResponseWriter, r *http.Request) {
  // 1. รับ ID จาก URL Parameter (สมมติว่าใช้ /items/{id})
  itemSku := r.URL.Query().Get("item_sku")
  if itemSku == "" {
    core.WriteError(w, http.StatusBadRequest, "item_sku is required", "40001")
    return
  }

  // 2. ค้นหาไอเทมก่อน (เพื่อให้แน่ใจว่ามีอยู่จริงและยังไม่ถูกลบ)
  var item models.MasterItem
  if err := core.DB.Where("item_sku = ?", itemSku).First(&item).Error; err != nil {
    core.WriteError(w, http.StatusNotFound,
      "Item not found or already deleted", "40402")
    return
  }

  // 3. สั่ง Delete (Soft Delete)
  // GORM จะเห็นว่ามีฟิลด์ DeletedAt อยู่ มันจะอัปเดต timestamp ให้เองอัตโนมัติ
  if err := core.DB.Delete(&item).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Failed to delete item", "50006")
    return
  }

  core.WriteSuccess(w, http.StatusOK,
    "Item deleted successfully (Soft delete)", "20005", nil)
}

func HandleRestoreMasterItem(w http.ResponseWriter, r *http.Request) {
  // รับ item_sku จาก Query String หรือ Body ก็ได้
  itemSku := r.URL.Query().Get("item_sku")
  if itemSku == "" {
    core.WriteError(w, http.StatusBadRequest, "item_sku is required", "40001")
    return
  }

  // 1. ค้นหาไอเทมโดยใช้ Unscoped() เพื่อหาตัวที่ deleted_at != null
  var item models.MasterItem
  if err := core.DB.Unscoped().Where("item_sku = ?", itemSku).First(&item).Error; err != nil {
    core.WriteError(w, http.StatusNotFound, "Item not found in archive", "40004")
    return
  }

  // 2. ตรวจสอบว่าไอเทมถูกลบไปจริงไหม (ถ้า deleted_at เป็น null อยู่แล้วไม่ต้องทำอะไร)
  if !item.DeletedAt.Valid {
    core.WriteError(w, http.StatusBadRequest, "Item is already active", "40009")
    return
  }

  // 3. ทำการ Restore (อัปเดต DeletedAt ให้กลับเป็น NULL)
  // ต้องใช้ Unscoped() ควบคู่ไปกับ Model ด้วยเพื่อให้ GORM อนุญาตให้แตะฟิลด์ DeletedAt
  if err := core.DB.Unscoped().Model(&item).Update("deleted_at", nil).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError, "Failed to restore item", "50009")
    return
  }

  core.WriteSuccess(w, http.StatusOK,
    "Item restored successfully", "20006", item.ToResponse())
}
