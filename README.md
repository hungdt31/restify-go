# GoFirstApp - RESTful API with Gin & GORM

Ứng dụng Go đầu tiên của bạn - Xây dựng RESTful API sử dụng Gin framework và GORM ORM, kết nối với MySQL database.

## 📋 Mục lục

- [Tính năng](#-tính-năng)
- [Công nghệ sử dụng](#-công-nghệ-sử-dụng)
- [Cấu trúc thư mục](#-cấu-trúc-thư-mục)
- [Yêu cầu hệ thống](#-yêu-cầu-hệ-thống)
- [Cài đặt](#-cài-đặt)
- [Cấu hình Database](#-cấu-hình-database)
- [Chạy ứng dụng](#-chạy-ứng-dụng)
 - [Chạy bằng Docker](#-chạy-bằng-docker)
- [Môi trường Staging/Production](#-môi-trường-stagingproduction)
- [Môi trường Development (Docker)](#-môi-trường-development-docker)
- [API Endpoints](#-api-endpoints)
- [Ví dụ sử dụng](#-ví-dụ-sử-dụng)

## ✨ Tính năng

- ✅ RESTful API với Gin Framework
- ✅ ORM với GORM (MySQL)
- ✅ Auto Migration Database
- ✅ Middleware Logger tùy chỉnh
- ✅ Cấu trúc MVC rõ ràng
- ✅ Hot Reload trong môi trường development
- ✅ CRUD operations cho User

## 🛠 Công nghệ sử dụng

- **Go** 1.25.1
- **Gin** v1.11.0 - HTTP web framework
- **GORM** v1.31.0 - ORM library
- **MySQL** - Database
- **Air** - Hot reload tool (development)

## 📁 Cấu trúc thư mục

```
GoFirstApp/
├── main.go                 # Entry point của ứng dụng
├── go.mod                  # Dependencies management
├── go.sum                  # Dependencies checksum
├── dev.ps1                 # Script chạy dev server (PowerShell)
├── Makefile               # Build & run tasks
├── .air.toml              # Cấu hình hot reload
├── README.md              # Tài liệu dự án
│
├── controllers/           # Business logic
│   └── userController.go  # User CRUD operations
│
├── database/              # Database setup
│   └── db.go             # Database connection & migration
│
├── middleware/            # Custom middleware
│   └── logger.go         # Request logging middleware
│
├── models/               # Data models
│   └── user.go          # User model (struct)
│
├── routes/              # Route definitions
│   ├── routes.go        # Main router setup
│   └── userRoutes.go    # User routes
│
└── tmp/                 # Build artifacts (gitignore)
    └── main.exe         # Compiled binary
```

## 📦 Yêu cầu hệ thống

- Go 1.25.1 hoặc cao hơn
- MySQL 8.0 hoặc cao hơn
- Git (optional)

## 🚀 Cài đặt

### 1. Clone hoặc tải dự án

```bash
git clone <repository-url>
cd GoFirstApp
```

### 2. Cài đặt dependencies

```bash
go mod download
```

### 3. Cài đặt Air (Hot reload tool - optional)

```bash
go install github.com/air-verse/air@latest
```

## 🗄 Cấu hình Database

### Tạo database MySQL với docker

```bash
docker run --name mysql-server -e MYSQL_ROOT_PASSWORD=rootpass -e MYSQL_DATABASE=mydb -e MYSQL_USER=myuser -e MYSQL_PASSWORD=mypass -p 3306:3306 -d mysql:8.0
```

- Connection type: MySQL
- Server name:
  - Nếu bạn kết nối từ máy host (Windows/Linux) → điền localhost hoặc 127.0.0.1
  - Nếu bạn kết nối từ một container khác trong cùng docker-compose → điền mysql-server (tên service/container).
- Authentication type: Password
- User name: myuser
- Password: mypass
- Database name: mydb
- Port: 3306 (mặc định)

### Hoặc sử dụng lệnh sql

```sql
CREATE DATABASE mydb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'myuser'@'localhost' IDENTIFIED BY 'mypass';
GRANT ALL PRIVILEGES ON mydb.* TO 'myuser'@'localhost';
FLUSH PRIVILEGES;
```

### Cấu hình kết nối

Mở file `database/db.go` và chỉnh sửa DSN connection string:

```go
dsn := "myuser:mypass@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
//      ^^^^^^  ^^^^^^      ^^^^^^^^^^^^^^^^ ^^^^
//      user    pass        host:port        database
```

## 🏃 Chạy ứng dụng

### Phương pháp 1: Chạy trực tiếp

```bash
go run main.go
```

### Phương pháp 2: Hot Reload (Development - Khuyên dùng)

```bash
# Cách 1: Dùng script PowerShell
.\dev.ps1

# Cách 2: Dùng Makefile (nếu đã cài Make)
make dev

# Cách 3: Dùng Air trực tiếp
go run github.com/air-verse/air@latest
```

### Phương pháp 3: Build và chạy

```bash
# Build
go build -o tmp/main.exe .

# Run
.\tmp\main.exe
```

Server sẽ chạy tại: **http://localhost:8080**

### Phương pháp 4: Docker Compose (Khuyên dùng cho môi trường chuẩn hóa)

1) Tạo file `.env` ở thư mục gốc với giá trị ví dụ:

```env
# App
APP_PORT=8080

# MySQL
DB_HOST=db
DB_PORT=3306
DB_NAME=mydb
DB_USER=myuser
DB_PASS=mypass
```

2) Khởi chạy docker-compose:

```bash
docker compose up -d --build
```

3) Truy cập API: `http://localhost:8080`

Ghi chú:
- Service `db` dùng `mysql:8.0` và map cổng `3306:3306` cho host.
- Ứng dụng chạy bằng binary đã build trong image (không mount source vào container) nhằm đảm bảo môi trường chạy ổn định.
- Biến `DB_HOST` mặc định được đặt là `db` để app kết nối MySQL nội bộ qua network của compose.

## 🌐 Môi trường Staging/Production

Từ nay ứng dụng hỗ trợ cấu hình đa môi trường qua biến `APP_ENV` và các file `.env.*`.

### Cơ chế nạp biến môi trường

- App đọc `APP_ENV` và tìm file `.env.<APP_ENV>` trước, nếu không tồn tại sẽ fallback `.env`.
- Ví dụ: `APP_ENV=staging` → ưu tiên `.env.staging`.

### Staging

1) Tạo file môi trường dựa trên mẫu:

```bash
copy .env.staging.example .env.staging  # Windows PowerShell/cmd
# hoặc
cp .env.staging.example .env.staging     # Linux/macOS
```

2) Khởi chạy môi trường staging:

```bash
make staging-up
# Dừng và xóa: make staging-down
```

- App sẽ chạy tại `http://localhost:8081`
- MySQL staging sẽ map `3307:3306` để không xung đột với dev.

### Production

1) Tạo file `.env.production` (không commit) dựa trên mẫu:

```bash
copy .env.production.example .env.production
```

2) Chạy production (app-only, không kèm MySQL nội bộ):

```bash
make prod-up    # build & run image với APP_ENV=production
# Dừng: make prod-down
```

Ghi chú Production:
- Sử dụng DB managed (RDS/Cloud SQL, v.v.). Thiết lập `DB_HOST`, `DB_USER`, `DB_PASS`, `DB_NAME` trong `.env.production`.
- Image build tối giản từ Dockerfile multi-stage, chạy user không phải root.

## 🧑‍💻 Môi trường Development (Docker)

Sử dụng `docker-compose.dev.yml` để chạy môi trường dev có hot reload (Air).

1) Tạo file môi trường:

```bash
copy .env.development.example .env.development  # Windows
# hoặc
cp .env.development.example .env.development    # Linux/macOS
```

2) Khởi chạy dev bằng Docker:

```bash
make docker-dev-up
# Dừng/xoá: make docker-dev-down
```

Ghi chú:
- App map cổng `8082:8080`, MySQL dev map `3308:3306` để tránh xung đột.
- Service `app` dùng image `golang:1.22`, bind-mount code `./:/src` và chạy `air` để hot reload.

## 📡 API Endpoints

| Method | Endpoint    | Description          | Request Body           |
|--------|-------------|----------------------|------------------------|
| POST   | /users      | Tạo user mới         | `{"name":"...", "email":"..."}` |
| GET    | /users      | Lấy danh sách users  | -                      |

### Request/Response Models

**User Model:**
```json
{
  "id": 1,
  "name": "Alice",
  "email": "alice@example.com"
}
```

## 💡 Ví dụ sử dụng

### 1. Tạo User mới

**PowerShell:**
```powershell
curl.exe -X POST http://localhost:8080/users `
  -H "Content-Type: application/json" `
  -d '{\"name\":\"Alice\", \"email\":\"alice@example.com\"}'
```

**Bash/Linux:**
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice", "email":"alice@example.com"}'
```

**Response:**
```json
{
  "id": 1,
  "name": "Alice",
  "email": "alice@example.com"
}
```

### 2. Lấy danh sách Users

**PowerShell:**
```powershell
curl.exe http://localhost:8080/users
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Alice",
    "email": "alice@example.com"
  },
  {
    "id": 2,
    "name": "Bob",
    "email": "bob@example.com"
  }
]
```

## 🔧 Development

### Hot Reload với Air

Air sẽ tự động phát hiện thay đổi trong các file `.go` và restart server. Cấu hình trong `.air.toml`:

- **Watched extensions:** `.go`, `.tpl`, `.tmpl`, `.html`
- **Excluded dirs:** `tmp`, `vendor`, `testdata`
- **Build delay:** 1 giây

### Thêm Middleware mới

Tạo file trong `middleware/`:

```go
package middleware

import "github.com/gin-gonic/gin"

func YourMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Logic trước khi xử lý request
        c.Next()
        // Logic sau khi xử lý request
    }
}
```

Đăng ký trong `routes/routes.go`:

```go
r.Use(middleware.YourMiddleware())
```

### Thêm Model mới

Tạo file trong `models/`, ví dụ `product.go`:

```go
package models

type Product struct {
    ID    uint    `json:"id" gorm:"primaryKey"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

Thêm vào Auto Migration trong `database/db.go`:

```go
DB.AutoMigrate(&models.User{}, &models.Product{})
```

## 📝 TODO

- [x] Thêm authentication (JWT)
- [ ] Thêm validation chi tiết hơn
- [ ] Thêm pagination cho GET endpoints
- [ ] Thêm unit tests
- [ ] Thêm Docker support
- [ ] Thêm logging vào file
- [ ] Thêm environment variables config

## 🤝 Contributing

Contributions, issues và feature requests đều được chào đón!

## 📄 License

This project is [MIT](LICENSE) licensed.

## 👤 Author

Your Name - [GitHub](https://github.com/yourusername)

---

**Happy Coding! 🚀**
