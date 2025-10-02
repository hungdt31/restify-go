package database

import (
	"log"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"myapp/config"
	"myapp/models" // import model để migrate
)

var DB *gorm.DB // Biến DB toàn cục

// Hàm khởi tạo DB
func InitDB() {
	// Lấy từ env
	user := config.GetEnv("DB_USER", "root")
	pass := config.GetEnv("DB_PASS", "mypass")
	host := config.GetEnv("DB_HOST", "127.0.0.1")
	port := config.GetEnv("DB_PORT", "3306")
	name := config.GetEnv("DB_NAME", "mydb")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Không kết nối được database: ", err)
	} else {
		log.Println("✅ Kết nối database thành công")
	}

	// Tự động migrate bảng User
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Lỗi migrate database: ", err)
	}
}
