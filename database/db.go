package database

import (
	"fmt"
	"log"
	"time"

	// "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"myapp/config"
)

var DB *gorm.DB // Biến DB toàn cục

// Hàm khởi tạo DB với retry logic và migration
func InitDB() {
	// Lấy từ env
	user := config.GetEnv("DB_USER", "root")
	pass := config.GetEnv("DB_PASS", "mypass")
	host := config.GetEnv("DB_HOST", "127.0.0.1")
	port := config.GetEnv("DB_PORT", "3306")
	name := config.GetEnv("DB_NAME", "mydb")

	// DSN cho GORM
	gormDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	// Retry logic để kết nối database
	maxRetries := 30
	retryInterval := 2 * time.Second

	log.Printf("🔄 Đang kết nối database tại %s:%s...", host, port)

	for i := range maxRetries {
		var err error
		DB, err = gorm.Open(mysql.Open(gormDsn), &gorm.Config{})
		if err != nil {
			log.Printf("⏳ Lần thử %d/%d: Database chưa sẵn sàng, thử lại sau %v...", i+1, maxRetries, retryInterval)
			if i < maxRetries-1 {
				time.Sleep(retryInterval)
				continue
			}
			log.Fatal("❌ Không thể kết nối database sau ", maxRetries, " lần thử: ", err)
		} else {
			log.Println("✅ Kết nối database thành công")
			break
		}
	}

	// --- Bắt đầu phần Migration ---
	// log.Println("🔄 Đang chạy database migrations...")

	// // URL cho golang-migrate/migrate
	// migrateDatabaseURL := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s",
	// 	user, pass, host, port, name)

	// migrationsPath := "file://database/migrations"

	// m, err := migrate.New(migrationsPath, migrateDatabaseURL)
	// if err != nil {
	// 	log.Fatal("❌ Lỗi khi khởi tạo instance migrate: ", err)
	// }

	// if err := m.Up(); err != nil && err != migrate.ErrNoChange {
	// 	log.Fatal("❌ Lỗi khi chạy migrate up: ", err)
	// }

	// log.Println("✅ Migrate database hoàn thành")
}
