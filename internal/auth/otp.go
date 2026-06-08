// internal/auth/otp.go
package auth

import (
	"fmt"
	"math/rand"
	"time"
	"totoro-app/internal/core"
)

const (
  OTP_PREFIX    = "otp:"           // Key สำหรับเก็บรหัส OTP
  VERIFY_PREFIX = "verified:"      // Key สำหรับยืนยันว่าเบอร์นี้ผ่านแล้ว
  OTP_TTL       = 5 * time.Minute  // รหัสมีอายุ 5 นาที
  VERIFY_TTL    = 10 * time.Minute // ยืนยันผ่านแล้ว มีเวลาสมัครสมาชิก 10 นาที
)

// GenerateAndSaveOTP สร้างรหัส 6 หลักและเก็บลง Redis
func GenerateAndSaveOTP(mobile string) (string, error) {
  // สุ่มเลข 6 หลัก
  otp := fmt.Sprintf("%06d", rand.Intn(1000000))

  // เก็บลง Redis (otp:0812345678 -> 123456)
  key := OTP_PREFIX + mobile
  err := core.RDB.Set(core.Ctx, key, otp, OTP_TTL).Err()

  return otp, err
}

// VerifyOTPCheck ตรวจสอบว่ารหัสที่กรอกมาตรงกับใน Redis ไหม
func VerifyOTPCheck(mobile, inputOTP string) (bool, error) {
  key := OTP_PREFIX + mobile
  storedOTP, err := core.RDB.Get(core.Ctx, key).Result()
  if err != nil {
    return false, err // อาจจะเพราะ Expired หรือไม่มี Key
  }

  if storedOTP == inputOTP {
    // ถ้าตรง -> ลบ OTP ทิ้ง และสร้าง "ตั๋วผ่านทาง" ชั่วคราว
    core.RDB.Del(core.Ctx, key)

    verifyKey := VERIFY_PREFIX + mobile
    core.RDB.Set(core.Ctx, verifyKey, "true", VERIFY_TTL)
    return true, nil
  }

  return false, nil
}

// IsMobileVerified ตรวจสอบว่าเบอร์นี้ผ่านการ Verify มาหรือยัง (ใช้ใน HandleRegister)
func IsMobileVerified(mobile string) bool {
  key := VERIFY_PREFIX + mobile
  exists, _ := core.RDB.Exists(core.Ctx, key).Result()
  return exists > 0
}
