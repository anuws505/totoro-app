// internal/totoro/models/promocode.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type PromotionCode struct {
  gorm.Model
  Code       string     `gorm:"uniqueIndex;type:varchar(20)" json:"code"`
  Value      float64    `json:"value"`
  MaxUsage   int        `json:"max_usage"`
  UsedCount  int        `gorm:"default:0" json:"used_count"`
  ExpiryDate *time.Time `json:"expiry_date"`
  IsFeatured bool       `gorm:"column:is_featured;default:false"`
}
