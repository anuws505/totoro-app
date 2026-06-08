// internal/auth/login.go
package auth

import (
	"encoding/json"
	"net/http"
	"totoro-app/internal/core"
)

type LoginRequest struct {
  Mobile   string `json:"mobile"`
  Password string `json:"password"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
  var req LoginRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    core.WriteError(w, http.StatusBadRequest, "Invalid request", "40000")
    return
  }

  var user User
  // 1. ค้นหา User จากเบอร์มือถือ
  err := core.DB.Where("mobile = ? AND is_active = ?", req.Mobile, true).First(&user).Error
  if err != nil {
    core.WriteError(w, http.StatusUnauthorized, "Invalid mobile or password", "40004")
    return
  }

  // 2. ตรวจสอบรหัสผ่านที่ Hash ไว้ด้วย bcrypt
  if !CheckPasswordHash(req.Password, user.Password) {
    core.WriteError(w, http.StatusUnauthorized, "Invalid mobile or password", "40001")
    return
  }

  // 3. สร้าง Token โดยใช้ UserID จริงจากฐานข้อมูล
  token, err := CreateToken(user)
  if err != nil {
    core.WriteError(w, http.StatusInternalServerError, "Could not create token", "50000")
    return
  }

  // ใช้ TokenDuration expired from auth package
  core.RDB.Set(core.Ctx, "active_token:"+user.UserID, token, TokenDuration)
  core.RDB.Set(core.Ctx, "user_role:"+user.UserID, user.Role, TokenDuration)

  // 4. ส่ง Token และข้อมูลเบื้องต้นกลับไป
  core.WriteSuccess(w, http.StatusOK, "Login success", "20000",
    map[string]interface{}{
      "token": token,
      "user": map[string]string{
        "user_id":   user.UserID,
        "user_name": user.UserName,
        "mobile":    user.Mobile,
      },
    },
  )
}
