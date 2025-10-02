package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	// Load .env file vào environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Không tìm thấy .env file, sử dụng biến môi trường hệ thống")
	} else {
		log.Println("✅ Đã load .env file thành công")
	}
}

// GetEnv gets environment variable with fallback value
func GetEnv(key, fallback string) string {
	// Lấy từ environment variables (đã được load từ .env)
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// MustGetEnv gets environment variable or panics if not found
func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("❌ Biến môi trường %s là bắt buộc nhưng không được thiết lập", key)
	}
	return value
}
