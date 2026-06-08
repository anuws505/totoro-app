// internal/auth/otp_handle.go
package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"totoro-app/internal/core"
)

type OTPRequest struct {
  Mobile string `json:"mobile"`
}

type OTPVerifyRequest struct {
  Mobile string `json:"mobile"`
  OTP    string `json:"otp"`
}

// HandleRequestOTP
func HandleRequestOTP(w http.ResponseWriter, r *http.Request) {
  var req OTPRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    core.WriteError(w, http.StatusBadRequest, "Invalid request", "40000")
    return
  }

  otp, err := GenerateAndSaveOTP(req.Mobile)
  if err != nil {
    core.WriteError(w, http.StatusInternalServerError,
      "Internal server error: Redis failure", "50003")
    return
  }

  // Mock SMS sending by printing to console
  fmt.Printf(" [SMS] OTP for mobile %s is: %s\n", req.Mobile, otp)

  core.WriteSuccess(w, http.StatusOK, "OTP has been sent successfully", "20000", nil)
}

// HandleVerifyOTP
func HandleVerifyOTP(w http.ResponseWriter, r *http.Request) {
  var req OTPVerifyRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    core.WriteError(w, http.StatusBadRequest, "Invalid request", "40000")
    return
  }

  ok, err := VerifyOTPCheck(req.Mobile, req.OTP)
  if err != nil || !ok {
    core.WriteError(w, http.StatusBadRequest, "Invalid or expired OTP code", "40005")
    return
  }

  // Return response
  core.WriteSuccess(w, http.StatusOK, "Mobile number verified successfully", "20000", nil)
}
