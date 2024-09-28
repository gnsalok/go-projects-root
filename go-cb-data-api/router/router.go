package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gnsalok/go-project-root/go-db-data-api/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter initializes the Gin router with all routes.
func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	// User routes
	r.GET("/users/:id", userHandler.GetUserByID)

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
