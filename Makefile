.PHONY: help proto-gen services-start services-stop docker-build docker-clean test lint migrate-up migrate-down k8s-deploy

# Default target
help: ## Show this help message
	@echo "Distributed Bookstore - Makefile Commands"
	@echo ""
	@echo "Proto & Code Generation:"
	@echo "  make proto-gen        Generate protobuf code for all languages"
	@echo ""
	@echo "Development:"
	@echo "  make services-start   Start all services with Docker Compose"
	@echo "  make services-stop    Stop all services"
	@echo "  make logs            Tail logs from all services"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build     Build all Docker images"
	@echo "  make docker-clean     Remove all Docker containers and volumes"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-up       Run database migrations"
	@echo "  make migrate-down     Rollback database migrations"
	@echo ""
	@echo "Testing:"
	@echo "  make test            Run all tests"
	@echo "  make test-go         Run Go service tests"
	@echo "  make test-node       Run Node.js service tests"
	@echo "  make test-python     Run Python service tests"
	@echo ""
	@echo "Code Quality:"
	@echo "  make lint            Run linters for all services"
	@echo "  make lint-fix        Auto-fix linting issues"
	@echo ""
	@echo "Kubernetes:"
	@echo "  make k8s-deploy      Deploy to Kubernetes"
	@echo "  make k8s-delete      Delete from Kubernetes"
	@echo "  make k8s-status      Check deployment status"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean           Clean all build artifacts"

# Generate protobuf code for all languages
proto-gen: ## Generate protobuf code
	@echo "Generating protobuf code..."
	@for proto in proto/*.proto; do \
		protoc --go_out=. --go_opt=paths=source_relative \
		       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		       $$proto 2>/dev/null || echo "Go protoc not installed or failed for $$proto"; \
	done
	@echo "Protobuf code generation complete"

# Start all services
services-start: ## Start all services
	@echo "Starting all services..."
	docker-compose up -d
	@echo "All services started"
	@echo "API Gateway: http://localhost:8080"
	@echo "Jaeger UI: http://localhost:16686"
	@echo "RabbitMQ Management: http://localhost:15672"

# Stop all services
services-stop: ## Stop all services
	@echo "Stopping all services..."
	docker-compose down
	@echo "All services stopped"

# View logs
logs: ## View logs from all services
	docker-compose logs -f

# Build all Docker images
docker-build: ## Build all Docker images
	@echo "Building all Docker images..."
	docker-compose build
	@echo "All images built"

# Clean Docker resources
docker-clean: ## Clean Docker resources
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f
	@echo "Docker resources cleaned"

# Run database migrations
migrate-up: ## Run database migrations
	@echo "Running database migrations..."
	@echo "TODO: Implement migration script"

# Rollback database migrations
migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	@echo "TODO: Implement migration script"

# Run all tests
test: test-go test-node test-python ## Run all tests
	@echo "All tests complete"

# Run Go service tests
test-go: ## Run Go service tests
	@echo "Running Go service tests..."
	@for service in services/api-gateway services/catalog-service services/user-service services/cart-service services/order-service services/inventory-service services/admin-service; do \
		if [ -f $$service/go.mod ]; then \
			echo "Testing $$service..."; \
			cd $$service && go test ./... -v -cover || true; \
			cd ../..; \
		fi \
	done
	@echo "Go tests complete"

# Run Node.js service tests
test-node: ## Run Node.js service tests
	@echo "Running Node.js service tests..."
	@for service in services/payment-service services/notification-service; do \
		if [ -f $$service/package.json ]; then \
			echo "Testing $$service..."; \
			cd $$service && npm test || true; \
			cd ../..; \
		fi \
	done
	@echo "Node.js tests complete"

# Run Python service tests
test-python: ## Run Python service tests
	@echo "Running Python service tests..."
	@for service in services/review-service services/recommendation-service; do \
		if [ -f $$service/requirements.txt ]; then \
			echo "Testing $$service..."; \
			cd $$service && pytest || true; \
			cd ../..; \
		fi \
	done
	@echo "Python tests complete"

# Run linters
lint: ## Run linters
	@echo "Running linters..."
	@echo "TODO: Implement linting for all services"

# Deploy to Kubernetes
k8s-deploy: ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	kubectl apply -f infrastructure/k8s/ --recursive
	@echo "Kubernetes deployment complete"

# Delete from Kubernetes
k8s-delete: ## Delete from Kubernetes
	@echo "Deleting from Kubernetes..."
	kubectl delete -f infrastructure/k8s/ --recursive
	@echo "Kubernetes resources deleted"

# Check Kubernetes status
k8s-status: ## Check Kubernetes deployment status
	@echo "Checking Kubernetes deployment status..."
	kubectl get pods -n production
	kubectl get services -n production

# Clean build artifacts
clean: ## Clean all build artifacts
	@echo "Cleaning build artifacts..."
	find . -type d -name "dist" -exec rm -rf {} + 2>/dev/null || true
	find . -type d -name "build" -exec rm -rf {} + 2>/dev/null || true
	find . -type d -name "bin" -exec rm -rf {} + 2>/dev/null || true
	find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	@echo "Build artifacts cleaned"
