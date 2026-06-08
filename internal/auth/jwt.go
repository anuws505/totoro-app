// internal/auth/jwt.go
package auth

import "time"

// กำหนด time expired เพื่อใช้ใน package auth
const TokenDuration = 24 * time.Hour
