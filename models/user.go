package models

import "gorm.io/gorm"

// Model User
type User struct {
	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}

// Nếu cần thêm method của User có thể viết ở đây
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// ví dụ: validate trước khi insert
	return nil
}
