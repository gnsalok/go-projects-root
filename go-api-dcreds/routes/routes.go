// routes/routes.go
package routes

import (
	"test-go/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	dynCreds := router.Group("/dyncreds")
	{
		dynCreds.POST("", handlers.CreateDynamicCredentialHandler)
		dynCreds.GET("/:dyncredId", handlers.GetDynamicCredentialHandler)
		dynCreds.PUT("/:dyncredId", handlers.UpdateDynamicCredentialHandler)
		dynCreds.DELETE("/:dyncredId", handlers.DeleteDynamicCredentialHandler)
		dynCreds.PATCH("/:dyncredId", handlers.PatchDynamicCredentialHandler)
	}
}
