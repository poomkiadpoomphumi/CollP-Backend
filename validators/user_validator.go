package validators

import (
	"collp-backend/utils"
	"errors"
)

// UserRegistrationRequest represents user registration data
type UserRegistrationRequest struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone,omitempty"`
	Address  string `json:"address,omitempty"`
}

// UserLoginRequest represents user login data
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ValidateUserRegistration validates user registration data
func ValidateUserRegistration(req UserRegistrationRequest) error {
	if utils.IsEmpty(req.Email) {
		return errors.New("email is required")
	}

	if !utils.IsValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if utils.IsEmpty(req.Name) {
		return errors.New("name is required")
	}

	if utils.IsEmpty(req.Password) {
		return errors.New("password is required")
	}

	if !utils.IsValidPassword(req.Password) {
		return errors.New("password must be at least 8 characters long and contain uppercase, lowercase, and number")
	}

	return nil
}

// ValidateUserLogin validates user login data
func ValidateUserLogin(req UserLoginRequest) error {
	if utils.IsEmpty(req.Email) {
		return errors.New("email is required")
	}

	if !utils.IsValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if utils.IsEmpty(req.Password) {
		return errors.New("password is required")
	}

	return nil
}
