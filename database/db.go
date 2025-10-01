package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"myapp/models" // import model để migrate
)

var DB *gorm.DB // Biến DB toàn cục

// Hàm khởi tạo DB
func InitDB() {
	dsn := "myuser:mypass@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Không kết nối được database: ", err)
	}

	// Tự động migrate bảng User
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Lỗi migrate database: ", err)
	}
}
