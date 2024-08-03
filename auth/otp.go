package auth

import (
    "net/http"
    "time"
    "github.com/citra-org/cypo/config"
)

func HandleOTP(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    otpInput := r.FormValue("otp")

    config.Mu.Lock()
    if otpInput != config.OTP || time.Since(config.OtpIssued) > config.OTPExpiration {
        config.Mu.Unlock()
        http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
        return
    }
    if config.Authenticated {
        config.Mu.Unlock()
        http.Error(w, "Already authenticated", http.StatusConflict)
        return
    }
    config.OtpIssued = time.Now()
    config.Authenticated = true
    config.Mu.Unlock()

    http.ServeFile(w, r, "assets/upload.html")
}