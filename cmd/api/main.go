// cmd/api/main.go
package main

import (
	"fmt"
	"net/http"
	"totoro-app/internal/auth"
	"totoro-app/internal/core"
	"totoro-app/internal/totoro/models"
	"totoro-app/internal/totoro/services"
	"totoro-app/internal/totoro/workers"

	"github.com/gorilla/mux"
)

func main() {
  // เชื่อมต่อ DB
  core.InitDB()

  // ทำ AutoMigrate ที่นี่ที่เดียว
  fmt.Println("Migrating database...")
  core.DB.AutoMigrate(
    &models.TotoroSheet{},
    &models.MasterItem{},
    &models.PromotionCode{},
    &models.MasterResult{},
    &auth.User{},
  )
  fmt.Println("Database Migrated: All tables are ready!")

  // กำหนดเส้นทาง (Routes) ของ API
  r := mux.NewRouter()

  // --- Public Routes (ไม่ต้องใช้ Middleware) ---
  r.HandleFunc("/api/v1/auth/login", auth.HandleLogin).Methods("POST")
  r.HandleFunc("/api/v1/auth/register", auth.HandleRegister).Methods("POST")
  r.HandleFunc("/api/v1/auth/otp/request", auth.HandleRequestOTP).Methods("POST")
  r.HandleFunc("/api/v1/auth/otp/verify", auth.HandleVerifyOTP).Methods("POST")

  // --- Protected Routes (ต้องผ่าน AuthMiddleware) ---
  r.HandleFunc("/api/v1/auth/logout",
    auth.AuthMiddleware(auth.HandleLogout),
  ).Methods("POST")

  r.HandleFunc("/api/v1/user/profile",
    auth.AuthMiddleware(auth.HandleUpdateUsername),
  ).Methods("PATCH")

  r.HandleFunc("/api/v1/user/change-password",
    auth.AuthMiddleware(auth.HandleChangePassword),
  ).Methods("POST")

  // Admin update user role
  r.HandleFunc("/api/v1/admin/update-role",
    auth.AuthMiddleware(auth.AdminMiddleware(auth.HandleUpdateUserRole)),
  ).Methods("PATCH")

  // MasterResult
  r.HandleFunc("/api/v1/admin/result-master",
    auth.AuthMiddleware(auth.AdminMiddleware(services.HandleCreateMasterResult)),
  ).Methods("POST")

  r.HandleFunc("/api/v1/admin/result-master/{id}",
    auth.AuthMiddleware(auth.AdminMiddleware(services.HandlePatchMasterResult)),
  ).Methods("PATCH")

  // MasterItem
  r.HandleFunc("/api/v1/items", services.HandleGetMasterItems).Methods("GET")

  r.HandleFunc("/api/v1/admin/item-master",
    auth.AuthMiddleware(auth.AdminMiddleware(services.HandleCreateMasterItem)),
  ).Methods("POST")

  r.HandleFunc("/api/v1/admin/item-master",
    auth.AuthMiddleware(auth.AdminMiddleware(services.HandleUpdateMasterItem)),
  ).Methods("PATCH")

  r.HandleFunc("/api/v1/admin/item-master",
    auth.AuthMiddleware(auth.AdminMiddleware(services.HandleDeleteMasterItem)),
  ).Methods("DELETE")

  r.HandleFunc("/api/v1/admin/item-master/restore",
    auth.AuthMiddleware(auth.AdminMiddleware(services.HandleRestoreMasterItem)),
  ).Methods("POST")

  // Promotion
  r.HandleFunc("/api/v1/admin/promotions",
    auth.AuthMiddleware(auth.AdminMiddleware(services.HandleGetPromotions)),
  ).Methods("GET")

  r.HandleFunc("/api/v1/admin/promotions",
    auth.AuthMiddleware(auth.AdminMiddleware(services.HandleCreatePromotion)),
  ).Methods("POST")

  r.HandleFunc("/api/v1/admin/promotions",
      auth.AuthMiddleware(auth.AdminMiddleware(services.HandleUpdatePromotion)),
  ).Methods("PATCH")

  r.HandleFunc("/api/v1/admin/promotions",
    auth.AuthMiddleware(auth.AdminMiddleware(services.HandleDeletePromotion)),
  ).Methods("DELETE")

  // SheetItem
  r.HandleFunc("/api/v1/sheets",
    auth.AuthMiddleware(services.HandleCreateSheet),
  ).Methods("POST")

  // Workers
  go workers.DoCleanupExpiredDrafts()
  go workers.DoCalculatePrizes()

  // เริ่มเปิด Server ที่ Port 8080
  port := ":8080"
  fmt.Printf("Totoro App is running on http://localhost%s\n", port)

  if err := http.ListenAndServe(port, r); err != nil {
    fmt.Printf("Server failed: %s\n", err)
  }
}
