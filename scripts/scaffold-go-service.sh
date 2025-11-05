#!/bin/bash

# Script to scaffold a Go microservice
# Usage: ./scaffold-go-service.sh <service-name> <port-http> <port-grpc>

SERVICE_NAME=$1
PORT_HTTP=$2
PORT_GRPC=$3

if [ -z "$SERVICE_NAME" ] || [ -z "$PORT_HTTP" ] || [ -z "$PORT_GRPC" ]; then
    echo "Usage: ./scaffold-go-service.sh <service-name> <port-http> <port-grpc>"
    exit 1
fi

SERVICE_DIR="services/$SERVICE_NAME"

# Create directory structure
mkdir -p $SERVICE_DIR/{cmd/server,internal/{config,domain,repository/postgres,service,handler/{http,grpc},middleware,events,tracing},pkg/{jwt,validator,errors,response},proto,migrations}

# Create go.mod
cat > $SERVICE_DIR/go.mod <<EOF
module github.com/youngermaster/distributed-bookstore/$SERVICE_NAME

go 1.21

require (
	github.com/gofiber/fiber/v2 v2.51.0
	gorm.io/gorm v1.25.5
	gorm.io/driver/postgres v1.5.4
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/redis/go-redis/v9 v9.3.0
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
	github.com/streadway/amqp v1.1.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/prometheus/client_golang v1.17.0
)
EOF

# Create Dockerfile
cat > $SERVICE_DIR/Dockerfile <<EOF
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN if [ -f go.sum ]; then go mod download; fi

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Final stage
FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Copy migrations (if needed)
COPY --from=builder /app/migrations ./migrations

# Expose ports
EXPOSE $PORT_HTTP $PORT_GRPC

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \\
  CMD wget --no-verbose --tries=1 --spider http://localhost:$PORT_HTTP/health || exit 1

# Run
CMD ["./main"]
EOF

# Create placeholder main.go
cat > $SERVICE_DIR/cmd/server/main.go <<EOF
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "$SERVICE_NAME",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "$SERVICE_NAME",
		})
	})

	// TODO: Add database connection
	// TODO: Add gRPC server
	// TODO: Add route handlers
	// TODO: Add RabbitMQ connection
	// TODO: Add distributed tracing
	// TODO: Add metrics

	// Get port from environment or use default
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "$PORT_HTTP"
	}

	// Start server in goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("$SERVICE_NAME started on port %s", port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down $SERVICE_NAME...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("$SERVICE_NAME stopped")
}
EOF

chmod +x $SERVICE_DIR/cmd/server/main.go

echo "Go service $SERVICE_NAME scaffolded successfully!"
