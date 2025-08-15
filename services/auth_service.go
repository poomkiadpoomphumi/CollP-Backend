package services

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	googleOauthConfig *oauth2.Config
	privateKey        *rsa.PrivateKey
}

type GoogleUserInfo struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
	Token         string `json:"token"`
	TokenExpiry   int64  `json:"token_expiry"`
}

type AuthServiceInterface interface {
	InitGoogleOauth() error
	GetGoogleAuthURL(state string) string
	HandleGoogleCallback(code, state string) (*GoogleUserInfo, error)
}

func NewAuthService() AuthServiceInterface {
	return &AuthService{}
}

// InitGoogleOauth โหลด Google OAuth config และ RSA Key
func (s *AuthService) InitGoogleOauth() error {
	// Setup Google OAuth Config
	s.googleOauthConfig = &oauth2.Config{
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
		return fmt.Errorf("failed to read rsa.pem: %v", err)
	}

	// แปลง PEM เป็น RSA Private Key
	s.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return fmt.Errorf("failed to parse RSA private key: %v", err)
	}

	return nil
}

// GetGoogleAuthURL สร้าง URL สำหรับ redirect ไป Google
func (s *AuthService) GetGoogleAuthURL(state string) string {
	return s.googleOauthConfig.AuthCodeURL(state)
}

// HandleGoogleCallback จัดการ callback จาก Google และสร้าง JWT
func (s *AuthService) HandleGoogleCallback(code, state string) (*GoogleUserInfo, error) {
	// Validate state
	if state != "random-state-string" {
		return nil, errors.New("state parameter doesn't match")
	}

	// Exchange code for token
	token, err := s.googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) {
			return nil, fmt.Errorf("network error: %v", err)
		}
		return nil, fmt.Errorf("token exchange failed: %v", err)
	}

	// Get user info from Google
	userInfo, err := s.fetchGoogleUserInfo(token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %v", err)
	}

	// Generate JWT token
	jwtToken, expiry, err := s.generateJWTToken(userInfo.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %v", err)
	}

	// Add token to user info
	userInfo.Token = jwtToken
	userInfo.TokenExpiry = expiry

	return userInfo, nil
}

// fetchGoogleUserInfo ดึงข้อมูล user จาก Google API
func (s *AuthService) fetchGoogleUserInfo(token *oauth2.Token) (*GoogleUserInfo, error) {
	client := s.googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get(os.Getenv("GOOGLE_USERINFO"))
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %v", err)
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}

	return &userInfo, nil
}

// generateJWTToken สร้าง JWT token
func (s *AuthService) generateJWTToken(email string) (string, int64, error) {
	expTime := time.Now().Add(2 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"email": email,
		"exp":   expTime,
		"iat":   time.Now().Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := jwtToken.SignedString(s.privateKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expTime, nil
}
