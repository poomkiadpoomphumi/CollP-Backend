package services

import (
	"fmt"
	"strings"

	"collp-backend/models"
	"collp-backend/repositories"
)

// UserService interface สำหรับ user business logic
type UserService interface {
	// User management
	CreateUser(email, name, googleID, avatar string) (*models.User, error)
	GetOrCreateUser(email, name, googleID, avatar string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUserProfile(id uint, name, avatar string) error
	DeactivateUser(id uint) error
	ActivateUser(id uint) error
	DeleteUser(id uint) error

	// User queries
	GetAllUsers(page, limit int) ([]*models.User, int64, error)
	GetActiveUsers() ([]*models.User, error)
	SearchUsers(keyword string, page, limit int) ([]*models.User, int64, error)

	// Statistics
	GetUserStats() (*UserStats, error)

	// Validation
	IsValidEmail(email string) bool
	IsUserActive(id uint) (bool, error)
}

// UserStats สถิติของ users
type UserStats struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	InactiveUsers int64 `json:"inactive_users"`
}

// userService struct implements UserService interface
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates new user service instance
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser สร้าง user ใหม่
func (s *userService) CreateUser(email, name, googleID, avatar string) (*models.User, error) {
	// Validate email
	if !s.IsValidEmail(email) {
		return nil, fmt.Errorf("invalid email format: %s", email)
	}

	// Validate required fields
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name is required")
	}

	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Create user
	user := &models.User{
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Name:     strings.TrimSpace(name),
		GoogleID: googleID,
		Avatar:   avatar,
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetOrCreateUser ดึง user หรือสร้างใหม่ถ้าไม่มี (สำหรับ OAuth)
func (s *userService) GetOrCreateUser(email, name, googleID, avatar string) (*models.User, error) {
	// Normalize email
	email = strings.ToLower(strings.TrimSpace(email))

	// Validate email
	if !s.IsValidEmail(email) {
		return nil, fmt.Errorf("invalid email format: %s", email)
	}

	// Create user data
	userData := &models.User{
		Email:    email,
		Name:     strings.TrimSpace(name),
		GoogleID: googleID,
		Avatar:   avatar,
		IsActive: true,
	}

	user, err := s.userRepo.CreateIfNotExists(email, googleID, userData)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create user: %w", err)
	}

	return user, nil
}

// GetUserByID ดึง user ด้วย ID
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid user id")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail ดึง user ด้วย email
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if !s.IsValidEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUserProfile อัพเดท profile ของ user
func (s *userService) UpdateUserProfile(id uint, name, avatar string) error {
	if id == 0 {
		return fmt.Errorf("invalid user id")
	}

	// Validate name
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("name is required")
	}

	// Check if user exists
	exists, err := s.userRepo.Exists(id)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Update user
	fields := map[string]interface{}{
		"name":   name,
		"avatar": avatar,
	}

	if err := s.userRepo.UpdateFields(id, fields); err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	return nil
}

// DeactivateUser ปิดการใช้งาน user
func (s *userService) DeactivateUser(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid user id")
	}

	if err := s.userRepo.UpdateStatus(id, false); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

// ActivateUser เปิดการใช้งาน user
func (s *userService) ActivateUser(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid user id")
	}

	if err := s.userRepo.UpdateStatus(id, true); err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	return nil
}

// DeleteUser ลบ user (soft delete)
func (s *userService) DeleteUser(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid user id")
	}

	// Check if user exists
	exists, err := s.userRepo.Exists(id)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user not found")
	}

	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// GetAllUsers ดึง users ทั้งหมดแบบ pagination
func (s *userService) GetAllUsers(page, limit int) ([]*models.User, int64, error) {
	// Validate pagination parameters
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	users, total, err := s.userRepo.GetAll(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}

// GetActiveUsers ดึง active users ทั้งหมด
func (s *userService) GetActiveUsers() ([]*models.User, error) {
	users, err := s.userRepo.GetActive()
	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}

	return users, nil
}

// SearchUsers ค้นหา users
func (s *userService) SearchUsers(keyword string, page, limit int) ([]*models.User, int64, error) {
	// Validate search keyword
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, 0, fmt.Errorf("search keyword is required")
	}

	// Validate pagination
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	users, total, err := s.userRepo.Search(keyword, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	return users, total, nil
}

// GetUserStats ดึงสถิติของ users
func (s *userService) GetUserStats() (*UserStats, error) {
	totalUsers, err := s.userRepo.Count()
	if err != nil {
		return nil, fmt.Errorf("failed to count total users: %w", err)
	}

	activeUsers, err := s.userRepo.CountActive()
	if err != nil {
		return nil, fmt.Errorf("failed to count active users: %w", err)
	}

	stats := &UserStats{
		TotalUsers:    totalUsers,
		ActiveUsers:   activeUsers,
		InactiveUsers: totalUsers - activeUsers,
	}

	return stats, nil
}

// IsValidEmail ตรวจสอบรูปแบบ email
func (s *userService) IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	// Basic email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".") && len(email) > 5
}

// IsUserActive ตรวจสอบว่า user active หรือไม่
func (s *userService) IsUserActive(id uint) (bool, error) {
	if id == 0 {
		return false, fmt.Errorf("invalid user id")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	return user.IsActive, nil
}
