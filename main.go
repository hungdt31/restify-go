package main

import (
	"myapp/config"
	"myapp/database"
	"myapp/routes"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Kết nối DB
	database.InitDB()

	// Setup routes
	r := routes.SetupRouter()

	// Lấy port từ env
	port := config.GetEnv("APP_PORT", "8080")
	r.Run(":" + port)
}
