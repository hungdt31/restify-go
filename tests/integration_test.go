package tests

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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"myapp/database"
	"myapp/models"
	"myapp/routes"
)

// setupTestEnvironment khởi tạo môi trường test
func setupTestEnvironment(t *testing.T) (sqlmock.Sqlmock, *gorm.DB, *gin.Engine) {
	gin.SetMode(gin.TestMode)

	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	database.DB = gormDB

	// Setup environment
	os.Setenv("JWT_SECRET", "test_secret")
	os.Setenv("APP_PORT", "8080")

	// Setup router
	router := routes.SetupRouter()

	return mock, gormDB, router
}

func teardownTestEnvironment() {
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("APP_PORT")
}

// TestUserRegistrationAndLogin test flow đăng ký và login user
func TestUserRegistrationAndLogin(t *testing.T) {
	mock, _, router := setupTestEnvironment(t)
	defer teardownTestEnvironment()

	// Step 1: Tạo user mới
	t.Run("Create User", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `users`").
			WithArgs("Alice Smith", "alice@example.com", "1234567890").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		user := models.User{
			Name:     "Alice Smith",
			Email:    "alice@example.com",
			Password: "1234567890",
		}
		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Alice Smith", response.Name)
		assert.Equal(t, "alice@example.com", response.Email)
	})

	// Step 2: Login với user vừa tạo
	t.Run("Login User", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
			AddRow(1, "Alice Smith", "alice@example.com", "1234567890")

		mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? ORDER BY").
			WithArgs("alice@example.com").
			WillReturnRows(rows)

		loginData := map[string]string{"email": "alice@example.com"}
		body, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response["token"])

		// Verify token
		token, err := jwt.Parse(response["token"], func(token *jwt.Token) (interface{}, error) {
			return []byte("test_secret"), nil
		})
		assert.NoError(t, err)
		assert.True(t, token.Valid)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestProtectedEndpointAccess test truy cập endpoint có authentication
func TestProtectedEndpointAccess(t *testing.T) {
	mock, _, router := setupTestEnvironment(t)
	defer teardownTestEnvironment()

	// Create valid token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("test_secret"))

	t.Run("Access Protected Endpoint With Valid Token", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
			AddRow(1, "John Doe", "john@example.com", "1234567890").
			AddRow(2, "Jane Doe", "jane@example.com", "0987654321")

		mock.ExpectQuery("SELECT \\* FROM `users`").
			WillReturnRows(rows)

		req, _ := http.NewRequest("GET", "/users", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Note: /users endpoint trong project hiện tại không yêu cầu auth
		// Test này để demo, có thể sửa route để thêm AuthRequired middleware
		assert.Equal(t, http.StatusOK, w.Code)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCompleteUserFlow test flow hoàn chỉnh: tạo, đọc, login
func TestCompleteUserFlow(t *testing.T) {
	mock, _, router := setupTestEnvironment(t)
	defer teardownTestEnvironment()

	// Step 1: Tạo nhiều users
	t.Run("Create Multiple Users", func(t *testing.T) {
		users := []models.User{
			{Name: "User 1", Email: "user1@example.com", Password: "1111111111"},
			{Name: "User 2", Email: "user2@example.com", Password: "2222222222"},
			{Name: "User 3", Email: "user3@example.com", Password: "3333333333"},
		}

		for _, user := range users {
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO `users`").
				WithArgs(user.Name, user.Email, user.Password).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			body, _ := json.Marshal(user)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	// Step 2: Lấy danh sách users
	t.Run("Get All Users", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
			AddRow(1, "User 1", "user1@example.com", "1111111111").
			AddRow(2, "User 2", "user2@example.com", "2222222222").
			AddRow(3, "User 3", "user3@example.com", "3333333333")

		mock.ExpectQuery("SELECT \\* FROM `users`").
			WillReturnRows(rows)

		req, _ := http.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Len(t, response, 3)
	})

	// Step 3: Login với từng user
	t.Run("Login Each User", func(t *testing.T) {
		emails := []string{"user1@example.com", "user2@example.com", "user3@example.com"}

		for i, email := range emails {
			rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
				AddRow(i+1, "User "+string(rune(i+1)), email, string(rune(i+1))+string(rune(i+1))+string(rune(i+1))+string(rune(i+1)))

			mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? ORDER BY").
				WithArgs(email).
				WillReturnRows(rows)

			loginData := map[string]string{"email": email}
			body, _ := json.Marshal(loginData)
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.NotEmpty(t, response["token"])
		}
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestHealthCheckIntegration test health check endpoint
func TestHealthCheckIntegration(t *testing.T) {
	_, _, router := setupTestEnvironment(t)
	defer teardownTestEnvironment()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "Application is running", response["message"])
}

// TestErrorHandling test xử lý lỗi
func TestErrorHandling(t *testing.T) {
	mock, _, router := setupTestEnvironment(t)
	defer teardownTestEnvironment()

	t.Run("Invalid JSON in Create User", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Database Error on Get Users", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `users`").
			WillReturnError(gorm.ErrInvalidDB)

		req, _ := http.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid Login Credentials", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? ORDER BY").
			WithArgs("notfound@example.com").
			WillReturnError(gorm.ErrRecordNotFound)

		loginData := map[string]string{"email": "notfound@example.com"}
		body, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestConcurrentRequests test xử lý concurrent requests
func TestConcurrentRequests(t *testing.T) {
	mock, _, router := setupTestEnvironment(t)
	defer teardownTestEnvironment()

	// Mock data for concurrent requests
	for i := 0; i < 10; i++ {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
			AddRow(1, "John Doe", "john@example.com", "1234567890")
		mock.ExpectQuery("SELECT \\* FROM `users`").
			WillReturnRows(rows)
	}

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			req, _ := http.NewRequest("GET", "/users", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}
