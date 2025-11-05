#!/bin/bash

set -euo pipefail

echo "=========================================="
echo "Undeploying Distributed Bookstore (Dev)"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_info "Tearing down workloads from namespace bookstore-dev..."

# Remove application workloads in reverse order of deployment
kubectl delete -f frontend/ --ignore-not-found=true
kubectl delete -f services/api-gateway/ --ignore-not-found=true
kubectl delete -f services/recommendation-service/ --ignore-not-found=true
kubectl delete -f services/review-service/ --ignore-not-found=true
kubectl delete -f services/inventory-service/ --ignore-not-found=true
kubectl delete -f services/order-service/ --ignore-not-found=true
kubectl delete -f services/cart-service/ --ignore-not-found=true
kubectl delete -f services/user-service/ --ignore-not-found=true
kubectl delete -f services/catalog-service/ --ignore-not-found=true
kubectl delete -f services/notification-service/ --ignore-not-found=true

# Infrastructure dependencies
kubectl delete -f messaging/redis/ --ignore-not-found=true
kubectl delete -f messaging/rabbitmq/ --ignore-not-found=true
kubectl delete -f databases/ --ignore-not-found=true

# Shared configuration
kubectl delete -f configmaps/ --ignore-not-found=true
kubectl delete -f secrets/ --ignore-not-found=true

# Persistent storage
print_info "Deleting PersistentVolumeClaims..."
kubectl delete pvc -l app=postgres -n bookstore-dev --ignore-not-found=true

print_warning "Namespace 'bookstore-dev' was left in place. Remove it manually if needed:"
echo "  kubectl delete namespace bookstore-dev"

echo ""
print_info "Undeployment complete."
echo ""
