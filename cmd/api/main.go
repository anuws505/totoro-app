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
  core.InitDB()

  fmt.Println("Migrating database...")
  core.DB.AutoMigrate(
    &models.TotoroSheet{},
    &models.MasterItem{},
    &models.PromotionCode{},
    &models.MasterResult{},
    &auth.User{},
  )
  fmt.Println("Database Migrated: All tables are ready!")

  r := mux.NewRouter()

  r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
  }).Methods("GET")

  // ----------------------------------------------------
  // 1. Public Routes (ไม่ต้องผูก Middleware)
  // ----------------------------------------------------
  publicRouter := r.PathPrefix("/api/v1").Subrouter()

  publicRouter.HandleFunc("/auth/create", auth.HandleCreate).Methods("POST")
  publicRouter.HandleFunc("/auth/login", auth.HandleLogin).Methods("POST")
  publicRouter.HandleFunc("/auth/otp/request", auth.HandleRequestOTP).Methods("POST")
  publicRouter.HandleFunc("/auth/otp/verify", auth.HandleVerifyOTP).Methods("POST")
  publicRouter.HandleFunc("/items", services.HandleGetMasterItems).Methods("GET")

  // ----------------------------------------------------
  // 2. Protected Routes (ต้องผ่าน AuthMiddleware)
  // ----------------------------------------------------
  protectedRouter := r.PathPrefix("/api/v1").Subrouter()
  protectedRouter.Use(auth.AuthMiddleware)

  // Logout & User Profile
  protectedRouter.HandleFunc("/auth/logout", auth.HandleLogout).Methods("POST")
  protectedRouter.HandleFunc("/user/profile", auth.HandleUpdateUsername).Methods("PATCH")
  protectedRouter.HandleFunc("/user/change-password", auth.HandleChangePassword).Methods("POST")

  // Sheet Service
  protectedRouter.HandleFunc("/sheets", services.HandleCreateSheet).Methods("POST")

  // ----------------------------------------------------
  // 3. Protected Admin Routes (ผ่าน AuthMiddleware + AdminMiddleware)
  // ----------------------------------------------------
  adminRouter := r.PathPrefix("/api/v1/admin").Subrouter()
  adminRouter.Use(auth.AuthMiddleware)
  adminRouter.Use(auth.AdminMiddleware)

  // User Management
  adminRouter.HandleFunc("/update-role", auth.HandleUpdateUserRole).Methods("PATCH")

  // Master Result Management
  adminRouter.HandleFunc("/result-master", services.HandleCreateMasterResult).Methods("POST")
  adminRouter.HandleFunc("/result-master/{id}", services.HandlePatchMasterResult).Methods("PATCH")

  // Master Item Management
  adminRouter.HandleFunc("/item-master", services.HandleCreateMasterItem).Methods("POST")
  adminRouter.HandleFunc("/item-master", services.HandleUpdateMasterItem).Methods("PATCH")
  adminRouter.HandleFunc("/item-master", services.HandleDeleteMasterItem).Methods("DELETE")
  adminRouter.HandleFunc("/item-master/restore", services.HandleRestoreMasterItem).Methods("POST")

  // Promotion Management
  adminRouter.HandleFunc("/promotions", services.HandleGetPromotions).Methods("GET")
  adminRouter.HandleFunc("/promotions", services.HandleCreatePromotion).Methods("POST")
  adminRouter.HandleFunc("/promotions", services.HandleUpdatePromotion).Methods("PATCH")
  adminRouter.HandleFunc("/promotions", services.HandleDeletePromotion).Methods("DELETE")

  // ----------------------------------------------------
  // Background Workers
  // ----------------------------------------------------
  go workers.DoCleanupExpiredDrafts()
  go workers.DoCalculatePrizes()


  // Start Server
  port := ":8080"
  fmt.Printf("Totoro App is running on http://localhost%s\n", port)

  if err := http.ListenAndServe(port, r); err != nil {
    fmt.Printf("Server failed: %s\n", err)
  }
}
