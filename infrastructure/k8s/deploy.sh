#!/bin/bash

set -e

echo "=========================================="
echo "Deploying Bookstore to Kubernetes"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if minikube is running
print_info "Checking if minikube is running..."
if ! minikube status &> /dev/null; then
    print_warning "Minikube is not running. Starting minikube..."
    minikube start
else
    print_info "Minikube is already running"
fi

# Configure docker environment to use minikube's docker daemon
print_info "Configuring Docker to use minikube's daemon..."
eval $(minikube docker-env)

# Build Docker images
print_info "Building Docker images..."

print_info "Building Catalog Service..."
docker build -t catalog-service:latest ../../services/catalog-service/

print_info "Building API Gateway..."
docker build -t api-gateway:latest ../../services/api-gateway/

print_info "Building Frontend..."
docker build -t frontend:latest ../../frontend/customer-app/

# Deploy to Kubernetes
print_info "Deploying to Kubernetes..."

# Create namespace
print_info "Creating namespace..."
kubectl apply -f namespaces/development.yaml

# Wait a moment for namespace to be ready
sleep 2

# Create secrets and configmaps
print_info "Creating secrets and configmaps..."
kubectl apply -f secrets/
kubectl apply -f configmaps/

# Deploy database
print_info "Deploying PostgreSQL..."
kubectl apply -f databases/

# Wait for PostgreSQL to be ready
print_info "Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres -n bookstore-dev --timeout=300s

# Deploy services
print_info "Deploying Catalog Service..."
kubectl apply -f services/catalog-service/

# Wait for Catalog Service to be ready
print_info "Waiting for Catalog Service to be ready..."
kubectl wait --for=condition=ready pod -l app=catalog-service -n bookstore-dev --timeout=300s

print_info "Deploying API Gateway..."
kubectl apply -f services/api-gateway/

# Wait for API Gateway to be ready
print_info "Waiting for API Gateway to be ready..."
kubectl wait --for=condition=ready pod -l app=api-gateway -n bookstore-dev --timeout=300s

# Deploy frontend
print_info "Deploying Frontend..."
kubectl apply -f frontend/

# Wait for Frontend to be ready
print_info "Waiting for Frontend to be ready..."
kubectl wait --for=condition=ready pod -l app=frontend -n bookstore-dev --timeout=300s

echo ""
echo "=========================================="
print_info "Deployment completed successfully!"
echo "=========================================="
echo ""

# Display access information
print_info "Access Information:"
echo ""
MINIKUBE_IP=$(minikube ip)
echo "  Frontend:    http://${MINIKUBE_IP}:30000"
echo "  API Gateway: http://${MINIKUBE_IP}:30080"
echo ""
print_info "To view pod status, run:"
echo "  kubectl get pods -n bookstore-dev"
echo ""
print_info "To view logs, run:"
echo "  kubectl logs -f <pod-name> -n bookstore-dev"
echo ""
print_info "To view services, run:"
echo "  kubectl get services -n bookstore-dev"
echo ""
