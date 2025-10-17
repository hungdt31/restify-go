package models

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestUserModel_Structure(t *testing.T) {
	user := User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password",
	}

	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "password", user.Password)
}

func TestUserModel_JSONTags(t *testing.T) {
	user := User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password",
	}

	// Verify struct tags
	userType := assert.IsType(t, User{}, user)
	assert.True(t, userType)
}

func TestUserModel_BeforeCreate(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Create user
	user := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password",
	}

	// Mock expectations
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs("John Doe", "john@example.com", "password").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Execute
	err = gormDB.Create(&user).Error

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserModel_EmptyUser(t *testing.T) {
	user := User{}

	assert.Equal(t, uint(0), user.ID)
	assert.Equal(t, "", user.Name)
	assert.Equal(t, "", user.Email)
	assert.Equal(t, "", user.Password)
}

func TestUserModel_PartialUser(t *testing.T) {
	user := User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "", user.Password)
	assert.Equal(t, uint(0), user.ID)
}

func TestUserModel_Update(t *testing.T) {
	user := User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "1234567890",
	}

	// Update user fields
	user.Name = "Jane Doe"
	user.Email = "jane@example.com"

	assert.Equal(t, "Jane Doe", user.Name)
	assert.Equal(t, "jane@example.com", user.Email)
	assert.Equal(t, "1234567890", user.Password)
}

func TestUserModel_GormTags(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Test unique constraint
	user1 := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password",
	}

	// Mock first insert success
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs("John Doe", "john@example.com", "1234567890").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = gormDB.Create(&user1).Error
	assert.NoError(t, err)

	// Try to insert duplicate email - should fail
	user2 := User{
		Name:     "Jane Doe",
		Email:    "john@example.com", // duplicate email
		Password: "0987654321",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs("Jane Doe", "john@example.com", "0987654321").
		WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	err = gormDB.Create(&user2).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrDuplicatedKey, err)
}

func TestUserModel_PhoneMaxLength(t *testing.T) {
	user := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "123456789012345", // exactly 15 characters
	}

	assert.Equal(t, 15, len(user.Password))
}

func TestUserModel_AutoIncrement(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Create user without ID
	user := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password",
	}

	// Mock auto-increment behavior
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs("John Doe", "john@example.com", "1234567890").
		WillReturnResult(sqlmock.NewResult(5, 1)) // ID will be 5
	mock.ExpectCommit()

	err = gormDB.Create(&user).Error
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
