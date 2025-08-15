package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"collp-backend/services"
)

var authService services.AuthServiceInterface

// InitAuthController initialize auth service
func InitAuthController() {
	authService = services.NewAuthService()
	if err := authService.InitGoogleOauth(); err != nil {
		log.Fatalf("Failed to initialize auth service: %v", err)
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// GoogleLogin redirect ผู้ใช้ไปหน้า Google Login
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// ใช้ service เพื่อสร้าง auth URL
	url := authService.GetGoogleAuthURL("random-state-string")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback รับ code จาก Google แล้วแลก token + ดึง user info
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// รับ parameters จาก request
	state := r.FormValue("state")
	code := r.FormValue("code")

	// เรียกใช้ service เพื่อ handle callback
	userInfo, err := authService.HandleGoogleCallback(code, state)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Convert userInfo to map สำหรับ query string
	values := url.Values{}
	values.Set("email", userInfo.Email)
	values.Set("name", userInfo.Name)
	values.Set("picture", userInfo.Picture)
	values.Set("verified_email", fmt.Sprintf("%v", userInfo.VerifiedEmail))
	values.Set("token", userInfo.Token)
	values.Set("token_expiry", fmt.Sprintf("%d", userInfo.TokenExpiry))

	// Redirect ไป frontend พร้อม query string
	frontendRedirectURL := fmt.Sprintf("%s?%s", os.Getenv("FRONTEND_REDIRECT"), values.Encode())
	http.Redirect(w, r, frontendRedirectURL, http.StatusSeeOther)
}
