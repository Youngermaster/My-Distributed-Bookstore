#!/bin/bash

set -euo pipefail

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_info() {
    echo -e "${GREEN}[SMOKE]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[SMOKE]${NC} $1"
}

print_error() {
    echo -e "${RED}[SMOKE]${NC} $1"
}

if ! command -v curl >/dev/null 2>&1; then
    print_error "curl is required to run the smoke tests."
    exit 1
fi

if ! command -v python3 >/dev/null 2>&1; then
    print_error "python3 is required to parse JSON responses."
    exit 1
fi

BASE_URL="${1:-}"

if [[ -z "${BASE_URL}" ]]; then
    if ! command -v minikube >/dev/null 2>&1; then
        print_error "Specify the gateway base URL as an argument when not running in Minikube."
        exit 1
    fi
    MINIKUBE_IP=$(minikube ip)
    BASE_URL="http://${MINIKUBE_IP}:30080"
    print_info "Detected Minikube endpoint: ${BASE_URL}"
else
    print_info "Using provided gateway endpoint: ${BASE_URL}"
fi

LAST_RESPONSE=""
http_request() {
    local description="$1"
    local expected_code="$2"
    shift 2

    local tmp
    tmp=$(mktemp)

    set +e
    local http_code
    http_code=$(curl -sS -o "$tmp" -w "%{http_code}" "$@")
    local exit_code=$?
    set -e

    LAST_RESPONSE=$(cat "$tmp")
    rm -f "$tmp"

    if [[ $exit_code -ne 0 ]]; then
        print_error "${description} failed (curl exit code $exit_code)"
        echo "$LAST_RESPONSE"
        exit 1
    fi

    if [[ "$http_code" != "$expected_code" ]]; then
        print_error "${description} failed (expected HTTP $expected_code, got $http_code)"
        echo "$LAST_RESPONSE"
        exit 1
    fi

    print_info "${description} (HTTP ${http_code})"
}

print_info "Running API Gateway smoke tests..."

http_request "Gateway health" 200 "${BASE_URL}/health"
http_request "Gateway readiness" 200 "${BASE_URL}/ready"
http_request "Catalog books listing" 200 "${BASE_URL}/api/v1/catalog/books"

TIMESTAMP=$(date +%s)
TEST_EMAIL="integration.${TIMESTAMP}@example.com"
TEST_PASSWORD="TestPass123!"

REGISTER_PAYLOAD=$(cat <<EOF
{
  "email": "${TEST_EMAIL}",
  "password": "${TEST_PASSWORD}",
  "full_name": "Integration User"
}
EOF
)

http_request "User registration via gateway" 201 \
    -X POST \
    -H "Content-Type: application/json" \
    -d "${REGISTER_PAYLOAD}" \
    "${BASE_URL}/api/v1/users/auth/register"

USER_ID=$(LAST_RESPONSE="$LAST_RESPONSE" python3 - <<'PY'
import json
import os

payload = os.environ.get("LAST_RESPONSE", "")
try:
    data = json.loads(payload)
    print(str(data["id"]))
except Exception as exc:
    raise SystemExit(f"failed to parse user registration response: {exc}")
PY
)
USER_ID=${USER_ID//$'\n'/}

if [[ -z "${USER_ID}" ]]; then
    print_error "Unable to extract user ID from registration response."
    echo "$LAST_RESPONSE"
    exit 1
fi

CART_ID=$(python3 - <<'PY'
import uuid
print(uuid.uuid4())
PY
)
CART_ID=${CART_ID//$'\n'/}

BOOK_ID=$(python3 - <<'PY'
import uuid
print(uuid.uuid4())
PY
)
BOOK_ID=${BOOK_ID//$'\n'/}

CART_PAYLOAD=$(cat <<EOF
{
  "book_id": "${BOOK_ID}",
  "quantity": 1,
  "price": 19.99
}
EOF
)

http_request "Add item to cart via gateway" 201 \
    -X POST \
    -H "Content-Type: application/json" \
    -d "${CART_PAYLOAD}" \
    "${BASE_URL}/api/v1/cart/${CART_ID}/items"

ORDER_PAYLOAD=$(cat <<EOF
{
  "user_id": "${USER_ID}",
  "items": [
    {
      "book_id": "${BOOK_ID}",
      "quantity": 1,
      "unit_price": 19.99
    }
  ],
  "shipping_address": "123 Main St, Test City",
  "payment_method": "credit_card"
}
EOF
)

http_request "Create order via gateway" 201 \
    -X POST \
    -H "Content-Type: application/json" \
    -d "${ORDER_PAYLOAD}" \
    "${BASE_URL}/api/v1/orders"

http_request "Fetch trending recommendations" 200 \
    "${BASE_URL}/api/v1/recommendations/trending"

print_info "Smoke tests completed successfully."
