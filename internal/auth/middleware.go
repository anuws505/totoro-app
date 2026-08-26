// internal/auth/middleware.go
package auth

import (
	"context"
	"net/http"
	"strings"
	"totoro-app/internal/core"

	"github.com/golang-jwt/jwt/v5"
)

// สร้าง Custom Context Key Type เพื่อความปลอดภัย
type contextKey string

const (
  UserIDKey contextKey = "userID"
  RoleKey   contextKey = "role"
)

// AuthMiddleware - ใช้เป็น gorilla/mux Middleware
func AuthMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
      core.WriteError(w, http.StatusUnauthorized, "No token provided", "40101")
      return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")

    token, err := jwt.ParseWithClaims(
      tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return JwtSecret, nil
      })

    if err != nil || !token.Valid {
      core.WriteError(w, http.StatusUnauthorized, "Invalid or expired token", "40102")
      return
    }

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
        var user User
        if err := core.DB.Select("role").Where("user_id = ?", claims.UserID).First(&user).Error; err != nil {
          core.WriteError(w, http.StatusUnauthorized, "No permissions", "40104")
          return
        }
        role = user.Role
        core.RDB.Set(core.Ctx, "user_role:"+claims.UserID, role, TokenDuration)
      }

      // ส่งผ่าน Context
      ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
      ctx = context.WithValue(ctx, RoleKey, role)

      next.ServeHTTP(w, r.WithContext(ctx))
      return
    }

    core.WriteError(w, http.StatusUnauthorized, "Unauthorized", "40100")
  })
}

// AdminMiddleware - ใช้เป็น gorilla/mux Middleware สำหรับ Admin
func AdminMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    role, ok := r.Context().Value(RoleKey).(string)

    if !ok || role != "ADMIN" {
      core.WriteError(w, http.StatusForbidden, "Access denied: Admin privileges required", "40301")
      return
    }

    next.ServeHTTP(w, r)
  })
}
