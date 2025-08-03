# Stage 1: The Build Stage
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker's layer caching.
# This ensures that these dependencies are only downloaded when they change.
COPY go.mod go.sum ./

# Download all the modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application binary statically.
# CGO_ENABLED=0 disables CGO, which means the binary won't link against any C libraries.
# This makes it a statically-linked binary, which can run on a minimal base image like scratch.
# -ldflags="-s -w" removes symbol and debug information, further reducing binary size.
RUN CGO_ENABLED=0 go build -a -ldflags="-s -w" -o /myapp .

# Stage 2: The Production Stage
# Use a minimal base image. 'scratch' is the smallest possible image.
FROM scratch

# Set the working directory
WORKDIR /

# Copy the static binary from the builder stage
COPY --from=builder /myapp /myapp

# Expose the port on which the application will listen
EXPOSE 8080

# The command to run the application
ENTRYPOINT ["/myapp"]