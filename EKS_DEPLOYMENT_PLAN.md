# 🚀 EKS Deployment Plan - URGENT (Due Today)

**Date:** November 4, 2025  
**Objective:** Deploy Distributed Bookstore to AWS EKS in the simplest way possible  
**Time Constraint:** Due today

---

## 📊 Current Status Review

### ✅ What's Ready (100%)

#### Backend Services - Implemented
1. **Catalog Service** ✅ (Go + PostgreSQL)
   - Full CRUD for books, authors, publishers, categories
   - Docker image ready
   - K8s manifests ready
   
2. **User Service** ✅ (Go + PostgreSQL)
   - Authentication (JWT)
   - Registration, login, roles (admin/customer)
   - Wishlist functionality
   - Docker image ready
   
3. **Cart Service** ✅ (Go + Redis)
   - Session-based cart management
   - 5 REST endpoints
   - Docker image ready
   - **K8s manifests MISSING**
   
4. **Order Service** ✅ (Go + PostgreSQL)
   - Complete order lifecycle
   - 6 REST endpoints
   - Docker image ready
   - **K8s manifests MISSING**

5. **API Gateway** ✅ (Go)
   - Routing, CORS, rate limiting
   - Docker image ready
   - K8s manifests ready

#### Frontend
6. **React App** ✅
   - All pages implemented
   - shadcn/ui components
   - Docker image ready
   - K8s manifests ready

#### Infrastructure - Partial
7. **Kubernetes Setup** ⚠️ (70% Ready)
   - ✅ Minikube tested and working
   - ✅ Automated deployment script (deploy.sh)
   - ✅ PostgreSQL StatefulSet
   - ✅ Catalog + API Gateway + Frontend manifests
   - ❌ Cart Service K8s manifests
   - ❌ Order Service K8s manifests
   - ❌ User Service K8s manifests
   - ❌ Redis deployment
   - ❌ EKS-specific configurations

### ❌ What's NOT Ready

#### Missing Services (Not Critical for MVP)
- Review Service (Python/ML)
- Recommendation Service (Python/ML)
- Notification Service (TypeScript)
- Payment Service (TypeScript)
- Inventory Service (Python)
- Admin Service (Go)

#### Missing Infrastructure
- RabbitMQ deployment
- Monitoring (Prometheus/Grafana)
- Logging (ELK)
- Tracing (Jaeger)

---

## 🎯 SIMPLEST EKS Deployment Strategy (TODAY)

### Phase 1: Minimal Viable Deployment (2-3 hours)

Deploy ONLY these 6 services to prove the system works:

1. **PostgreSQL** (AWS RDS or in-cluster)
2. **Redis** (AWS ElastiCache or in-cluster)
3. **Catalog Service**
4. **User Service**
5. **Cart Service**
6. **Order Service**
7. **API Gateway**
8. **Frontend**

---

## 📋 Step-by-Step EKS Deployment (FASTEST PATH)

### Prerequisites (15 mins)

```bash
# 1. Install AWS CLI
choco install awscli  # Windows

# 2. Install eksctl
choco install eksctl

# 3. Install kubectl (if not installed)
choco install kubernetes-cli

# 4. Configure AWS credentials
aws configure
# Enter: Access Key ID, Secret Access Key, Region (us-east-1), Output (json)
```

### Step 1: Create EKS Cluster (20-30 mins) ⏱️

**Option A: Using eksctl (RECOMMENDED - FASTEST)**

```bash
# Create cluster with minimal config
eksctl create cluster \
  --name bookstore-cluster \
  --region us-east-1 \
  --nodegroup-name standard-workers \
  --node-type t3.medium \
  --nodes 2 \
  --nodes-min 1 \
  --nodes-max 3 \
  --managed

# This creates:
# - EKS cluster
# - VPC with subnets
# - Security groups
# - Node group (2 t3.medium instances)
# - Automatically configures kubectl
```

**Verify:**
```bash
kubectl get nodes
# Should show 2 nodes in Ready state
```

### Step 2: Setup Container Registry (10 mins)

```bash
# Create ECR repositories for each service
aws ecr create-repository --repository-name catalog-service --region us-east-1
aws ecr create-repository --repository-name user-service --region us-east-1
aws ecr create-repository --repository-name cart-service --region us-east-1
aws ecr create-repository --repository-name order-service --region us-east-1
aws ecr create-repository --repository-name api-gateway --region us-east-1
aws ecr create-repository --repository-name frontend --region us-east-1

# Login to ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com
```

### Step 3: Build and Push Images (15-20 mins)

```bash
# Navigate to project root
cd "C:\Users\sebas\Documents\Eafit\Semestre 10\Telemática\My-Distributed-Bookstore"

# Set your AWS account ID
$AWS_ACCOUNT_ID = "<your-account-id>"
$AWS_REGION = "us-east-1"
$ECR_REGISTRY = "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com"

# Build and push each service
# Catalog Service
docker build -t catalog-service:latest services/catalog-service/
docker tag catalog-service:latest $ECR_REGISTRY/catalog-service:latest
docker push $ECR_REGISTRY/catalog-service:latest

# User Service
docker build -t user-service:latest services/user-service/
docker tag user-service:latest $ECR_REGISTRY/user-service:latest
docker push $ECR_REGISTRY/user-service:latest

# Cart Service
docker build -t cart-service:latest services/cart-service/
docker tag cart-service:latest $ECR_REGISTRY/cart-service:latest
docker push $ECR_REGISTRY/cart-service:latest

# Order Service
docker build -t order-service:latest services/order-service/
docker tag order-service:latest $ECR_REGISTRY/order-service:latest
docker push $ECR_REGISTRY/order-service:latest

# API Gateway
docker build -t api-gateway:latest services/api-gateway/
docker tag api-gateway:latest $ECR_REGISTRY/api-gateway:latest
docker push $ECR_REGISTRY/api-gateway:latest

# Frontend
docker build -t frontend:latest frontend/customer-app/
docker tag frontend:latest $ECR_REGISTRY/frontend:latest
docker push $ECR_REGISTRY/frontend:latest
```

### Step 4: Create Missing K8s Manifests (30 mins)

**I'll create these for you now - see below**

### Step 5: Deploy to EKS (15 mins)

```bash
cd infrastructure/k8s

# Create namespace
kubectl apply -f namespaces/development.yaml

# Create secrets (update with real values)
kubectl apply -f secrets/

# Create configmaps
kubectl apply -f configmaps/

# Deploy databases
kubectl apply -f databases/

# Wait for databases
kubectl wait --for=condition=ready pod -l app=postgres -n bookstore-dev --timeout=300s

# Deploy services
kubectl apply -f services/catalog-service/
kubectl apply -f services/user-service/
kubectl apply -f services/cart-service/
kubectl apply -f services/order-service/
kubectl apply -f services/api-gateway/
kubectl apply -f frontend/

# Check deployment status
kubectl get pods -n bookstore-dev
```

### Step 6: Expose Services (10 mins)

**Option A: LoadBalancer (SIMPLEST)**

```bash
# Create LoadBalancer for API Gateway
kubectl expose deployment api-gateway \
  --type=LoadBalancer \
  --port=80 \
  --target-port=8080 \
  --name=api-gateway-lb \
  -n bookstore-dev

# Create LoadBalancer for Frontend
kubectl expose deployment frontend \
  --type=LoadBalancer \
  --port=80 \
  --target-port=80 \
  --name=frontend-lb \
  -n bookstore-dev

# Get URLs (wait ~2 mins for AWS to provision)
kubectl get svc -n bookstore-dev
```

**Option B: Ingress (More Professional)**

See separate section below if time permits.

---

## 🛠️ Missing K8s Manifests to Create

### 1. Redis Deployment

**File:** `infrastructure/k8s/databases/redis-deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
  namespace: bookstore-dev
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - name: redis
        image: redis:7-alpine
        ports:
        - containerPort: 6379
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: redis
  namespace: bookstore-dev
spec:
  type: ClusterIP
  ports:
  - port: 6379
    targetPort: 6379
  selector:
    app: redis
```

### 2. User Service Deployment

**File:** `infrastructure/k8s/services/user-service/deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  namespace: bookstore-dev
spec:
  replicas: 2
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/user-service:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8081
        env:
        - name: DB_HOST
          value: "postgres"
        - name: DB_PORT
          value: "5432"
        - name: DB_USER
          value: "userservice"
        - name: DB_PASSWORD
          value: "userpass123"  # Use secret in production
        - name: DB_NAME
          value: "userdb"
        - name: PORT
          value: "8081"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: jwt-secret
              key: secret
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
  namespace: bookstore-dev
spec:
  type: ClusterIP
  ports:
  - port: 8081
    targetPort: 8081
  selector:
    app: user-service
```

### 3. Cart Service Deployment

**File:** `infrastructure/k8s/services/cart-service/deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cart-service
  namespace: bookstore-dev
spec:
  replicas: 2
  selector:
    matchLabels:
      app: cart-service
  template:
    metadata:
      labels:
        app: cart-service
    spec:
      containers:
      - name: cart-service
        image: <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/cart-service:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8083
        env:
        - name: REDIS_HOST
          value: "redis"
        - name: REDIS_PORT
          value: "6379"
        - name: REDIS_PASSWORD
          value: ""
        - name: REDIS_DB
          value: "0"
        - name: PORT
          value: "8083"
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8083
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8083
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: cart-service
  namespace: bookstore-dev
spec:
  type: ClusterIP
  ports:
  - port: 8083
    targetPort: 8083
  selector:
    app: cart-service
```

### 4. Order Service Deployment

**File:** `infrastructure/k8s/services/order-service/deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-service
  namespace: bookstore-dev
spec:
  replicas: 2
  selector:
    matchLabels:
      app: order-service
  template:
    metadata:
      labels:
        app: order-service
    spec:
      containers:
      - name: order-service
        image: <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/order-service:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8084
        env:
        - name: DB_HOST
          value: "postgres"
        - name: DB_PORT
          value: "5432"
        - name: DB_USER
          value: "orderuser"
        - name: DB_PASSWORD
          value: "orderpass123"  # Use secret in production
        - name: DB_NAME
          value: "orderdb"
        - name: PORT
          value: "8084"
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8084
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8084
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: order-service
  namespace: bookstore-dev
spec:
  type: ClusterIP
  ports:
  - port: 8084
    targetPort: 8084
  selector:
    app: order-service
```

---

## ⚡ Quick Deployment Script (All-in-One)

**File:** `infrastructure/k8s/eks-deploy.sh`

```bash
#!/bin/bash
set -e

echo "🚀 Starting EKS Deployment..."

# Variables (UPDATE THESE)
AWS_ACCOUNT_ID="your-account-id"
AWS_REGION="us-east-1"
ECR_REGISTRY="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"

# Login to ECR
echo "📦 Logging into ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_REGISTRY

# Build and push images
echo "🔨 Building and pushing images..."
services=("catalog-service" "user-service" "cart-service" "order-service" "api-gateway")

for service in "${services[@]}"; do
    echo "Building $service..."
    docker build -t $service:latest services/$service/
    docker tag $service:latest $ECR_REGISTRY/$service:latest
    docker push $ECR_REGISTRY/$service:latest
done

# Build frontend
echo "Building frontend..."
docker build -t frontend:latest frontend/customer-app/
docker tag frontend:latest $ECR_REGISTRY/frontend:latest
docker push $ECR_REGISTRY/frontend:latest

# Deploy to K8s
echo "☸️ Deploying to Kubernetes..."
kubectl apply -f namespaces/development.yaml
kubectl apply -f secrets/
kubectl apply -f configmaps/
kubectl apply -f databases/

# Wait for databases
echo "⏳ Waiting for databases..."
kubectl wait --for=condition=ready pod -l app=postgres -n bookstore-dev --timeout=300s

# Deploy services
kubectl apply -f services/catalog-service/
kubectl apply -f services/user-service/
kubectl apply -f services/cart-service/
kubectl apply -f services/order-service/
kubectl apply -f services/api-gateway/
kubectl apply -f frontend/

echo "✅ Deployment complete!"
echo ""
echo "Check status with:"
echo "  kubectl get pods -n bookstore-dev"
echo ""
echo "Get service URLs with:"
echo "  kubectl get svc -n bookstore-dev"
```

---

## 🎯 Expected Timeline

| Task | Time | Status |
|------|------|--------|
| Setup AWS CLI & eksctl | 15 min | ⏳ |
| Create EKS cluster | 30 min | ⏳ |
| Create ECR repos | 10 min | ⏳ |
| Build & push images | 20 min | ⏳ |
| Create missing K8s manifests | 30 min | ⏳ |
| Deploy to EKS | 15 min | ⏳ |
| Expose services | 10 min | ⏳ |
| Testing & debugging | 30 min | ⏳ |
| **TOTAL** | **2.5-3 hrs** | |

---

## 🆘 Emergency Fallback (If EKS Fails)

If you run out of time for EKS:

### Option 1: Docker Compose (5 mins)
```bash
docker-compose -f docker-compose.yml up -d
```
Accessible at `http://localhost:8080`

### Option 2: Minikube (10 mins)
```bash
cd infrastructure/k8s
./deploy.sh
```
Accessible at `http://$(minikube ip):30000`

---

## 📝 What to Present

### Minimum Viable Demo:
1. **Show running pods** in EKS
2. **API Gateway health check** working
3. **Frontend** loading
4. **One complete flow:** Register → Login → Browse Books → Add to Cart → Create Order

### Architecture Diagram (Already have):
- Show microservices architecture
- Highlight what's deployed vs. what's planned

### Mention Future Work:
- ML services (recommendation, review sentiment)
- Monitoring (Prometheus/Grafana)
- CI/CD pipeline
- Production secrets management

---

## 🚨 Critical Notes

1. **Replace `<AWS_ACCOUNT_ID>`** in all manifests with your actual AWS account ID
2. **Use proper secrets** (not hardcoded passwords) if you have time
3. **Monitor costs** - remember to delete the cluster after demo:
   ```bash
   eksctl delete cluster --name bookstore-cluster --region us-east-1
   ```

---

## Next Steps After Today

1. Add Ingress with SSL
2. Setup RDS instead of in-cluster PostgreSQL
3. Setup ElastiCache instead of in-cluster Redis
4. Add monitoring (Prometheus/Grafana)
5. Setup CI/CD with GitHub Actions
6. Implement auto-scaling (HPA)

---

**PRIORITY:** Start with eksctl cluster creation NOW - it takes 20-30 minutes!

Good luck! 🍀
