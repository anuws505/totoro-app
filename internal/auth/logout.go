// internal/auth/logout.go
package auth

import (
	"net/http"
	"totoro-app/internal/core"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
  userID := r.Context().Value("userID").(string)
  core.RDB.Del(core.Ctx, "active_token:"+userID)
  core.RDB.Del(core.Ctx, "user_role:"+userID)

  core.WriteSuccess(w, http.StatusOK, "Logged out successfully", "20000", nil)
}
