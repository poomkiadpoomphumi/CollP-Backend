package repositories

import (
	"errors"
	"fmt"

	"collp-backend/models"

	"gorm.io/gorm"
)

// UserRepository interface สำหรับ User CRUD operations
type UserRepository interface {
	// Create operations
	Create(user *models.User) error
	CreateIfNotExists(email, googleID string, user *models.User) (*models.User, error)

	// Read operations
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByGoogleID(googleID string) (*models.User, error)
	GetAll(page, limit int) ([]*models.User, int64, error)
	GetActive() ([]*models.User, error)

	// Update operations
	Update(id uint, user *models.User) error
	UpdateFields(id uint, fields map[string]interface{}) error
	UpdateStatus(id uint, isActive bool) error
	UpdateAvatar(id uint, avatarURL string) error

	// Delete operations
	Delete(id uint) error     // Soft delete
	HardDelete(id uint) error // Hard delete
	Restore(id uint) error    // Restore soft deleted

	// Search operations
	Search(keyword string, page, limit int) ([]*models.User, int64, error)

	// Utility operations
	Exists(id uint) (bool, error)
	ExistsByEmail(email string) (bool, error)
	Count() (int64, error)
	CountActive() (int64, error)
}

// userRepository struct implements UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates new user repository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create สร้าง user ใหม่
func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// CreateIfNotExists สร้าง user ใหม่ถ้ายังไม่มี
func (r *userRepository) CreateIfNotExists(email, googleID string, user *models.User) (*models.User, error) {
	// ตรวจสอบว่ามี user อยู่แล้วหรือไม่
	existingUser := &models.User{}

	// หา user ด้วย email หรือ google_id
	err := r.db.Where("email = ? OR google_id = ?", email, googleID).First(existingUser).Error
	if err == nil {
		// User มีอยู่แล้ว return user เดิม
		return existingUser, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// User ไม่มี สร้างใหม่
	if err := r.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID หา user ด้วย ID
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	user := &models.User{}
	if err := r.db.First(user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetByEmail หา user ด้วย email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := r.db.Where("email = ?", email).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetByGoogleID หา user ด้วย Google ID
func (r *userRepository) GetByGoogleID(googleID string) (*models.User, error) {
	user := &models.User{}
	if err := r.db.Where("google_id = ?", googleID).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with google_id %s not found", googleID)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetAll ดึง users ทั้งหมดแบบ pagination
func (r *userRepository) GetAll(page, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	// Count total records
	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated results
	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}

// GetActive ดึง users ที่ active
func (r *userRepository) GetActive() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Where("is_active = ?", true).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}
	return users, nil
}

// Update อัพเดท user
func (r *userRepository) Update(id uint, user *models.User) error {
	if err := r.db.Where("id = ?", id).Updates(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// UpdateFields อัพเดท fields เฉพาะ
func (r *userRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	if err := r.db.Model(&models.User{}).Where("id = ?", id).Updates(fields).Error; err != nil {
		return fmt.Errorf("failed to update user fields: %w", err)
	}
	return nil
}

// UpdateStatus อัพเดทสถานะ active/inactive
func (r *userRepository) UpdateStatus(id uint, isActive bool) error {
	if err := r.db.Model(&models.User{}).Where("id = ?", id).Update("is_active", isActive).Error; err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	return nil
}

// UpdateAvatar อัพเดท avatar URL
func (r *userRepository) UpdateAvatar(id uint, avatarURL string) error {
	if err := r.db.Model(&models.User{}).Where("id = ?", id).Update("avatar", avatarURL).Error; err != nil {
		return fmt.Errorf("failed to update user avatar: %w", err)
	}
	return nil
}

// Delete soft delete user
func (r *userRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// HardDelete ลบ user ออกจากฐานข้อมูลถาวร
func (r *userRepository) HardDelete(id uint) error {
	if err := r.db.Unscoped().Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to hard delete user: %w", err)
	}
	return nil
}

// Restore คืนค่า soft deleted user
func (r *userRepository) Restore(id uint) error {
	if err := r.db.Unscoped().Model(&models.User{}).Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
		return fmt.Errorf("failed to restore user: %w", err)
	}
	return nil
}

// Search ค้นหา users ด้วย keyword
func (r *userRepository) Search(keyword string, page, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	query := r.db.Model(&models.User{}).Where(
		"name ILIKE ? OR email ILIKE ?",
		"%"+keyword+"%",
		"%"+keyword+"%",
	)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Get results with pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	return users, total, nil
}

// Exists ตรวจสอบว่า user มีอยู่หรือไม่
func (r *userRepository) Exists(id uint) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByEmail ตรวจสอบว่า email มีอยู่หรือไม่
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

// Count นับจำนวน users ทั้งหมด
func (r *userRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// CountActive นับจำนวน active users
func (r *userRepository) CountActive() (int64, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("is_active = ?", true).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count active users: %w", err)
	}
	return count, nil
}
