package models

import (
	"gorm.io/gorm"
)

// User model tương ứng với bảng `users`
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Email    string `json:"email" gorm:"unique;not null"`
	Name     string `json:"name" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	// CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	// UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	// DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Tùy chọn: cho soft delete
}

// BeforeCreate hook — có thể dùng để hash password hoặc validate
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Ví dụ: validate email hoặc hash password
	return nil
}
