# Unit Testing Guide - GoFirstApp

## 📋 Tổng quan

Project này đã được tích hợp unit tests và integration tests toàn diện cho tất cả các thành phần chính.

## 🧪 Cấu trúc Tests

```
GoFirstApp/
├── controllers/
│   ├── userController.go
│   ├── userController_test.go      # Tests cho user controller
│   ├── authController.go
│   └── authController_test.go      # Tests cho auth controller
├── middleware/
│   ├── auth.go
│   ├── auth_test.go                # Tests cho auth middleware
│   ├── logger.go
│   └── logger_test.go              # Tests cho logger middleware
├── models/
│   ├── user.go
│   └── user_test.go                # Tests cho user model
├── routes/
│   ├── routes.go
│   └── routes_test.go              # Tests cho router setup
└── tests/
    └── integration_test.go         # Integration tests
```

## 📦 Dependencies cho Testing

Các thư viện testing đã được cài đặt:

- **testify** - Assertions và mocking framework
- **go-sqlmock** - Mock SQL database cho testing
- **httptest** - HTTP testing utilities (built-in Go)

## 🚀 Chạy Tests

### 1. Chạy tất cả tests

```powershell
# Sử dụng Makefile
make test

# Hoặc sử dụng go test trực tiếp
go test ./...
```

### 2. Chạy tests với output chi tiết

```powershell
make test-verbose

# Hoặc
go test ./... -v
```

### 3. Chạy tests với coverage report

```powershell
make test-coverage

# Kết quả sẽ tạo file coverage.html để xem trong browser
```

### 4. Chạy unit tests riêng biệt

```powershell
make test-unit

# Hoặc test từng package
go test ./controllers/... -v
go test ./models/... -v
go test ./middleware/... -v
```

### 5. Chạy integration tests

```powershell
make test-integration

# Hoặc
go test ./tests/... -v
```

### 6. Chạy test cho một file cụ thể

```powershell
# Test user controller
go test ./controllers/userController_test.go ./controllers/userController.go -v

# Test auth middleware
go test ./middleware/auth_test.go ./middleware/auth.go -v
```

### 7. Chạy một test function cụ thể

```powershell
# Chạy test CreateUser_Success
go test ./controllers -run TestCreateUser_Success -v

# Chạy test Login_Success
go test ./controllers -run TestLogin_Success -v
```

## 📊 Test Coverage

Để xem coverage report:

```powershell
# Tạo coverage report
make test-coverage

# Mở file coverage.html trong browser
start coverage.html  # Windows
```

Target coverage: **>80%** cho tất cả các package chính

## 🧩 Các Test Cases

### Controllers Tests

#### UserController (`controllers/userController_test.go`)
- ✅ `TestCreateUser_Success` - Tạo user thành công
- ✅ `TestCreateUser_InvalidJSON` - Xử lý JSON không hợp lệ
- ✅ `TestCreateUser_DatabaseError` - Xử lý lỗi database
- ✅ `TestGetUsers_Success` - Lấy danh sách users
- ✅ `TestGetUsers_DatabaseError` - Xử lý lỗi khi query
- ✅ `TestGetUsers_EmptyResult` - Xử lý kết quả rỗng

#### AuthController (`controllers/authController_test.go`)
- ✅ `TestLogin_Success` - Login thành công
- ✅ `TestLogin_InvalidJSON` - JSON không hợp lệ
- ✅ `TestLogin_UserNotFound` - User không tồn tại
- ✅ `TestLogin_EmptyEmail` - Email rỗng
- ✅ `TestLogin_TokenExpiration` - Kiểm tra token expiration
- ✅ `TestLogin_DatabaseError` - Xử lý lỗi database

### Middleware Tests

#### Auth Middleware (`middleware/auth_test.go`)
- ✅ `TestAuthRequired_ValidToken` - Token hợp lệ
- ✅ `TestAuthRequired_MissingAuthHeader` - Thiếu Authorization header
- ✅ `TestAuthRequired_InvalidAuthHeaderFormat` - Format header sai
- ✅ `TestAuthRequired_InvalidTokenFormat` - Token không hợp lệ
- ✅ `TestAuthRequired_ExpiredToken` - Token hết hạn
- ✅ `TestAuthRequired_WrongSecret` - Sai JWT secret
- ✅ `TestAuthRequired_BearerCaseInsensitive` - Kiểm tra case sensitivity
- ✅ `TestAuthRequired_TokenWithoutExpiration` - Token không có expiration

#### Logger Middleware (`middleware/logger_test.go`)
- ✅ `TestRequestLogger` - Logger middleware hoạt động
- ✅ `TestRequestLogger_POST` - Log POST requests
- ✅ `TestRequestLogger_ErrorStatus` - Log error status
- ✅ `TestRequestLogger_WithQueryParams` - Log với query params
- ✅ `TestRequestLogger_MultipleMiddlewares` - Hoạt động với nhiều middleware
- ✅ `TestRequestLogger_DifferentPaths` - Log các paths khác nhau
- ✅ `TestRequestLogger_WithHeaders` - Log với custom headers

### Models Tests

#### User Model (`models/user_test.go`)
- ✅ `TestUserModel_Structure` - Cấu trúc model
- ✅ `TestUserModel_JSONTags` - JSON tags
- ✅ `TestUserModel_BeforeCreate` - Hooks trước khi tạo
- ✅ `TestUserModel_EmptyUser` - User rỗng
- ✅ `TestUserModel_PartialUser` - User thiếu field
- ✅ `TestUserModel_Update` - Cập nhật user
- ✅ `TestUserModel_GormTags` - GORM tags (unique, etc.)
- ✅ `TestUserModel_PhoneMaxLength` - Kiểm tra độ dài phone
- ✅ `TestUserModel_AutoIncrement` - Auto increment ID

### Routes Tests

#### Router (`routes/routes_test.go`)
- ✅ `TestSetupRouter` - Setup router
- ✅ `TestHealthCheckEndpoint` - Health check endpoint
- ✅ `TestRouterMiddleware` - Middleware được apply
- ✅ `TestNonExistentRoute` - Route không tồn tại
- ✅ `TestMethodNotAllowed` - Method không được phép
- ✅ `TestHealthCheckResponseFormat` - Format response
- ✅ `TestRouterWithMultipleRequests` - Multiple requests
- ✅ `TestRouterConcurrentRequests` - Concurrent requests

### Integration Tests

#### Complete Flows (`tests/integration_test.go`)
- ✅ `TestUserRegistrationAndLogin` - Flow đăng ký và login
- ✅ `TestProtectedEndpointAccess` - Truy cập endpoint có auth
- ✅ `TestCompleteUserFlow` - Flow hoàn chỉnh: tạo, đọc, login
- ✅ `TestHealthCheckIntegration` - Integration health check
- ✅ `TestErrorHandling` - Xử lý lỗi toàn diện
- ✅ `TestConcurrentRequests` - Concurrent requests handling

## 🎯 Best Practices

### 1. Naming Convention
- Test file: `<filename>_test.go`
- Test function: `Test<FunctionName>_<Scenario>`
- Helper function: `setup<Something>` hoặc `mock<Something>`

### 2. Test Structure (AAA Pattern)
```go
func TestSomething(t *testing.T) {
    // Arrange - Setup
    gin.SetMode(gin.TestMode)
    mock, db := setupTestDB(t)
    
    // Act - Execute
    result := DoSomething()
    
    // Assert - Verify
    assert.Equal(t, expected, result)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

### 3. Table-Driven Tests
```go
func TestMultipleScenarios(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected string
    }{
        {"Case 1", "input1", "output1"},
        {"Case 2", "input2", "output2"},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := Process(tc.input)
            assert.Equal(t, tc.expected, result)
        })
    }
}
```

### 4. Mock Database
Sử dụng `go-sqlmock` để mock database:
```go
mock.ExpectQuery("SELECT \\* FROM `users`").
    WillReturnRows(rows)
```

### 5. Cleanup
Luôn cleanup sau test:
```go
defer teardownTestEnvironment()
defer os.Unsetenv("JWT_SECRET")
```

## 📝 Thêm Tests Mới

### 1. Thêm test cho controller mới

```go
// controllers/productController_test.go
package controllers

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateProduct_Success(t *testing.T) {
    // Setup
    // ...
    
    // Execute
    // ...
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
}
```

### 2. Thêm integration test

```go
// tests/product_integration_test.go
package tests

func TestProductFlow(t *testing.T) {
    mock, _, router := setupTestEnvironment(t)
    defer teardownTestEnvironment()
    
    // Test complete flow
    // ...
}
```

## 🐛 Debugging Tests

### Chạy test với debug output

```powershell
go test ./... -v -count=1
```

### Xem chi tiết lỗi

```powershell
go test ./controllers -v -run TestCreateUser_Success
```

### Skip cache khi test

```powershell
go test ./... -count=1
```

## 🔄 CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.25'
      - run: go test ./... -v -cover
```

## 📚 Resources

- [Go Testing Documentation](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Go-SQLMock Documentation](https://github.com/DATA-DOG/go-sqlmock)
- [Gin Testing Guide](https://gin-gonic.com/docs/testing/)

## ✅ Checklist khi thêm feature mới

- [ ] Viết unit tests cho các functions mới
- [ ] Đảm bảo coverage >= 80%
- [ ] Viết integration test nếu cần
- [ ] Chạy tất cả tests trước khi commit
- [ ] Update documentation nếu cần

---

**Happy Testing! 🧪**
