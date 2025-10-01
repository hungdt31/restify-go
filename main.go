package main

import (
	"myapp/database"
	"myapp/routes"
)

func main() {
	// Kết nối DB
	database.InitDB()

	// Setup routes
	r := routes.SetupRouter()

	// Run server
	r.Run(":8080")
}
