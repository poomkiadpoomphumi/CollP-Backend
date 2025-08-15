package utils

import (
	"regexp"
	"strings"
)

// IsValidEmail checks if email format is valid
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidPassword checks if password meets requirements
func IsValidPassword(password string) bool {
	// At least 8 characters, contains uppercase, lowercase, and number
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}

// SanitizeString removes leading/trailing whitespace and converts to lowercase
func SanitizeString(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}

// IsEmpty checks if string is empty after trimming whitespace
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}
