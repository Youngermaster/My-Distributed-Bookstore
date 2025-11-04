#!/bin/bash

set -euo pipefail

echo "=========================================="
echo "Deploying Distributed Bookstore (Dev)"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging helpers
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Ensure Minikube is running
print_info "Checking Minikube status..."
if ! minikube status &> /dev/null; then
    print_warning "Minikube is not running. Starting Minikube..."
    minikube start
else
    print_info "Minikube is already running."
fi

# Point Docker at the Minikube daemon
print_info "Configuring Docker to use Minikube's daemon..."
eval "$(minikube docker-env)"

# Build Docker images for all services
print_info "Building service images (this can take a few minutes)..."
IMAGE_SPECS=(
    "catalog-service|../../services/catalog-service"
    "user-service|../../services/user-service"
    "cart-service|../../services/cart-service"
    "order-service|../../services/order-service"
    "recommendation-service|../../services/recommendation-service"
    "inventory-service|../../services/inventory-service"
    "review-service|../../services/review-service"
    "api-gateway|../../services/api-gateway"
    "frontend|../../frontend/customer-app"
)

for spec in "${IMAGE_SPECS[@]}"; do
    IFS='|' read -r image path <<< "$spec"
    print_info "Building ${image}:latest from ${path}..."
    docker build -t "${image}:latest" "${path}"
done

# Namespace & shared configuration
print_info "Creating/Updating namespace..."
kubectl apply -f namespaces/development.yaml
sleep 2

print_info "Applying global secrets..."
kubectl apply -f secrets/

print_info "Applying global configmaps..."
kubectl apply -f configmaps/

# Data stores
print_info "Deploying PostgreSQL..."
kubectl apply -f databases/

print_info "Waiting for PostgreSQL to become ready..."
kubectl wait --for=condition=ready pod -l app=postgres -n bookstore-dev --timeout=300s

print_info "Deploying Redis for cart/recommendations..."
kubectl apply -f messaging/redis/

print_info "Waiting for Redis to become ready..."
kubectl wait --for=condition=ready pod -l app=redis -n bookstore-dev --timeout=300s

# Helper to deploy services and wait for readiness
deploy_component() {
    local friendly_name="$1"
    local manifest_path="$2"
    local label_selector="$3"

    print_info "Deploying ${friendly_name}..."
    kubectl apply -f "${manifest_path}"

    print_info "Waiting for ${friendly_name} pods to become ready..."
    kubectl wait --for=condition=ready pod -l "${label_selector}" -n bookstore-dev --timeout=300s
}

# Core services
deploy_component "Catalog Service" "services/catalog-service/" "app=catalog-service"
deploy_component "User Service" "services/user-service/" "app=user-service"
deploy_component "Cart Service" "services/cart-service/" "app=cart-service"
deploy_component "Order Service" "services/order-service/" "app=order-service"
deploy_component "Recommendation Service" "services/recommendation-service/" "app=recommendation-service"
deploy_component "Inventory Service" "services/inventory-service/" "app=inventory-service"
deploy_component "Review Service" "services/review-service/" "app=review-service"
deploy_component "API Gateway" "services/api-gateway/" "app=api-gateway"
deploy_component "Frontend" "frontend/" "app=frontend"

echo ""
MINIKUBE_IP=$(minikube ip)
BASE_URL="http://${MINIKUBE_IP}:30080"
echo "=========================================="
print_info "Deployment completed successfully!"
echo "=========================================="
echo ""
print_info "Access Information:"
echo "  Frontend:    http://${MINIKUBE_IP}:30000"
echo "  API Gateway: ${BASE_URL}"
echo ""
print_info "To inspect pods:"
echo "  kubectl get pods -n bookstore-dev"
echo ""
