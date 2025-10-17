package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"myapp/database"
	"myapp/models"
)

// setupTestDB tạo mock database cho testing
func setupTestDB(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return mock, gormDB
}

func TestCreateUser_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mock, gormDB := setupTestDB(t)
	database.DB = gormDB

	// Mock SQL expectations
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs("John Doe", "john@example.com", "1234567890").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Create request
	user := models.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "1234567890",
	}
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	CreateUser(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", response.Name)
	assert.Equal(t, "john@example.com", response.Email)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	_, gormDB := setupTestDB(t)
	database.DB = gormDB

	// Create invalid request
	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	CreateUser(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "invalid")
}

func TestCreateUser_DatabaseError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mock, gormDB := setupTestDB(t)
	database.DB = gormDB

	// Mock SQL expectations với lỗi
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WillReturnError(gorm.ErrInvalidData)
	mock.ExpectRollback()

	// Create request
	user := models.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password",
	}
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	CreateUser(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["error"])
}

func TestGetUsers_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mock, gormDB := setupTestDB(t)
	database.DB = gormDB

	// Mock SQL expectations
	rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
		AddRow(1, "John Doe", "john@example.com", "1234567890").
		AddRow(2, "Jane Doe", "jane@example.com", "0987654321")

	mock.ExpectQuery("SELECT \\* FROM `users`").
		WillReturnRows(rows)

	// Create request
	req, _ := http.NewRequest("GET", "/users", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	GetUsers(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "John Doe", response[0].Name)
	assert.Equal(t, "Jane Doe", response[1].Name)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUsers_DatabaseError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mock, gormDB := setupTestDB(t)
	database.DB = gormDB

	// Mock SQL expectations với lỗi
	mock.ExpectQuery("SELECT \\* FROM `users`").
		WillReturnError(gorm.ErrInvalidDB)

	// Create request
	req, _ := http.NewRequest("GET", "/users", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	GetUsers(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["error"])
}

func TestGetUsers_EmptyResult(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mock, gormDB := setupTestDB(t)
	database.DB = gormDB

	// Mock SQL expectations với kết quả rỗng
	rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"})
	mock.ExpectQuery("SELECT \\* FROM `users`").
		WillReturnRows(rows)

	// Create request
	req, _ := http.NewRequest("GET", "/users", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	GetUsers(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Empty(t, response)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
