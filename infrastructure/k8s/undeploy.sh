#!/bin/bash

set -e

echo "=========================================="
echo "Undeploying Bookstore from Kubernetes"
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

# Delete all resources
print_info "Deleting all resources from bookstore-dev namespace..."

kubectl delete -f frontend/ --ignore-not-found=true
kubectl delete -f services/api-gateway/ --ignore-not-found=true
kubectl delete -f services/catalog-service/ --ignore-not-found=true
kubectl delete -f databases/ --ignore-not-found=true
kubectl delete -f configmaps/ --ignore-not-found=true
kubectl delete -f secrets/ --ignore-not-found=true

# Delete PVCs
print_info "Deleting PersistentVolumeClaims..."
kubectl delete pvc -l app=postgres -n bookstore-dev --ignore-not-found=true

# Optionally delete namespace (uncomment if you want to delete the namespace)
# print_warning "Deleting namespace..."
# kubectl delete -f namespaces/development.yaml --ignore-not-found=true

echo ""
print_info "Undeployment completed!"
echo ""
print_info "To delete the namespace as well, run:"
echo "  kubectl delete namespace bookstore-dev"
echo ""
