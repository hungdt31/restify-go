package config

import (
    "log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() {
    // Ưu tiên theo APP_ENV: .env.<env> -> .env
    // Ví dụ: APP_ENV=staging => load .env.staging, nếu không có thì fallback .env
    appEnv := os.Getenv("APP_ENV")
    var candidates []string
    if appEnv != "" {
        candidates = append(candidates, ".env."+appEnv)
    }
    candidates = append(candidates, ".env")

    // Thử load lần lượt các file ứng viên
    loadedAny := false
    for _, file := range candidates {
        if _, statErr := os.Stat(file); statErr == nil {
            if err := godotenv.Overload(file); err == nil {
                log.Printf("✅ Đã load biến môi trường từ %s", file)
                loadedAny = true
                break
            }
        }
    }

    if !loadedAny {
        log.Println("⚠️ Không tìm thấy file môi trường (.env*, dùng biến môi trường hệ thống)")
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
