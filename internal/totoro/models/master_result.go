// internal/totoro/models/master_result.go
package models

import "time"

type MasterResult struct {
  ID         uint      `gorm:"primaryKey" json:"id"`
  DateReward time.Time `gorm:"index;type:timestamp" json:"date_reward"` // เวลาที่ออกรางวัล (เก็บ UTC แต่ตอนเช็คค่อยแปลง)
  IsActive   bool      `gorm:"default:true" json:"is_active"`

  // เก็บผลรางวัลแยกตามประเภท (Serializer JSON เพื่อความง่ายและยืดหยุ่น)
  // ตัวอย่างข้อมูล: {"DD": "10", "TTT": "160"}
  Results    ResultMap `gorm:"serializer:json" json:"results"`

  CreatedAt  time.Time `json:"created_at"`
  UpdatedAt  time.Time `json:"updated_at"`
}

type ResultMap map[string]string
