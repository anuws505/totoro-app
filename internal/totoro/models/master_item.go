// internal/totoro/models/master_item.go
package models

import (
	"gorm.io/gorm"
)

const (
  MultipleDirect = "DIRECT"
  MultipleMixed  = "MIXED"
)

type ItemType struct {
  Direct float64 `json:"DIRECT"`
  Mixed  float64 `json:"MIXED"`
}

// MasterItem สำหรับ Database
type MasterItem struct {
  gorm.Model
  ItemSku    string   `gorm:"uniqueIndex;type:varchar(50);column:item_sku"`
  ItemName   string   `gorm:"column:item_name"`
  Price      float64  `gorm:"column:price;default:0"`
  Type       ItemType `gorm:"serializer:json;column:type"`
  Code       string   `gorm:"column:code"`
  CodeValue  string   `gorm:"column:code_value"`
  IsFeatured bool     `gorm:"column:is_featured;default:false"`
  // หมายเหตุ: Multiple, Reward, Prize มักจะเกิดขึ้นในจังหวะคำนวณของ Sheet
  // ถ้าใน Master ไม่มีค่าคงที่สำหรับพวกนี้ สามารถตัดออกเพื่อลดความซ้ำซ้อนได้ครับ
}

// MasterItemResponse สำหรับส่งไปให้ Frontend
type MasterItemResponse struct {
  ItemSku    string   `json:"item_sku"`
  ItemName   string   `json:"item_name"`
  Price      float64  `json:"price"`
  Type       ItemType `json:"type"`
  Code       string   `json:"code"`
  IsFeatured bool     `json:"is_featured"`
}

// ToResponse: Helper สำหรับแปลงจาก Model เป็น Response (แบบตัวเดียว)
func (m *MasterItem) ToResponse() MasterItemResponse {
  return MasterItemResponse{
    ItemSku:    m.ItemSku,
    ItemName:   m.ItemName,
    Price:      m.Price,
    Type:       m.Type,
    Code:       m.Code,
    IsFeatured: m.IsFeatured,
  }
}

// โครงสร้างสำหรับรับข้อมูลขาเข้าจาก API
type CreateOrPatchItemInput struct {
  ItemSku    string   `json:"item_sku"`
  ItemName   string   `json:"item_name"`
  Price      float64  `json:"price"`
  Type       ItemType `json:"type"`
  Code       string   `json:"code"`
  CodeValue  string   `json:"code_value"`
  IsFeatured bool     `json:"is_featured"`
}
