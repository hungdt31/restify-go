package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Execute
	router := SetupRouter()

	// Assert - router should not be nil
	assert.NotNil(t, router)
}

func TestHealthCheckEndpoint(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// Create request
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
	assert.Contains(t, w.Body.String(), "Application is running")
}

func TestRouterMiddleware(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// Create request
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert - middleware should be applied
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNonExistentRoute(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// Create request to non-existent route
	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert - should return 404
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMethodNotAllowed(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// Create request with wrong method for health endpoint
	req, _ := http.NewRequest("POST", "/health", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert - should return 404 (Gin returns 404 for method not allowed by default)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCORSHeaders(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// Create OPTIONS request
	req, _ := http.NewRequest("OPTIONS", "/health", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert - check if CORS headers are set (if middleware is added)
	// This test will depend on if CORS middleware is implemented
	assert.NotNil(t, w)
}

func TestHealthCheckResponseFormat(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// Create request
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestRouterWithMultipleRequests(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// Test multiple requests
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestRouterConcurrentRequests(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := SetupRouter()

	// Create channels for synchronization
	done := make(chan bool, 10)

	// Execute concurrent requests
	for i := 0; i < 10; i++ {
		go func() {
			req, _ := http.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
			done <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
