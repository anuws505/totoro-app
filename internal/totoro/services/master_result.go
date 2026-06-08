// internal/totoro/services/master_result.go
package services

import (
	"encoding/json"
	"net/http"
	"time"
	"totoro-app/internal/core"
	"totoro-app/internal/totoro/models"

	"github.com/gorilla/mux"
)

type CreateResultInput struct {
  DateReward string            `json:"date_reward"` // Format: "2026-05-16 15:00:00"
  Results    models.ResultMap  `json:"results"`     // {"DD": "10", "TTT": "160"}
}

func HandleCreateMasterResult(w http.ResponseWriter, r *http.Request) {
  var input CreateResultInput
  if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
    core.WriteError(w, http.StatusBadRequest, "Invalid request", "40000")
    return
  }

  loc, _ := time.LoadLocation("Asia/Bangkok")
  parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", input.DateReward, loc)
  if err != nil {
    core.WriteError(w, http.StatusBadRequest,
      "Invalid date format. Use YYYY-MM-DD HH:MM:SS", "40402")
    return
  }

  masterResult := models.MasterResult{
    DateReward: parsedTime.UTC(),
    Results:    input.Results,
    IsActive:   true,
  }

  if err := core.DB.Create(&masterResult).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Failed to create master result", "50000")
    return
  }

  core.WriteSuccess(w, http.StatusCreated,
    "Master result created successfully", "20001", masterResult)
}

type PatchResultInput struct {
  DateReward *string          `json:"date_reward"`
  Results    models.ResultMap `json:"results"`
  IsActive   *bool            `json:"is_active"`
}

func HandlePatchMasterResult(w http.ResponseWriter, r *http.Request) {
  resultID := mux.Vars(r)["id"]

  if resultID == "" {
    core.WriteError(w, http.StatusBadRequest, "Master result ID is required", "40001")
    return
  }

  // 1. ดึงข้อมูลเดิมจาก DB มารอไว้
  var result models.MasterResult
  if err := core.DB.First(&result, resultID).Error; err != nil {
    core.WriteError(w, http.StatusNotFound, "Master result not found", "40401")
    return
  }

  var input PatchResultInput
  if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
    core.WriteError(w, http.StatusBadRequest, "Invalid request payload", "40000")
    return
  }

  // 2. ตรวจสอบและแทนที่ค่าทีละฟิลด์ลงใน Struct `result` เดิมตรงๆ
  if input.IsActive != nil {
    result.IsActive = *input.IsActive
  }

  // ตรวจสอบ Results: ไม่ใช้ค่า nil เพื่อให้ยังสามารถส่ง {} มาล้างค่าได้
  if input.Results != nil {
    result.Results = input.Results
  }

  if input.DateReward != nil {
    loc, _ := time.LoadLocation("Asia/Bangkok")
    parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", *input.DateReward, loc)
    if err != nil {
      core.WriteError(w, http.StatusBadRequest,
        "Invalid date format. Use YYYY-MM-DD HH:MM:SS", "40402")
      return
    }
    result.DateReward = parsedTime.UTC()
  }

  // 3. ใช้คำสั่ง .Save(&result) เพื่อเซฟ Object Struct ทั้งตัว
  // GORM Serializer จะทำงานร่วมกับขั้นตอนนี้ได้อย่างสมบูรณ์แบบ ไม่มีเออเรอร์แน่นอนครับ
  if err := core.DB.Save(&result).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Failed to update master result", "50000")
    return
  }

  core.WriteSuccess(w, http.StatusOK,
    "Master result updated successfully", "20002", result)
}
