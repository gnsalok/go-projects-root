package main

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/gnsalok/go-project-root/go-db-data-api/handler"
	"github.com/gnsalok/go-project-root/go-db-data-api/repository"
	"github.com/gnsalok/go-project-root/go-db-data-api/router"
	"github.com/swaggo/swag/example/basic/docs"
	// Update with your module path
)

// @title Gin Couchbase API
// @version 1.0
// @description This is a sample server for managing users.
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize Couchbase
	cluster, err := gocb.Connect("couchbase://localhost", gocb.ClusterOptions{
		Username: "Administrator",
		Password: "password",
	})

	if err != nil {
		log.Fatalf("Failed to connect to Couchbase: %v", err)
	}

	// Open bucket
	bucket := cluster.Bucket("users")
	err = bucket.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		log.Fatalf("Bucket not ready: %v", err)
	}

	var user *gocb.Bucket
	// Initialize repository and handler
	userRepo := repository.NewUserRepository(user)

	userHandler := &handler.UserHandler{Repo: userRepo}

	// Setup router
	r := router.SetupRouter(userHandler)

	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "Gin Couchbase API"
	docs.SwaggerInfo.Description = "API documentation for Gin Couchbase API"
	docs.SwaggerInfo.Version = "1.0"

	// Start server
	fmt.Println("Server is running at http://localhost:8080")
	fmt.Println("Swagger docs available at http://localhost:8080/swagger/index.html")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
