// main.go
package main

import (
	"test-go/middleware"
	"test-go/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Apply middlewares
	router.Use(middleware.LoggerMiddleware())
	// router.Use(middleware.AuthenticationMiddleware()) // Uncomment if authentication is implemented

	// Setup routes
	routes.SetupRoutes(router)

	// Start server on port 8080
	router.Run(":8080")
}
