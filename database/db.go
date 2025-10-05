package database

import (
	"log"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"myapp/config"
	"myapp/models" // import model để migrate
)

var DB *gorm.DB // Biến DB toàn cục

// Hàm khởi tạo DB với retry logic
func InitDB() {
	// Lấy từ env
	user := config.GetEnv("DB_USER", "root")
	pass := config.GetEnv("DB_PASS", "mypass")
	host := config.GetEnv("DB_HOST", "127.0.0.1")
	port := config.GetEnv("DB_PORT", "3306")
	name := config.GetEnv("DB_NAME", "mydb")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	// Retry logic để kết nối database
	maxRetries := 30
	retryInterval := 2 * time.Second
	
	log.Printf("🔄 Đang kết nối database tại %s:%s...", host, port)
	
	for i := 0; i < maxRetries; i++ {
		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("⏳ Lần thử %d/%d: Database chưa sẵn sàng, thử lại sau %v...", i+1, maxRetries, retryInterval)
			if i < maxRetries-1 {
				time.Sleep(retryInterval)
				continue
			}
			log.Fatal("❌ Không thể kết nối database sau", maxRetries, "lần thử: ", err)
		} else {
			log.Println("✅ Kết nối database thành công")
			break
		}
	}

	// Tự động migrate bảng User
	log.Println("🔄 Đang migrate database...")
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("❌ Lỗi migrate database: ", err)
	}
	log.Println("✅ Migrate database hoàn thành")
}
