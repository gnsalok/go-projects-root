# Makefile for Go Couchbase API

.PHONY: all run test swagger docker-up docker-down clean

APP_NAME=go-couchbase-api
DOCKER_COUCHBASE_NAME=couchbase-local
COUCHBASE_IMAGE=couchbase
COUCHBASE_VERSION=latest
COUCHBASE_ADMIN=admin
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
	sleep 15
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
	swag init --dir=cmd --output=./docs --parseDependency --parseInternal
	@echo "Swagger docs generated."

# Run unit tests
test:
	@echo "Running unit tests..."
	go test ./... -v
	@echo "Unit tests completed."

# Run the application
run:
	@echo "Starting the application..."
	go run main.go

# Clean generated files
clean:
	@echo "Cleaning up generated files..."
	docker-down
	rm -rf docs/
	rm -f docs.go
	@echo "Cleaned."
