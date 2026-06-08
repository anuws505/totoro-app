// internal/auth/user_handle.go
package auth

import (
	"encoding/json"
	"net/http"
	"time"
	"totoro-app/internal/core"
)

type UpdateUsernameRequest struct {
  NewUserName string `json:"new_user_name"`
}

func HandleUpdateUsername(w http.ResponseWriter, r *http.Request) {
  // 1. ดึง userID จาก Context (ที่ถูกฉีดเข้ามาโดย AuthMiddleware)
  userID, ok := r.Context().Value("userID").(string)
  if !ok {
    core.WriteError(w, http.StatusUnauthorized, "Unauthorized", "40001")
    return
  }

  // 2. รับข้อมูลชื่อใหม่จาก Body
  var req UpdateUsernameRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.NewUserName == "" {
    core.WriteError(w, http.StatusBadRequest, "Invalid name", "40000")
    return
  }

  // 3. อัปเดตข้อมูลใน Database
  result := core.DB.Model(&User{}).
    Where("user_id = ?", userID).
    Update("user_name", req.NewUserName)

  if result.Error != nil {
    core.WriteError(w, http.StatusInternalServerError, "Database error", "50000")
    return
  }

  if result.RowsAffected == 0 {
    core.WriteError(w, http.StatusNotFound, "User not found", "40004")
    return
  }

  // 4. ส่ง Success Response
  core.WriteSuccess(w, http.StatusOK, "Username updated successfully", "20002",
    map[string]string{
      "user_id":   userID,
      "user_name": req.NewUserName,
    },
  )
}

type UpdateRoleRequest struct {
  TargetUserID string `json:"target_user_id"` // ID ของคนที่จะถูกเปลี่ยนสิทธิ์
  NewRole      string `json:"new_role"`       // "ADMIN" หรือ "USER"
}

func HandleUpdateUserRole(w http.ResponseWriter, r *http.Request) {
  var req UpdateRoleRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    core.WriteError(w, http.StatusBadRequest, "Invalid request", "40000")
    return
  }

  // 1. ตรวจสอบว่ามี User นี้อยู่จริงไหมก่อนอัปเดต
  var user User
  if err := core.DB.Where("user_id = ?", req.TargetUserID).First(&user).Error; err != nil {
    core.WriteError(w, http.StatusNotFound, "Target user not found", "40401")
    return
  }

  // 2. Update ใน Database
  if err := core.DB.Model(&User{}).Where("user_id = ?", req.TargetUserID).
  Update("role", req.NewRole).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Failed to update role in database", "50004")
    return
  }

  // 3. อัปเดตใน Redis ทันที (เพื่อให้สิทธิ์ใหม่มีผลทันทีไม่ต้องรอ Re-login)
  // เราใช้ Key เดียวกับที่ AuthMiddleware ตรวจสอบ
  core.RDB.Set(core.Ctx, "user_role:"+req.TargetUserID, req.NewRole, TokenDuration)

  // 4. Return response
  core.WriteSuccess(w, http.StatusOK, "User role updated successfully", "20002",
    map[string]string{
      "user_id":  req.TargetUserID,
      "new_role": req.NewRole,
    },
  )
}

type ChangePasswordRequest struct {
  OldPassword string `json:"old_password"`
  NewPassword string `json:"new_password"`
}

func HandleChangePassword(w http.ResponseWriter, r *http.Request) {
  // 1. ดึง userID จาก Context (ที่ AuthMiddleware แกะให้)
  userID, ok := r.Context().Value("userID").(string)
  if !ok {
    core.WriteError(w, http.StatusUnauthorized, "Unauthorized", "40001")
    return
  }

  var req ChangePasswordRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    core.WriteError(w, http.StatusBadRequest, "Invalid request", "40000")
    return
  }

  // 2. ค้นหา User ใน DB
  var user User
  if err := core.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
    core.WriteError(w, http.StatusNotFound, "User not found", "40004")
    return
  }

  // 3. ตรวจสอบรหัสผ่านเก่าว่าถูกต้องไหม
  if !CheckPasswordHash(req.OldPassword, user.Password) {
    core.WriteError(w, http.StatusUnauthorized,
      "Old password is incorrect", "40010")
    return
  }

  // 4. Hash รหัสผ่านใหม่
  newHashedPassword, err := HashPassword(req.NewPassword)
  if err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Failed password error", "50004")
    return
  }

  // 5. อัปเดต DB (Password และ PasswordChangedAt)
  now := time.Now()
  updates := map[string]interface{}{
    "password":            newHashedPassword,
    "password_changed_at": &now,
  }

  if err := core.DB.Model(&user).Updates(updates).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Failed update password", "50000")
    return
  }

  // 6. SECURITY STEP: ลบ Token และ Role ใน Redis ทันที
  // เพื่อให้ User ต้อง Log in ใหม่ด้วยรหัสผ่านใหม่ (ป้องกันกรณีคนอื่นแอบเปลี่ยน)
  core.RDB.Del(core.Ctx, "active_token:"+userID)
  core.RDB.Del(core.Ctx, "user_role:"+userID)

  // 7. Return response
  core.WriteSuccess(w, http.StatusOK,
    "Password changed successfully. Please login again.", "20003", nil)
}
