package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestLogger(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create request
	req, _ := http.NewRequest("GET", "/test", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create router with logger middleware
	router := gin.New()
	router.Use(RequestLogger())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// Execute
	router.ServeHTTP(w, req)

	// Assert - middleware không thay đổi response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test")
}

func TestRequestLogger_POST(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create POST request
	req, _ := http.NewRequest("POST", "/create", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create router with logger middleware
	router := gin.New()
	router.Use(RequestLogger())
	router.POST("/create", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"message": "created"})
	})

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRequestLogger_ErrorStatus(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create request
	req, _ := http.NewRequest("GET", "/error", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create router with logger middleware
	router := gin.New()
	router.Use(RequestLogger())
	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	})

	// Execute
	router.ServeHTTP(w, req)

	// Assert - middleware vẫn hoạt động với error status
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRequestLogger_WithQueryParams(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create request with query params
	req, _ := http.NewRequest("GET", "/search?q=test&page=1", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create router with logger middleware
	router := gin.New()
	router.Use(RequestLogger())
	router.GET("/search", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"query": c.Query("q")})
	})

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequestLogger_MultipleMiddlewares(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create request
	req, _ := http.NewRequest("GET", "/test", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create router with multiple middlewares
	router := gin.New()
	router.Use(RequestLogger())
	router.Use(func(c *gin.Context) {
		c.Set("test_key", "test_value")
		c.Next()
	})
	router.GET("/test", func(c *gin.Context) {
		value, exists := c.Get("test_key")
		assert.True(t, exists)
		assert.Equal(t, "test_value", value)
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequestLogger_DifferentPaths(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name       string
		method     string
		path       string
		statusCode int
	}{
		{"Root Path", "GET", "/", http.StatusOK},
		{"Users Path", "GET", "/users", http.StatusOK},
		{"Create User", "POST", "/users", http.StatusCreated},
		{"Update User", "PUT", "/users/1", http.StatusOK},
		{"Delete User", "DELETE", "/users/1", http.StatusNoContent},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req, _ := http.NewRequest(tc.method, tc.path, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Create router
			router := gin.New()
			router.Use(RequestLogger())
			router.Handle(tc.method, tc.path, func(c *gin.Context) {
				c.Status(tc.statusCode)
			})

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.statusCode, w.Code)
		})
	}
}

func TestRequestLogger_WithHeaders(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create request with custom headers
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "TestAgent/1.0")
	req.Header.Set("X-Custom-Header", "custom-value")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create router
	router := gin.New()
	router.Use(RequestLogger())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}
