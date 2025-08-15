package controllers

import (
	"collp-backend/repositories"
	"collp-backend/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

var userService services.UserService

// InitUserController initialize user service
func InitUserController(db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	userService = services.NewUserService(userRepo)
}

// GetUserByID ดึงข้อมูล user ตาม ID
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL parameter
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Get user from service
	user, err := userService.GetUserByID(uint(id))
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	// Return user data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// GetAllUsers ดึงรายการ users ทั้งหมดแบบ pagination
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Get users from service
	users, total, err := userService.GetAllUsers(page, limit)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Calculate pagination info
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// Return paginated data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"users":       users,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	})
}

// UpdateUserProfile อัพเดท profile ของ user
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSONError(w, http.StatusMethodNotAllowed, "Only PUT method allowed")
		return
	}

	// Parse user ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Parse request body
	var reqBody struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Update user profile
	if err := userService.UpdateUserProfile(uint(id), reqBody.Name, reqBody.Avatar); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User profile updated successfully",
	})
}

// DeactivateUser ปิดการใช้งาน user
func DeactivateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeJSONError(w, http.StatusMethodNotAllowed, "Only PATCH method allowed")
		return
	}

	// Parse user ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Deactivate user
	if err := userService.DeactivateUser(uint(id)); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User deactivated successfully",
	})
}

// ActivateUser เปิดการใช้งาน user
func ActivateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeJSONError(w, http.StatusMethodNotAllowed, "Only PATCH method allowed")
		return
	}

	// Parse user ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Activate user
	if err := userService.ActivateUser(uint(id)); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User activated successfully",
	})
}

// DeleteUser ลบ user (soft delete)
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSONError(w, http.StatusMethodNotAllowed, "Only DELETE method allowed")
		return
	}

	// Parse user ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Delete user
	if err := userService.DeleteUser(uint(id)); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User deleted successfully",
	})
}

// SearchUsers ค้นหา users
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	// Parse search parameters
	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		writeJSONError(w, http.StatusBadRequest, "Search keyword is required")
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Search users
	users, total, err := userService.SearchUsers(keyword, page, limit)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Calculate pagination info
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// Return search results
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"users":       users,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
			"keyword":     keyword,
		},
	})
}

// GetUserStats ดึงสถิติของ users
func GetUserStats(w http.ResponseWriter, r *http.Request) {
	// Get stats from service
	stats, err := userService.GetUserStats()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return stats
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

// CollPLogin (Legacy - keep for backward compatibility)
func CollPLogin(w http.ResponseWriter, r *http.Request) {
	// For now, return a placeholder response
	// TODO: Implement proper login logic with userService
	response := map[string]interface{}{
		"success": false,
		"message": "Login endpoint is being refactored. Please use OAuth login instead.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CollPRegister (Legacy - keep for backward compatibility)
func CollPRegister(w http.ResponseWriter, r *http.Request) {
	// Get headers for logging purposes
	username := r.Header.Get("username")
	password := r.Header.Get("password")
	phone := r.Header.Get("phone")
	address := r.Header.Get("address")

	// Log registration attempt
	log.Printf("Legacy registration attempt - Username: %s, Phone: %s, Address: %s", username, phone, address)
	if password != "" {
		log.Printf("Password provided: yes")
	}

	// For now, return a placeholder response
	// TODO: Implement proper registration logic with userService
	response := map[string]interface{}{
		"success": false,
		"message": "Registration endpoint is being refactored. Please use OAuth registration instead.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
