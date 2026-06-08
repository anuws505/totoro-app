// internal/auth/middleware.go
package auth

import (
	"context"
	"net/http"
	"strings"
	"totoro-app/internal/core"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
      core.WriteError(w, http.StatusUnauthorized, "No token provided", "40101")
      return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")

    // 1. ตรวจสอบ Signature และ Expiration ของ JWT ตามปกติ
    token, err := jwt.ParseWithClaims(
      tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
      return JwtSecret, nil
    })

    if err != nil || !token.Valid {
      core.WriteError(w, http.StatusUnauthorized, "Invalid or expired token", "40102")
      return
    }

    // 2. ดึง UserID จาก Claims
    if claims, ok := token.Claims.(*CustomClaims); ok {

      // ตรวจสอบ Active Token (Single Session)
      latestToken, err := core.RDB.Get(core.Ctx, "active_token:"+claims.UserID).Result()

      if err != nil || latestToken != tokenString {
        core.WriteError(w, http.StatusUnauthorized, "Session expired", "40103")
        return
      }

      // ตรวจสอบ Role จาก Redis
      role, err := core.RDB.Get(core.Ctx, "user_role:"+claims.UserID).Result()
      if err != nil {
        // ถ้าใน Redis หาย (เช่นโดนล้าง) ให้ไปดึงจาก DB เป็น Fallback
        var user User
        if err := core.DB.Select("role").Where("user_id = ?", claims.UserID).First(&user).Error; err != nil {
          core.WriteError(w, http.StatusUnauthorized, "No permissions", "40104")
          return
        }
        role = user.Role
        // บันทึกกลับลง Redis อีกครั้ง
        core.RDB.Set(core.Ctx, "user_role:"+claims.UserID, role, TokenDuration)
      }
      // ----------------------------------------

      // ถ้าผ่านทุกด่าน ส่ง userID เข้า Context ต่อไป
      ctx := context.WithValue(r.Context(), "userID", claims.UserID)
      ctx = context.WithValue(ctx, "role", role)

      next.ServeHTTP(w, r.WithContext(ctx))
      return
    }

    core.WriteError(w, http.StatusUnauthorized, "Unauthorized", "40100")
  }
}

func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    // ดึง Role จาก Context (ที่ AuthMiddleware แกะออกมาให้)
    role, ok := r.Context().Value("role").(string)

    if !ok || role != "ADMIN" {
      core.WriteError(w, http.StatusForbidden,
        "Access denied: Admin privileges required", "40301")
      return
    }

    next.ServeHTTP(w, r)
  }
}
