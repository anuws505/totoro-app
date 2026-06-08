// internal/auth/token.go
package auth

import (
	"time"
	"totoro-app/internal/core"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var host = core.GetEnv("AUTH_SECRET_KEY", "a-string-secret-at-least-256-bits-long")
var JwtSecret = []byte(host)

type CustomClaims struct {
  UserID string `json:"user_id"`
  jwt.RegisteredClaims
}

func CreateToken(user User) (string, error) {
  claims := CustomClaims{
    UserID: user.UserID,
    RegisteredClaims: jwt.RegisteredClaims{
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenDuration)),
      IssuedAt:  jwt.NewNumericDate(time.Now()),
    },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString(JwtSecret)
}

// HashPassword แปลงรหัสผ่านเป็น Hash
func HashPassword(password string) (string, error) {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  return string(bytes), err
}

// CheckPasswordHash ตรวจสอบรหัสผ่านตอน Login
func CheckPasswordHash(password, hash string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
  return err == nil
}
