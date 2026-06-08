// internal/auth/model.go
package auth

import (
	"time"

	"gorm.io/gorm"
)

const (
  RoleUser  = "USER"
  RoleAdmin = "ADMIN"
)

type User struct {
  gorm.Model
  ID                uint       `gorm:"primaryKey" json:"-"`
  UserID            string     `gorm:"uniqueIndex;type:varchar(12)" json:"user_id"`
  Mobile            string     `gorm:"index;type:varchar(15)" json:"mobile"`
  UserName          string     `json:"user_name"`
  Password          string     `json:"-"`
  Role              string     `gorm:"type:varchar(10);default:'USER'"`
  IsActive          bool       `gorm:"default:true" json:"is_active"`
  DeactivatedAt     *time.Time `json:"deactivated_at"`
  PasswordChangedAt *time.Time `json:"password_changed_at"`
}
