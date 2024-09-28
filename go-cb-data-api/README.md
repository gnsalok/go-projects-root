Creating a REST API in Golang using the Gin framework that interacts with a local Couchbase instance involves several key steps. This comprehensive guide will walk you through:

1. **Setting up the Go project with Gin and Couchbase integration**
2. **Implementing the REST API**
3. **Writing unit tests using Mockery**
4. **Generating Swagger documentation**
5. **Creating a Makefile to automate tasks**

By the end of this guide, you'll have a functional API with comprehensive testing and documentation, all orchestrated through a Makefile.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [1. Create a REST API in Golang Using Gin That Reads Data from Local Couchbase](#1-create-a-rest-api-in-golang-using-gin-that-reads-data-from-local-couchbase)
    - [a. Set Up Your Go Project](#a-set-up-your-go-project)
    - [b. Directory Structure](#b-directory-structure)
    - [c. Install Necessary Packages](#c-install-necessary-packages)
    - [d. Implement the API Components](#d-implement-the-api-components)
        - [i. Define the User Model](#i-define-the-user-model)
        - [ii. Create the Repository Layer](#ii-create-the-repository-layer)
        - [iii. Implement the Handler](#iii-implement-the-handler)
        - [iv. Setup the Router](#iv-setup-the-router)
        - [v. Initialize Couchbase and Start the Server](#v-initialize-couchbase-and-start-the-server)
    - [e. Populate Couchbase with Sample Data](#e-populate-couchbase-with-sample-data)
3. [2. Write Unit Test Cases for the REST API Using Mockery](#2-write-unit-test-cases-for-the-rest-api-using-mockery)
    - [a. Install Test Dependencies and Mockery](#a-install-test-dependencies-and-mockery)
    - [b. Generate Mocks with Mockery](#b-generate-mocks-with-mockery)
    - [c. Write the Handler Tests](#c-write-the-handler-tests)
4. [3. Setup Swagger Documentation Generation for the Handler](#3-setup-swagger-documentation-generation-for-the-handler)
    - [a. Install Swag CLI](#a-install-swag-cli)
    - [b. Annotate Your Handlers](#b-annotate-your-handlers)
    - [c. Generate Swagger Documentation](#c-generate-swagger-documentation)
    - [d. Serve Swagger UI](#d-serve-swagger-ui)
5. [4. Create a Makefile to Automate Tasks](#4-create-a-makefile-to-automate-tasks)
    - [a. Install Docker (If Not Already Installed)](#a-install-docker-if-not-already-installed)
    - [b. Create the Makefile](#b-create-the-makefile)
    - [c. Explanation of Makefile Targets](#c-explanation-of-makefile-targets)
    - [d. Usage Examples](#d-usage-examples)
6. [5. Running the Application](#5-running-the-application)
7. [6. Additional Tips](#6-additional-tips)

---

## Prerequisites

Before diving into the implementation, ensure you have the following installed on your machine:

- **Go**: [Download and install Go](https://golang.org/dl/)
- **Couchbase Server**: [Download and install Couchbase](https://www.couchbase.com/downloads)
- **Docker**: [Download and install Docker](https://www.docker.com/get-started) (for managing Couchbase via Docker)
- **Swag CLI**: Used for generating Swagger docs
    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```
    Ensure `$GOPATH/bin` is in your `PATH` to use the `swag` command.
- **Mockery**: Used for generating mocks
    ```bash
    go install github.com/vektra/mockery/v2@latest
    ```
    Ensure `$GOPATH/bin` is in your `PATH` to use the `mockery` command.

---

## 1. Create a REST API in Golang Using Gin That Reads Data from Local Couchbase

### a. Set Up Your Go Project

1. **Initialize the Go Module**

    Open your terminal and run:

    ```bash
    mkdir gin-couchbase-api
    cd gin-couchbase-api
    go mod init github.com/yourusername/gin-couchbase-api
    ```

    Replace `yourusername` with your actual GitHub username or desired module path.

### b. Directory Structure

Organize your project with the following structure:

```
gin-couchbase-api/
├── cmd/
│   └── main.go
├── handler/
│   └── user.go
├── model/
│   └── user.go
├── repository/
│   ├── user_repository.go
│   └── mock_user_repository.go
├── router/
│   └── router.go
├── docs/
│   └── (generated Swagger docs)
├── tests/
│   └── handler/
│       └── user_test.go
├── go.mod
├── go.sum
├── Makefile
└── ... (other files)
```

### c. Install Necessary Packages

Install the required dependencies using `go get`:

```bash
go get github.com/gin-gonic/gin
go get github.com/couchbase/gocb/v2
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/files
go get github.com/swaggo/swag/cmd/swag
go get github.com/stretchr/testify
go get github.com/vektra/mockery/v2
```

### d. Implement the API Components

#### i. Define the User Model

_Create `model/user.go`:_

```go
package model

// User represents a user entity in the system.
type User struct {
    ID    string `json:"id" couchbase:"id"`
    Name  string `json:"name" couchbase:"name"`
    Email string `json:"email" couchbase:"email"`
}
```

#### ii. Create the Repository Layer

_Create `repository/user_repository.go`:_

```go
package repository

import (
    "context"
    "errors"

    "github.com/couchbase/gocb/v2"
    "github.com/yourusername/gin-couchbase-api/model"
)

var (
    // ErrNotFound is returned when a user is not found in the database.
    ErrNotFound = errors.New("user not found")
)

// UserRepository defines the methods that any
// data storage provider must implement to get User information.
type UserRepository interface {
    GetUserByID(ctx context.Context, id string) (*model.User, error)
}

// userRepository implements UserRepository interface.
type userRepository struct {
    bucket *gocb.Bucket
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(bucket *gocb.Bucket) UserRepository {
    return &userRepository{bucket: bucket}
}

// GetUserByID retrieves a user by their ID from Couchbase.
func (r *userRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
    var user model.User
    collection := r.bucket.DefaultCollection()
    getResult, err := collection.Get(id, nil)
    if err != nil {
        if errors.Is(err, gocb.ErrDocumentNotFound) {
            return nil, ErrNotFound
        }
        return nil, err
    }
    err = getResult.Content(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

#### iii. Implement the Handler

_Create `handler/user.go`:_

```go
package handler

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/yourusername/gin-couchbase-api/model"
    "github.com/yourusername/gin-couchbase-api/repository"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
    Repo repository.UserRepository
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Retrieve user details by ID from Couchbase
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
    id := c.Param("id")

    user, err := h.Repo.GetUserByID(c.Request.Context(), id)
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    c.JSON(http.StatusOK, user)
}
```

**Note:** Ensure to import the `errors` package in `handler/user.go` if not already:

```go
import "errors"
```

#### iv. Setup the Router

_Create `router/router.go`:_

```go
package router

import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"

    "github.com/yourusername/gin-couchbase-api/handler"
    _ "github.com/yourusername/gin-couchbase-api/docs" // Swagger docs
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
```

#### v. Initialize Couchbase and Start the Server

_Create `cmd/main.go`:_

```go
package main

import (
    "fmt"
    "log"

    "github.com/couchbase/gocb/v2"
    "github.com/yourusername/gin-couchbase-api/handler"
    "github.com/yourusername/gin-couchbase-api/repository"
    "github.com/yourusername/gin-couchbase-api/router"
    "github.com/swaggo/swag/example/basic/docs" // Update with your module path
    _ "github.com/yourusername/gin-couchbase-api/docs"
)

// @title Gin Couchbase API
// @version 1.0
// @description This is a sample server for managing users.
// @host localhost:8080
// @BasePath /
func main() {
    // Initialize Couchbase
    cluster, err := gocb.Connect("localhost", gocb.ClusterOptions{
        Username: "Administrator",
        Password: "password",
    })
    if err != nil {
        log.Fatalf("Failed to connect to Couchbase: %v", err)
    }

    bucket := cluster.Bucket("users")
    err = bucket.WaitUntilReady(nil)
    if err != nil {
        log.Fatalf("Bucket not ready: %v", err)
    }

    // Initialize repository and handler
    userRepo := repository.NewUserRepository(bucket)
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
    log.Fatal(r.Run(":8080"))
}
```

**Notes:**

- Ensure you replace `github.com/swaggo/swag/example/basic/docs` with the correct import path based on your module name.
- The `@title`, `@version`, etc., are Swagger annotations necessary for generating the documentation.

### e. Populate Couchbase with Sample Data

Before running the API, ensure that Couchbase is running locally and contains a `users` bucket with some sample user documents.

1. **Start Couchbase Locally**

    If you haven't installed Couchbase, download and install it from [Couchbase Downloads](https://www.couchbase.com/downloads). After installation:

    - Start Couchbase Server.
    - Create a bucket named `users`.
    - Insert sample documents into the `users` bucket.

2. **Insert Sample Documents**

    You can use the Couchbase Web Console to insert documents manually:

    - Navigate to the Web Console (typically at `http://localhost:8091`).
    - Go to the `Users` bucket.
    - Insert a document with the following JSON:

      ```json
      {
        "id": "user1",
        "name": "John Doe",
        "email": "john.doe@example.com"
      }
      ```

    Alternatively, use the Couchbase CLI or SDK to insert documents programmatically.

---

## 2. Write Unit Test Cases for the REST API Using Mockery

Unit testing is crucial for ensuring the reliability and maintainability of your code. To effectively unit test the REST API, we'll mock interactions with Couchbase to isolate the handler logic.

### a. Install Test Dependencies and Mockery

Ensure you have the testing packages installed. We'll use [Testify](https://github.com/stretchr/testify) for assertions and [Mockery](https://github.com/vektra/mockery) for generating mocks.

```bash
go get github.com/stretchr/testify
go get github.com/vektra/mockery/v2
```

### b. Generate Mocks with Mockery

Mockery is a tool for generating mocks of interfaces. We'll use it to generate a mock for the `UserRepository` interface.

1. **Generate the Mock**

    Navigate to the root of your project and run:

    ```bash
    mockery --name=UserRepository --output=repository/mocks --outpkg=mocks --case=underscore
    ```

    **Explanation:**

    - `--name=UserRepository`: Specifies the interface to mock.
    - `--output=repository/mocks`: Specifies the output directory for the generated mock.
    - `--outpkg=mocks`: Sets the package name for the mocks.
    - `--case=underscore`: Formats the generated file names with underscores.

    This command will generate a mock file `repository/mocks/mock_user_repository.go`.

2. **Directory Structure After Mock Generation**

    ```
    gin-couchbase-api/
    ├── repository/
    │   ├── user_repository.go
    │   └── mocks/
    │       └── mock_user_repository.go
    └── ...
    ```

### c. Write the Handler Tests

_Create `tests/handler/user_test.go`:_

```go
package handler_test

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/yourusername/gin-couchbase-api/handler"
    "github.com/yourusername/gin-couchbase-api/model"
    "github.com/yourusername/gin-couchbase-api/repository/mocks"
)

// TestGetUserByID tests the GetUserByID handler.
func TestGetUserByID(t *testing.T) {
    // Initialize the mock repository
    mockRepo := new(mocks.UserRepository)

    // Sample user data
    sampleUser := &model.User{
        ID:    "user1",
        Name:  "John Doe",
        Email: "john.doe@example.com",
    }

    // Setup expectations
    mockRepo.On("GetUserByID", mock.Anything, "user1").Return(sampleUser, nil)
    mockRepo.On("GetUserByID", mock.Anything, "user2").Return(nil, repository.ErrNotFound)
    mockRepo.On("GetUserByID", mock.Anything, "user3").Return(nil, errors.New("database error"))

    // Initialize handler with mock repository
    userHandler := &handler.UserHandler{Repo: mockRepo}

    // Setup Gin router with the handler
    router := gin.Default()
    router.GET("/users/:id", userHandler.GetUserByID)

    // Define test cases
    testCases := []struct {
        name         string
        userID       string
        expectedCode int
        expectedBody interface{}
    }{
        {
            name:         "Existing User",
            userID:       "user1",
            expectedCode: http.StatusOK,
            expectedBody: sampleUser,
        },
        {
            name:         "Non-Existing User",
            userID:       "user2",
            expectedCode: http.StatusNotFound,
            expectedBody: gin.H{"error": "User not found"},
        },
        {
            name:         "Database Error",
            userID:       "user3",
            expectedCode: http.StatusInternalServerError,
            expectedBody: gin.H{"error": "Internal server error"},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create a request
            req, err := http.NewRequest("GET", "/users/"+tc.userID, nil)
            assert.NoError(t, err)

            // Create a response recorder
            rr := httptest.NewRecorder()

            // Serve the HTTP request
            router.ServeHTTP(rr, req)

            // Check the status code
            assert.Equal(t, tc.expectedCode, rr.Code)

            // Check the response body
            if tc.expectedCode == http.StatusOK {
                var user model.User
                err := json.Unmarshal(rr.Body.Bytes(), &user)
                assert.NoError(t, err)
                assert.Equal(t, tc.expectedBody, &user)
            } else {
                var respBody map[string]string
                err := json.Unmarshal(rr.Body.Bytes(), &respBody)
                assert.NoError(t, err)
                assert.Equal(t, tc.expectedBody, respBody)
            }
        })
    }

    // Assert that the expectations were met
    mockRepo.AssertExpectations(t)
}
```

**Explanation:**

- **Import Statements:** Import necessary packages, including the generated mock.
- **Test Function:** `TestGetUserByID` sets up the mock repository, defines test cases, and asserts the responses.
- **Test Cases:**
  - **Existing User:** Expects a 200 OK response with user data.
  - **Non-Existing User:** Expects a 404 Not Found response.
  - **Database Error:** Expects a 500 Internal Server Error response.

**Note:** Ensure that the import path `github.com/yourusername/gin-couchbase-api/repository/mocks` correctly points to the generated mocks.

---

## 3. Setup Swagger Documentation Generation for the Handler

Swagger (now known as OpenAPI) allows you to document your REST API. We'll use `swaggo/swag` to generate Swagger docs from annotations in your code.

### a. Install Swag CLI

Ensure `swag` is installed globally:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Ensure `$GOPATH/bin` is in your `PATH` to use the `swag` command.

### b. Annotate Your Handlers

Ensure your handler functions have proper Swagger annotations. This was partially done in the handler earlier.

**Example in `handler/user.go`:**

```go
// GetUserByID godoc
// @Summary Get a user by ID
// @Description Retrieve user details by ID from Couchbase
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
    // Handler implementation
}
```

Additionally, add general API information in `cmd/main.go` as shown earlier with Swagger annotations.

### c. Generate Swagger Documentation

Run the `swag` command to generate the Swagger docs.

1. **Navigate to Project Root**

    ```bash
    cd gin-couchbase-api
    ```

2. **Initialize Swagger Docs**

    ```bash
    swag init --parseDependency --parseInternal
    ```

    **Options:**

    - `--parseDependency`: Parses dependencies.
    - `--parseInternal`: Parses internal packages.

    This command will generate a `docs` folder containing `docs.go` and Swagger JSON/YAML files.

**Note:** Ensure that the `docs` directory is not tracked by version control or is appropriately managed.

### d. Serve Swagger UI

Ensure that the Swagger UI is accessible via your API. This was set up in the `router/router.go`:

```go
// Swagger route
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

After running the server, navigate to `http://localhost:8080/swagger/index.html` to view the Swagger UI.

---

## 4. Create a Makefile to Automate Tasks

A Makefile can help automate the generation of Swagger docs, launching Couchbase locally, and running unit tests.

### a. Install Docker (If Not Already Installed)

To launch Couchbase locally using Docker, ensure Docker is installed on your machine. Download from [Docker Official Site](https://www.docker.com/get-started).

### b. Create the Makefile

_Create `Makefile` in the project root:_

```makefile
# Makefile for Gin Couchbase API

.PHONY: all run test swagger docker-up docker-down generate-mocks clean

APP_NAME=gin-couchbase-api
DOCKER_COUCHBASE_NAME=couchbase-local
COUCHBASE_IMAGE=couchbase
COUCHBASE_VERSION=latest
COUCHBASE_ADMIN=Administrator
COUCHBASE_PASSWORD=password

all: docker-up run

# Launch Couchbase using Docker
docker-up:
	@echo "Starting Couchbase Docker container..."
	docker run -d \
		--name $(DOCKER_COUCHBASE_NAME) \
		-p 8091-8094:8091-8094 \
		-p 11210:11210 \
		-e COUCHBASE_ADMINISTRATOR_USERNAME=$(COUCHBASE_ADMIN) \
		-e COUCHBASE_ADMINISTRATOR_PASSWORD=$(COUCHBASE_PASSWORD) \
		$(COUCHBASE_IMAGE):$(COUCHBASE_VERSION)
	@echo "Waiting for Couchbase to initialize..."
	sleep 20
	@echo "Couchbase is up."

# Stop and remove Couchbase Docker container
docker-down:
	@echo "Stopping Couchbase Docker container..."
	docker stop $(DOCKER_COUCHBASE_NAME) || true
	@echo "Removing Couchbase Docker container..."
	docker rm $(DOCKER_COUCHBASE_NAME) || true
	@echo "Couchbase container removed."

# Generate Swagger docs
swagger:
	@echo "Generating Swagger documentation..."
	swag init --parseDependency --parseInternal
	@echo "Swagger docs generated."

# Generate mocks using Mockery
generate-mocks:
	@echo "Generating mocks with Mockery..."
	mockery --all --output=repository/mocks --outpkg=mocks --case=underscore
	@echo "Mocks generated."

# Run unit tests
test: generate-mocks
	@echo "Running unit tests..."
	go test ./tests/... -v
	@echo "Unit tests completed."

# Run the application
run:
	@echo "Starting the application..."
	go run cmd/main.go

# Clean generated files and Docker containers
clean:
	@echo "Cleaning up generated files and Docker containers..."
	make docker-down
	rm -rf docs/
	rm -rf repository/mocks/
	rm -f docs.go
	@echo "Cleaned."
```

**Explanation of Makefile Targets:**

- **`all`**: Default target that starts Couchbase and runs the application.
- **`docker-up`**: Launches a Couchbase server in a Docker container.
- **`docker-down`**: Stops and removes the Couchbase Docker container.
- **`swagger`**: Generates Swagger documentation using `swag`.
- **`generate-mocks`**: Generates mocks using Mockery.
- **`test`**: Runs all unit tests after generating mocks.
- **`run`**: Starts the Go application.
- **`clean`**: Stops Couchbase, removes the Docker container, and cleans generated docs and mocks.

### c. Explanation of Makefile Targets

- **`all`**: Executes `docker-up` and then `run` to start Couchbase and the application.
- **`docker-up`**:
  - Pulls and runs the Couchbase Docker image.
  - Maps necessary ports:
    - `8091-8094`: Couchbase management and data ports.
    - `11210`: Data access port.
  - Sets environment variables for Couchbase admin credentials.
  - Waits for 20 seconds to allow Couchbase to initialize.
- **`docker-down`**:
  - Stops the Couchbase Docker container.
  - Removes the Couchbase Docker container.
- **`swagger`**:
  - Generates Swagger documentation by parsing code annotations.
- **`generate-mocks`**:
  - Uses Mockery to generate mocks for all interfaces.
- **`test`**:
  - Runs `generate-mocks` to ensure mocks are up-to-date.
  - Executes all unit tests in the `tests/` directory.
- **`run`**:
  - Runs the Go application.
- **`clean`**:
  - Executes `docker-down` to stop and remove Couchbase.
  - Removes the `docs/` directory containing Swagger docs.
  - Removes the `repository/mocks/` directory containing generated mocks.
  - Removes the `docs.go` file.

### d. Usage Examples

- **Start Couchbase and Run the Application**

    ```bash
    make all
    ```

- **Generate Swagger Documentation Only**

    ```bash
    make swagger
    ```

- **Generate Mocks and Run Unit Tests**

    ```bash
    make test
    ```

- **Run Unit Tests Without Generating Mocks**

    *Note:* The current `test` target depends on `generate-mocks`. If you want to run tests without regenerating mocks, you can modify the Makefile or run the tests directly.

    ```bash
    go test ./tests/... -v
    ```

- **Stop and Remove Couchbase Docker Container**

    ```bash
    make docker-down
    ```

- **Clean All Generated Files and Containers**

    ```bash
    make clean
    ```

---

## 5. Running the Application

With everything set up, you can now run your application.

### 1. Start Couchbase and the API Server

You can use the Makefile to start Couchbase and run the application simultaneously.

```bash
make all
```

This command will:

- Launch Couchbase in a Docker container.
- Wait for Couchbase to initialize.
- Start the Gin-based API server.

**Alternatively**, you can run the steps separately:

1. **Start Couchbase**

    ```bash
    make docker-up
    ```

2. **Run the Application**

    ```bash
    make run
    ```

### 2. Access the API

- **Get User by ID**

    Use `curl` or any API client (like Postman) to retrieve a user.

    ```bash
    curl http://localhost:8080/users/user1
    ```

    **Expected Response:**

    ```json
    {
      "id": "user1",
      "name": "John Doe",
      "email": "john.doe@example.com"
    }
    ```

- **Swagger UI**

    Open your browser and navigate to [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) to view the Swagger UI and interact with your API documentation.

### 3. Run Unit Tests

Execute unit tests using the Makefile:

```bash
make test
```

This command will:

- Generate mocks using Mockery.
- Run all unit tests in the `tests/` directory.

**Alternatively**, run tests directly:

```bash
go test ./tests/... -v
```

### 4. Generate Swagger Docs

To regenerate Swagger documentation (e.g., after updating handler annotations):

```bash
make swagger
```

### 5. Stop Couchbase

To stop and remove the Couchbase Docker container:

```bash
make docker-down
```

### 6. Clean All Generated Files and Containers

To clean up all generated files and Docker containers:

```bash
make clean
```

---

## 6. Additional Tips

- **Environment Variables**: For better security and flexibility, consider using environment variables to manage configurations such as Couchbase credentials and ports. You can use packages like [`godotenv`](https://github.com/joho/godotenv) to load environment variables from a `.env` file.

- **Error Handling**: Enhance error handling in your handlers to cover more cases and provide more informative responses.

- **Logging**: Implement structured logging using packages like [`logrus`](https://github.com/sirupsen/logrus) or [`zap`](https://github.com/uber-go/zap) for better observability.

- **Configuration Management**: Use configuration files or libraries like [`viper`](https://github.com/spf13/viper) to manage application configurations.

- **Docker Compose**: For more complex setups or additional services, consider using Docker Compose to manage multiple containers. This can help in orchestrating Couchbase alongside other dependencies.

- **CI/CD Integration**: Integrate your tests and documentation generation into a CI/CD pipeline (e.g., GitHub Actions, GitLab CI) to automate testing and deployment.

- **API Versioning**: Implement API versioning to manage changes and maintain backward compatibility.

- **Security**: Implement security best practices, such as input validation, rate limiting, and authentication/authorization mechanisms.

- **Database Migrations**: While Couchbase is schemaless, consider implementing scripts or tools to manage data migrations or seed data.

- **Containerization of the Application**: Consider containerizing your Go application using Docker for easier deployment and scalability.

---

By following the above steps, you will have a functional REST API in Golang using the Gin framework that interacts with Couchbase, complete with unit tests using Mockery, Swagger documentation, and automation through a Makefile. This setup promotes good development practices, making your API scalable, maintainable, and well-documented.



docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' 89e9eb381e31 | grep Hostname
