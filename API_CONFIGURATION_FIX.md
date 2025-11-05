# API Configuration Fix - Frontend to Backend Communication

## Problem Identified
The frontend was configured to connect to `http://localhost:8080` which doesn't work when running in Kubernetes. Browser requests from users' machines need to connect to services within the cluster.

## Solution Implemented

### 1. Frontend API Configuration
**Changed:** `frontend-config.yaml`
- **Before:** `VITE_API_URL: "http://localhost:30080"`
- **After:** `VITE_API_URL: ""`

Empty string makes the frontend use relative URLs (same domain).

### 2. Nginx Proxy Configuration
The frontend's nginx configuration already includes API proxying:
```nginx
location /api/ {
    proxy_pass http://api-gateway.bookstore-dev.svc.cluster.local:8080;
    # ... proxy headers and CORS
}
```

This means:
- Browser requests to `http://<frontend-lb>/api/v1/catalog/books`
- Nginx proxies internally to `http://api-gateway:8080/api/v1/catalog/books`
- No CORS issues, no external API Gateway exposure needed

### 3. Frontend Image Rebuilt
```bash
# Rebuild with empty API URL
docker build --build-arg VITE_API_URL="" -t 905418472239.dkr.ecr.us-east-1.amazonaws.com/customer-frontend:latest .

# Push to ECR
docker push 905418472239.dkr.ecr.us-east-1.amazonaws.com/customer-frontend:latest

# Restart deployment
kubectl rollout restart deployment/frontend -n bookstore-dev
```

## How It Works Now

### Request Flow:
```
User Browser
    ↓
Frontend LoadBalancer (http://ab1a1c3c5b1ca49768c26f26e92ca780-844836377.us-east-1.elb.amazonaws.com)
    ↓
Frontend Pod (Nginx)
    ↓
/api/* requests → Nginx Proxy → API Gateway Service (internal)
    ↓
API Gateway → Backend Services (Catalog, User, Cart, Order, etc.)
```

### API Request Example:
```javascript
// Frontend code
booksAPI.list() // GET /api/v1/catalog/books

// Browser makes request to:
http://ab1a1c3c5b1ca49768c26f26e92ca780-844836377.us-east-1.elb.amazonaws.com/api/v1/catalog/books

// Nginx proxies internally to:
http://api-gateway.bookstore-dev.svc.cluster.local:8080/api/v1/catalog/books
```

## Benefits of This Approach

✅ **No CORS issues** - Same-origin requests (frontend and API on same domain)
✅ **Secure** - API Gateway not exposed directly to internet
✅ **Simple** - One LoadBalancer instead of two
✅ **Production-ready** - Standard pattern for SPAs with backends

## Services Still Need Fixing

### 1. user-service (CreateContainerConfigError)
- Status: Pods cannot start
- Likely Issue: Missing or incorrect secret keys
- Next Step: Check pod events and secrets

### 2. order-service (CrashLoopBackOff)
- Status: Pods crashing on startup
- Likely Issue: Database connection error
- Next Step: Check pod logs

## Testing the Application

### 1. Access Frontend
```
URL: http://ab1a1c3c5b1ca49768c26f26e92ca780-844836377.us-east-1.elb.amazonaws.com
```

### 2. Test Catalog API (should work now)
Open browser console and check:
- Network tab → `/api/v1/catalog/books` requests
- Should return book data (not ERR_CONNECTION_REFUSED)

### 3. What Should Work
✅ Browse books
✅ View categories
✅ Search books
✅ View book details

### 4. What Won't Work Yet
❌ User registration/login (user-service down)
❌ Shopping cart (needs authentication)
❌ Orders (order-service down)
❌ Wishlist (user-service down)

## Next Steps

1. **Fix user-service**
   ```bash
   kubectl describe pod <user-service-pod> -n bookstore-dev
   kubectl get secret user-service-secrets -n bookstore-dev -o yaml
   ```

2. **Fix order-service**
   ```bash
   kubectl logs <order-service-pod> -n bookstore-dev
   ```

3. **Deploy new services** (when core is stable)
   - admin-service
   - inventory-service
   - review-service
   - recommendation-service

## Files Modified

1. `infrastructure/k8s/configmaps/frontend-config.yaml` - Empty API URL
2. `infrastructure/k8s/services/api-gateway/service.yaml` - Changed to LoadBalancer (reverted - not needed)
3. `frontend/customer-app/` - Rebuilt Docker image with new config

## Deployment Status

| Service | Status | Replicas | Notes |
|---------|--------|----------|-------|
| catalog-service | ✅ Running | 2/2 | Working |
| api-gateway | ✅ Running | 2/2 | Working |
| cart-service | ✅ Running | 2/2 | Working |
| frontend | ✅ Running | 2/2 | **Updated with new config** |
| postgres | ✅ Running | 1/1 | Working |
| redis | ✅ Running | 1/1 | Working |
| user-service | ❌ Error | 0/2 | CreateContainerConfigError |
| order-service | ❌ Error | 0/2 | CrashLoopBackOff (7 restarts) |

**Overall Status:** 6/8 services running (75%), API configuration fixed, catalog browsing should now work in browser!
