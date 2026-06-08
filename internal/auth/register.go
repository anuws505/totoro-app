// internal/auth/register.go
package auth

import (
	"encoding/json"
	"net/http"
	"totoro-app/internal/core"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type RegisterRequest struct {
  Mobile   string `json:"mobile"`
  Password string `json:"password"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
  var req RegisterRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    core.WriteError(w, http.StatusBadRequest, "Invalid request", "40000")
    return
  }

  // --- CHECK OTP VERIFICATION ---
  if !IsMobileVerified(req.Mobile) {
    core.WriteError(w, http.StatusForbidden, "Verify OTP, please", "40003")
    return
  }

  // 1. ตรวจสอบเบอร์ซ้ำ
  var existingUser User
  if err := core.DB.Where("mobile = ?", req.Mobile).First(&existingUser).Error; err == nil {
    core.WriteError(w, http.StatusConflict,
      "Mobile number already registered", "40009")
    return
  }

  // 2. Hash Password ก่อนเก็บลง DB
  hashedPassword, err := HashPassword(req.Password)
  if err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Could not process password", "50001")
    return
  }

  // 3. สร้าง User Object ใหม่
  userName := "User-Unknown"
  if len(req.Mobile) >= 4 {
    userName = "User-" + req.Mobile[len(req.Mobile)-4:]
  }

  newID, err := gonanoid.New(12)
  if err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Internal server error", "50003")
    return
  }

  newUser := User{
    UserID:   newID,
    Mobile:   req.Mobile,
    UserName: userName,
    Password: hashedPassword,
    Role:     "USER",
    IsActive: true,
  }

  // 4. บันทึกลง Database
  if err := core.DB.Create(&newUser).Error; err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Could not create user", "50002")
    return
  }

  // --- ลบสิทธิ์ Verified ทิ้งทันทีที่สมัครสำเร็จ ---
  core.RDB.Del(core.Ctx, VERIFY_PREFIX + req.Mobile)

  // 5. ส่ง Response กลับ (อาจจะส่ง Token ให้เลยเพื่อให้ Login อัตโนมัติ)
  token, _ := CreateToken(newUser)

  // ใช้ TokenDuration expired from auth package
  core.RDB.Set(core.Ctx, "active_token:"+newUser.UserID, token, TokenDuration)
  core.RDB.Set(core.Ctx, "user_role:"+newUser.UserID, newUser.Role, TokenDuration)

  // 6. Return response
  core.WriteSuccess(w, http.StatusCreated,
    "User registered successfully", "20001",
    map[string]interface{}{
      "token": token,
      "user": map[string]string{
        "user_id":   newUser.UserID,
        "user_name": newUser.UserName,
        "mobile":    newUser.Mobile,
      },
    },
  )
}
