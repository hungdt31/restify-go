package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"myapp/database"

	"gorm.io/gorm"
)

// TestMain setup cho tất cả tests trong package
func TestMain(m *testing.M) {
	// Set JWT secret trước khi chạy bất kỳ test nào
	os.Setenv("JWT_SECRET", "test_secret")

	// Chạy tests
	code := m.Run()

	// Cleanup
	os.Unsetenv("JWT_SECRET")
	os.Exit(code)
}

func TestLogin_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mock, gormDB := setupTestDB(t)
	database.DB = gormDB

	// Mock SQL expectations - tìm user theo email
	rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
		AddRow(1, "John Doe", "john@example.com", "1234567890")

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs("john@example.com", 1).
		WillReturnRows(rows)

	// Create request
	loginData := map[string]string{
		"email": "john@example.com",
	}
	body, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	Login(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["token"])

	// Verify token is valid
	token, err := jwt.Parse(response["token"], func(token *jwt.Token) (interface{}, error) {
		return []byte("test_secret"), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	// Verify claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		assert.Equal(t, float64(1), claims["user_id"])
		assert.NotNil(t, claims["exp"])
	}

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_InvalidJSON(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	_, gormDB := setupTestDB(t)
	database.DB = gormDB

	// Create invalid request
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	Login(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "invalid")
}

func TestLogin_UserNotFound(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mock, gormDB := setupTestDB(t)
	database.DB = gormDB

	// Mock SQL expectations - user không tồn tại
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs("notfound@example.com", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Create request
	loginData := map[string]string{
		"email": "notfound@example.com",
	}
	body, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	Login(c)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid email", response["error"])

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_EmptyEmail(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	_, gormDB := setupTestDB(t)
	database.DB = gormDB

	// Create request with empty email
	loginData := map[string]string{
		"email": "",
	}
	body, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	Login(c)

	// Assert - vì email rỗng nên sẽ không tìm thấy user
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_TokenExpiration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mock, gormDB := setupTestDB(t)
	database.DB = gormDB

	// Set JWT secret for testing
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	// Mock SQL expectations
	rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
		AddRow(1, "John Doe", "john@example.com", "1234567890")

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs("john@example.com", 1).
		WillReturnRows(rows)

	// Create request
	loginData := map[string]string{
		"email": "john@example.com",
	}
	body, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	Login(c)

	// Assert
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	// Parse token và kiểm tra expiration time
	token, _ := jwt.Parse(response["token"], func(token *jwt.Token) (interface{}, error) {
		return []byte("test_secret"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		exp := int64(claims["exp"].(float64))
		expectedExp := time.Now().Add(time.Hour * 24).Unix()
		// Cho phép sai lệch 5 giây
		assert.InDelta(t, expectedExp, exp, 5)
	}
}

func TestLogin_DatabaseError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mock, gormDB := setupTestDB(t)
	database.DB = gormDB

	// Mock SQL expectations với lỗi database
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs("john@example.com", 1).
		WillReturnError(assert.AnError)

	// Create request
	loginData := map[string]string{
		"email": "john@example.com",
	}
	body, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	Login(c)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid email", response["error"])
}
