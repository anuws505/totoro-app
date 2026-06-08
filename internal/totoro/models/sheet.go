// internal/totoro/models/sheet.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Status Types สำหรับการจัดการ Workflow
type SheetStatus string

const (
  StatusDraft               SheetStatus = "DRAFT"
  StatusExpiring            SheetStatus = "EXPIRING"
  StatusExpired             SheetStatus = "EXPIRED"
  StatusWaitingPayment      SheetStatus = "WAITING_PAYMENT"
  StatusFailed              SheetStatus = "FAILED"
  StatusPaid                SheetStatus = "PAID"
  StatusRewardProcessing    SheetStatus = "REWARD_PROCESSING"
  StatusRewardCompleted     SheetStatus = "REWARD_COMPLETED"

  StatusMissed              SheetStatus = "MISSED"
  StatusReward              SheetStatus = "REWARD"
  StatusPendingPayment      SheetStatus = "PENDING_PAYMENT"
  StatusPreparePayment      SheetStatus = "PREPARE_PAYMENT"
  StatusPaymentApply        SheetStatus = "PAYMENT_APPLY"
  StatusPaymentProcessing   SheetStatus = "PAYMENT_PROCESSING"
  StatusPaymentCompleted    SheetStatus = "PAYMENT_COMPLETED"
  StatusPaymentFailed       SheetStatus = "PAYMENT_FAILED"
  StatusPaymentRetry        SheetStatus = "PAYMENT_RETRY"
  StatusPaymentRetryApply   SheetStatus = "PAYMENT_RETRY_APPLY"
  StatusCancelled           SheetStatus = "CANCELLED"
)

type SheetItem struct {
  ItemSku     string   `json:"item_sku"`
  ItemName    string   `json:"item_name"`
  Price       float64  `json:"price"`
  Type        ItemType `gorm:"serializer:json" json:"type"`
  Multiple    string   `json:"multiple"`
  Reward      float64  `json:"reward"`
  Prize       bool     `json:"prize"`
  Code        string   `json:"code"`
  CodeValue   string   `json:"code_value"`
  IsFeatured  bool     `json:"is_featured"`
}

// Totals สรุปยอดรวมและส่วนลด
type Totals struct {
  Total      float64  `json:"total"`
  Discount   Discount `json:"discount" gorm:"embedded;embeddedPrefix:discount_"`
  GrandTotal float64  `json:"grand_total"`
  Reward     float64  `json:"reward"`
}

type Discount struct {
  Code  string  `json:"code"`
  Value float64 `json:"value"`
}

// TotoroSheet โครงสร้างหลัก (The Master Model)
type TotoroSheet struct {
  gorm.Model
  ID           uint   `gorm:"primaryKey" json:"-"`
  SheetID      string `gorm:"uniqueIndex;not null" json:"sheet_id"`
  SheetCodeDD  string `json:"sheet_code_dd"`
  SheetCodeTTT string `json:"sheet_code_ttt"`

  // Data Sections
  Items   []SheetItem `json:"items" gorm:"serializer:json"`
  Totals  Totals      `json:"totals" gorm:"serializer:json"`

  // Ownership
  UserID    string `gorm:"index" json:"user_id"`
  UserName  string `json:"user_name"`
  AgentID   string `gorm:"index" json:"agent_id"`
  AgentName string `json:"agent_name"`

  // Status & Workflow
  Status SheetStatus `gorm:"default:'DRAFT'" json:"status"`

  // Timestamps (ISO8601 รองรับใน JSON อัตโนมัติ)
  ExpiresAt          *time.Time    `json:"expires_at"`
  PaidAt             *time.Time    `json:"paid_at"`
  RewardAt           time.Time     `json:"reward_at"`
  RewardCompletedAt  *time.Time    `json:"reward_completed_at"`
  PaymentCompletedAt *time.Time    `json:"payment_completed_at"`
  CompletedAt        *time.Time    `json:"completed_at"`
}

// BeforeSave จะทำงานอัตโนมัติก่อนบันทึกลง DB
func (s *TotoroSheet) BeforeSave(tx *gorm.DB) (err error) {
  var total float64
  var totalReward float64

  // วนลูปคำนวณจากรายการ Items
  for _, item := range s.Items {
    total += item.Price

    // ถ้า Prize เป็น true ให้บวก Reward เข้ายอดรวมรางวัล
    if item.Prize {
      totalReward += item.Reward
    }
  }

  // อัปเดตค่าใน Totals Object
  s.Totals.Total = total
  s.Totals.Reward = totalReward
  s.Totals.GrandTotal = total - s.Totals.Discount.Value

  return nil
}
