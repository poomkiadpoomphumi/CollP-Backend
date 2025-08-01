package controller

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config
var privateKey *rsa.PrivateKey // ใช้เก็บ RSA Key
// InitGoogleOauth โหลด Google OAuth config และ RSA Key
func InitGoogleOauth() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			os.Getenv("GOOGLE_USERINFO_EMAIL"),
			os.Getenv("GOOGLE_USERINFO_PROFILE"),
		},
		Endpoint: google.Endpoint,
	}

	// โหลด RSA Private Key จากไฟล์
	keyData, err := ioutil.ReadFile("rsa.pem")
	if err != nil {
		log.Fatalf("Failed to read rsa.pem: %v", err)
	}

	// แปลง PEM เป็น RSA Private Key
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		log.Fatalf("Failed to parse RSA private key: %v", err)
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// GoogleLogin redirect ผู้ใช้ไปหน้า Google Login
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL("random-state-string")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback รับ code จาก Google แล้วแลก token + ดึง user info
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != "random-state-string" {
		writeJSONError(w, http.StatusBadRequest, "State parameter doesn't match")
		return
	}
	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		var netErr net.Error
		if ok := errors.As(err, &netErr); ok {
			writeJSONError(w, http.StatusGatewayTimeout, fmt.Sprintf("Network error: %v", err))
			return
		}
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Token exchange failed: %v", err))
		return
	}
	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get(os.Getenv("GOOGLE_USERINFO"))
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed getting user info: %v", err))
		return
	}
	defer resp.Body.Close()
	userInfo := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to parse user info: %v", err))
		return
	}
	expTime := time.Now().Add(2 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"email": userInfo["email"],
		"exp":   expTime,
		"iat":   time.Now().Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := jwtToken.SignedString(privateKey)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate JWT: %v", err))
		return
	}
	// เพิ่ม token และ expiry เข้าไปใน userInfo map
	userInfo["token"] = tokenString
	userInfo["token_expiry"] = expTime

	// send json data
	/* 	w.Header().Set("Content-Type", "application/json")
   	w.WriteHeader(http.StatusOK)
   	json.NewEncoder(w).Encode(userInfo) */

	// แปลงข้อมูล userInfo map เป็น query string
	values := url.Values{}
	for k, v := range userInfo {
		// แปลงค่าเป็น string ให้หมดก่อน แล้ว url.QueryEscape จะทำให้ส่ง query string ได้ปลอดภัย
		strVal := fmt.Sprintf("%v", v)
		// ถ้า value เป็น string แล้วมี & หรือ ตัวอักษรพิเศษก็ถูก escape อยู่แล้ว
		values.Set(k, strVal)
	}
	// กำหนด URL frontend ที่ต้องการ redirect ไป พร้อม query string ทั้งหมด
	frontendRedirectURL := fmt.Sprintf("%s?%s", os.Getenv("FRONTEND_REDIRECT"), values.Encode())
	// redirect ไป frontend พร้อมส่ง token + userInfo ผ่าน query string
	http.Redirect(w, r, frontendRedirectURL, http.StatusSeeOther)
}


