# Complete EKS Deployment Guide
## Distributed Bookstore - AWS Academy Edition

This guide walks through the complete deployment process for deploying the Distributed Bookstore application to AWS EKS using AWS Academy accounts.

---

## Table of Contents
- [Prerequisites](#prerequisites)
- [Architecture Overview](#architecture-overview)
- [AWS Academy Setup](#aws-academy-setup)
- [Install Required Tools](#install-required-tools)
- [Create EKS Cluster](#create-eks-cluster)
- [Build and Push Docker Images](#build-and-push-docker-images)
- [Deploy to Kubernetes](#deploy-to-kubernetes)
  - [Database Configuration](#database-configuration)
  - [Service Deployment](#service-deployment)
  - [New Services Overview](#new-services-overview)
- [Access Your Application](#access-your-application)
- [Updating Services After Code Changes](#updating-services-after-code-changes)
  - [Method 1: Quick Update with Image Tag](#method-1-quick-update-with-image-tag)
  - [Method 2: Update with Manifest Change](#method-2-update-with-manifest-change)
  - [Method 3: Force Pull Latest Image](#method-3-force-pull-latest-image)
  - [Frontend-Specific Updates](#frontend-specific-updates)
  - [Best Practices for Updates](#best-practices-for-updates)
  - [Complete Update Example](#complete-update-example)
- [Troubleshooting](#troubleshooting)

---

## Architecture Overview

### Services Deployed (12 Total)

**Core Microservices (6):**
- **catalog-service** (Go) - Book catalog management on port 8081
- **user-service** (Go) - User authentication & profiles on port 8082
- **cart-service** (Go) - Shopping cart functionality on port 8083
- **order-service** (Go) - Order processing on port 8084
- **admin-service** (Go) - Analytics & admin dashboard on port 8090
- **inventory-service** (Python/FastAPI) - Stock tracking & reservations on port 8086

**ML-Powered Services (2):**
- **review-service** (Python/FastAPI) - Book reviews with sentiment analysis on port 8088
- **recommendation-service** (Python/FastAPI) - Personalized recommendations on port 8089

**Infrastructure Services (4):**
- **api-gateway** (Go) - Routes all requests on port 8080
- **frontend** (React/Vite) - Customer-facing UI on port 80
- **PostgreSQL** - Shared database server on port 5432
- **Redis** - Caching layer on port 6379

### Database Architecture

**PostgreSQL Databases (6):**
- `catalog_db` - Product catalog data
- `userdb` - User accounts and authentication
- `orderdb` - Order history and transactions
- `inventory_db` - Stock levels and reservations
- `reviews_db` - Book reviews and ratings
- `bookstore` - Admin analytics and aggregated data

**Service Users:**
- `userservice` - Access to `userdb`
- `orderservice` + `orderuser` - Access to `orderdb`
- Service users automatically created via init scripts

---

## Prerequisites

### System Requirements
- Windows 10/11 with PowerShell
- Docker Desktop installed and running
- Git installed
- Internet connection

### AWS Academy Account
- Active AWS Academy Learner Lab session
- Access to AWS Console
- Ability to create EKS clusters and ECR repositories

---

## AWS Academy Setup

### 1. Start Your AWS Academy Lab
```powershell
# Open AWS Academy Learner Lab
# Click "Start Lab" button
# Wait for the AWS indicator to turn green
```

### 2. Get AWS Credentials
In AWS Academy, click "AWS Details" → "Show" under AWS CLI credentials.

Copy the credentials and paste into PowerShell:

```powershell
# Configure AWS CLI with temporary credentials
aws configure

# When prompted, enter:
# AWS Access Key ID: [from AWS Academy]
# AWS Secret Access Key: [from AWS Academy]
# Default region name: us-east-1
# Default output format: json

# Set session token (REQUIRED for AWS Academy)
$env:AWS_SESSION_TOKEN = "YOUR_SESSION_TOKEN_HERE"
```

### 3. Verify AWS Access
```powershell
# Test AWS connectivity
aws sts get-caller-identity

# Should show your AWS Academy account details
```

### 4. Get Your VPC and Subnet Information
```powershell
# Get VPC ID
aws ec2 describe-vpcs --query "Vpcs[0].VpcId" --output text

# Get Subnet IDs
aws ec2 describe-subnets --query "Subnets[?MapPublicIpOnLaunch==\`true\`].[SubnetId,AvailabilityZone]" --output table
```

Save these values - you'll need them later:
- **VPC ID**: `vpc-0718312e7f365e4b8` (your value will differ)
- **Subnet IDs**: 3 public subnets in different availability zones

---

## Install Required Tools

### 1. Install eksctl
```powershell
# Create eksctl directory
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\eksctl"

# Download eksctl
$eksctlUrl = "https://github.com/eksctl-io/eksctl/releases/latest/download/eksctl_Windows_amd64.zip"
Invoke-WebRequest -Uri $eksctlUrl -OutFile "$env:USERPROFILE\eksctl\eksctl.zip"

# Extract
Expand-Archive -Path "$env:USERPROFILE\eksctl\eksctl.zip" -DestinationPath "$env:USERPROFILE\eksctl" -Force

# Add to PATH for current session
$env:PATH += ";$env:USERPROFILE\eksctl"

# Verify installation
eksctl version
```

### 2. Install kubectl (if not already installed)
```powershell
# Download kubectl
$kubectlUrl = "https://dl.k8s.io/release/v1.28.0/bin/windows/amd64/kubectl.exe"
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\kubectl"
Invoke-WebRequest -Uri $kubectlUrl -OutFile "$env:USERPROFILE\kubectl\kubectl.exe"

# Add to PATH
$env:PATH += ";$env:USERPROFILE\kubectl"

# Verify
kubectl version --client
```

---

## Create EKS Cluster

### 1. Create Cluster Configuration File

Create `cluster-config.yaml`:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: my-bookstore
  region: us-east-1
  version: "1.28"

# Use existing VPC (AWS Academy restriction)
vpc:
  id: "vpc-0718312e7f365e4b8"  # Replace with YOUR VPC ID
  subnets:
    public:
      us-east-1a: { id: subnet-0b0a497aca8761715 }  # Replace with YOUR subnet IDs
      us-east-1b: { id: subnet-08d684911389ce0b4 }
      us-east-1d: { id: subnet-0229448bb3e7df8fd }

# Use existing IAM role (AWS Academy provides LabRole)
iam:
  withOIDC: false
  serviceRoleARN: arn:aws:iam::905418472239:role/LabRole  # Replace with YOUR account ID

# Don't create node group in cluster creation (will be created separately)
nodeGroups: []

# Install required addons
addons:
  - name: vpc-cni
  - name: kube-proxy
  - name: coredns
```

**Important**: Update these values:
- `vpc.id`: Your VPC ID from step 4 above
- `vpc.subnets.public`: Your 3 public subnet IDs
- `iam.serviceRoleARN`: Replace `905418472239` with your AWS account ID

### 2. Create EKS Cluster
```powershell
# Create the cluster (takes 10-15 minutes)
eksctl create cluster -f cluster-config.yaml

# Wait for completion - you'll see:
# "EKS cluster 'my-bookstore' in 'us-east-1' region is ready"
```

### 3. Install EBS CSI Driver Addon
```powershell
# Required for persistent storage (PostgreSQL)
aws eks create-addon --cluster-name my-bookstore --addon-name aws-ebs-csi-driver --region us-east-1

# Wait 2-3 minutes for installation
aws eks describe-addon --cluster-name my-bookstore --addon-name aws-ebs-csi-driver --region us-east-1 --query "addon.status"
# Wait until status shows "ACTIVE"
```

### 4. Create Node Group
```powershell
# Create worker nodes using AWS CLI (console creation may have issues)
aws eks create-nodegroup `
  --cluster-name my-bookstore `
  --nodegroup-name my-bookstore-workers `
  --node-role arn:aws:iam::905418472239:role/LabRole `
  --subnets subnet-0b0a497aca8761715 subnet-08d684911389ce0b4 subnet-0229448bb3e7df8fd `
  --instance-types t3.small `
  --scaling-config minSize=2,maxSize=3,desiredSize=2 `
  --disk-size 20 `
  --region us-east-1

# Wait 5-10 minutes for node provisioning
```

**Important**: Replace these values:
- `--node-role`: Use `arn:aws:iam::YOUR_ACCOUNT_ID:role/LabRole`
- `--subnets`: Use your 3 public subnet IDs (space-separated)

### 5. Configure kubectl
```powershell
# Update kubeconfig to access the cluster
aws eks update-kubeconfig --name my-bookstore --region us-east-1

# Verify connection
kubectl get nodes

# Should show 2 nodes in Ready state:
# NAME                            STATUS   ROLES    AGE
# ip-172-31-1-130.ec2.internal    Ready    <none>   2m
# ip-172-31-31-249.ec2.internal   Ready    <none>   2m
```

---

## Build and Push Docker Images

### 1. Create ECR Repositories
```powershell
# Create repositories for all services
aws ecr create-repository --repository-name catalog-service --region us-east-1
aws ecr create-repository --repository-name user-service --region us-east-1
aws ecr create-repository --repository-name cart-service --region us-east-1
aws ecr create-repository --repository-name order-service --region us-east-1
aws ecr create-repository --repository-name api-gateway --region us-east-1
aws ecr create-repository --repository-name customer-frontend --region us-east-1

# New services
aws ecr create-repository --repository-name inventory-service --region us-east-1
aws ecr create-repository --repository-name review-service --region us-east-1
aws ecr create-repository --repository-name recommendation-service --region us-east-1
aws ecr create-repository --repository-name admin-service --region us-east-1
```

### 2. Login to ECR
```powershell
# Get your AWS Account ID
$ACCOUNT_ID = (aws sts get-caller-identity --query Account --output text)

# Login to ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin "$ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com"

# Should see: "Login Succeeded"
```

### 3. Build and Push Images
```powershell
# Navigate to your project directory
cd C:\Users\YOUR_USERNAME\path\to\My-Distributed-Bookstore

# Get your account ID
$ACCOUNT_ID = (aws sts get-caller-identity --query Account --output text)
$ECR_REGISTRY = "$ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com"

# Build and push catalog-service
cd services/catalog-service
docker build -t "$ECR_REGISTRY/catalog-service:latest" .
docker push "$ECR_REGISTRY/catalog-service:latest"
cd ../..

# Build and push user-service
cd services/user-service
docker build -t "$ECR_REGISTRY/user-service:latest" .
docker push "$ECR_REGISTRY/user-service:latest"
cd ../..

# Build and push cart-service
cd services/cart-service
docker build -t "$ECR_REGISTRY/cart-service:latest" .
docker push "$ECR_REGISTRY/cart-service:latest"
cd ../..

# Build and push order-service
cd services/order-service
docker build -t "$ECR_REGISTRY/order-service:latest" .
docker push "$ECR_REGISTRY/order-service:latest"
cd ../..

# Build and push api-gateway
cd services/api-gateway
docker build -t "$ECR_REGISTRY/api-gateway:latest" .
docker push "$ECR_REGISTRY/api-gateway:latest"
cd ../..

# Build and push frontend
cd frontend/customer-app
docker build -t "$ECR_REGISTRY/customer-frontend:latest" .
docker push "$ECR_REGISTRY/customer-frontend:latest"
cd ../..

# Build and push inventory-service (Python/FastAPI)
docker build -t "$ECR_REGISTRY/inventory-service:latest" -f services/inventory-service/Dockerfile services/inventory-service/
docker push "$ECR_REGISTRY/inventory-service:latest"

# Build and push review-service (Python/FastAPI)
docker build -t "$ECR_REGISTRY/review-service:latest" -f services/review-service/Dockerfile services/review-service/
docker push "$ECR_REGISTRY/review-service:latest"

# Build and push recommendation-service (Python/FastAPI)
docker build -t "$ECR_REGISTRY/recommendation-service:latest" -f services/recommendation-service/Dockerfile services/recommendation-service/
docker push "$ECR_REGISTRY/recommendation-service:latest"

# Build and push admin-service (Go)
docker build -t "$ECR_REGISTRY/admin-service:latest" -f services/admin-service/Dockerfile services/admin-service/
docker push "$ECR_REGISTRY/admin-service:latest"
```

**Note:** The new services (inventory, review, recommendation, admin) use updated Dockerfiles with `go mod tidy && go mod download` for proper dependency management.

### 4. Update Kubernetes Manifests with ECR URLs
```powershell
# Update catalog-service
$manifestPath = ".\infrastructure\k8s\services\catalog-service\deployment.yaml"
(Get-Content $manifestPath) -replace 'image: catalog-service:latest', "image: $ECR_REGISTRY/catalog-service:latest" -replace 'imagePullPolicy: Never', 'imagePullPolicy: Always' | Set-Content $manifestPath

# Update user-service
$manifestPath = ".\infrastructure\k8s\services\user-service\deployment.yaml"
(Get-Content $manifestPath) -replace 'image: user-service:latest', "image: $ECR_REGISTRY/user-service:latest" -replace 'imagePullPolicy: Never', 'imagePullPolicy: Always' | Set-Content $manifestPath

# Update cart-service
$manifestPath = ".\infrastructure\k8s\services\cart-service\deployment.yaml"
(Get-Content $manifestPath) -replace 'image: cart-service:latest', "image: $ECR_REGISTRY/cart-service:latest" -replace 'imagePullPolicy: Never', 'imagePullPolicy: Always' | Set-Content $manifestPath

# Update order-service
$manifestPath = ".\infrastructure\k8s\services\order-service\deployment.yaml"
(Get-Content $manifestPath) -replace 'image: order-service:latest', "image: $ECR_REGISTRY/order-service:latest" -replace 'imagePullPolicy: Never', 'imagePullPolicy: Always' | Set-Content $manifestPath

# Update api-gateway
$manifestPath = ".\infrastructure\k8s\services\api-gateway\deployment.yaml"
(Get-Content $manifestPath) -replace 'image: api-gateway:latest', "image: $ECR_REGISTRY/api-gateway:latest" -replace 'imagePullPolicy: Never', 'imagePullPolicy: Always' | Set-Content $manifestPath

# Update frontend
$manifestPath = ".\infrastructure\k8s\frontend\deployment.yaml"
(Get-Content $manifestPath) -replace 'image: customer-frontend:latest', "image: $ECR_REGISTRY/customer-frontend:latest" -replace 'imagePullPolicy: Never', 'imagePullPolicy: Always' | Set-Content $manifestPath

# Update inventory-service
$manifestPath = ".\infrastructure\k8s\services\inventory-service\deployment.yaml"
(Get-Content $manifestPath) -replace 'image: inventory-service:latest', "image: $ECR_REGISTRY/inventory-service:latest" -replace 'imagePullPolicy: Never', 'imagePullPolicy: Always' | Set-Content $manifestPath

# Update review-service
$manifestPath = ".\infrastructure\k8s\services\review-service\deployment.yaml"
(Get-Content $manifestPath) -replace 'image: review-service:latest', "image: $ECR_REGISTRY/review-service:latest" -replace 'imagePullPolicy: Never', 'imagePullPolicy: Always' | Set-Content $manifestPath

# Update recommendation-service
$manifestPath = ".\infrastructure\k8s\services\recommendation-service\deployment.yaml"
(Get-Content $manifestPath) -replace 'image: recommendation-service:latest', "image: $ECR_REGISTRY/recommendation-service:latest" -replace 'imagePullPolicy: Never', 'imagePullPolicy: Always' | Set-Content $manifestPath

# Update admin-service
$manifestPath = ".\infrastructure\k8s\services\admin-service\deployment.yaml"
(Get-Content $manifestPath) -replace 'image: admin-service:latest', "image: $ECR_REGISTRY/admin-service:latest" -replace 'imagePullPolicy: Never', 'imagePullPolicy: Always' | Set-Content $manifestPath
```

---

## Deploy to Kubernetes

### 1. Create Namespace
```powershell
kubectl apply -f .\infrastructure\k8s\namespaces\development.yaml
```

### 2. Create Secrets
```powershell
kubectl apply -f .\infrastructure\k8s\secrets\
```

### 3. Deploy Databases
```powershell
kubectl apply -f .\infrastructure\k8s\databases\
```

**Important:** Wait for PostgreSQL to be fully running before proceeding:
```powershell
# Wait for postgres to be ready
kubectl wait --for=condition=ready pod/postgres-0 -n bookstore-dev --timeout=300s
```

### 3.1. Initialize PostgreSQL Databases

The PostgreSQL init scripts will automatically create the required databases and users. However, if you need to manually verify or create additional databases:

```powershell
# List all databases
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d postgres -c "\l"

# Expected databases:
# - catalog_db
# - inventory_db
# - reviews_db
# - userdb
# - orderdb
# - bookstore (for admin-service)
```

If any database is missing, create it manually:
```powershell
# Create missing database (example)
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d postgres -c "CREATE DATABASE database_name;"

# Grant permissions to service user
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d postgres -c "GRANT ALL PRIVILEGES ON DATABASE database_name TO service_user;"
```

### 4. Deploy ConfigMaps
```powershell
kubectl apply -f .\infrastructure\k8s\configmaps\
```

### 5. Deploy Services

Deploy all services in the following order:

```powershell
# Core services (required by others)
kubectl apply -f .\infrastructure\k8s\services\catalog-service\
kubectl apply -f .\infrastructure\k8s\services\user-service\
kubectl apply -f .\infrastructure\k8s\services\cart-service\
kubectl apply -f .\infrastructure\k8s\services\order-service\

# New services (deployed after core services)
kubectl apply -f .\infrastructure\k8s\services\inventory-service\
kubectl apply -f .\infrastructure\k8s\services\review-service\
kubectl apply -f .\infrastructure\k8s\services\recommendation-service\
kubectl apply -f .\infrastructure\k8s\services\admin-service\

# API Gateway (depends on all services)
kubectl apply -f .\infrastructure\k8s\services\api-gateway\
```

**Service Overview:**
- **catalog-service**: Book catalog management
- **user-service**: User authentication & profiles
- **cart-service**: Shopping cart functionality
- **order-service**: Order processing
- **inventory-service**: ✨ Stock tracking & reservations
- **review-service**: ✨ Book reviews with ML sentiment analysis
- **recommendation-service**: ✨ Personalized book recommendations
- **admin-service**: ✨ Analytics & admin dashboard
- **api-gateway**: Routes requests to all services

✨ = New services added in this deployment

### 6. Deploy Frontend
```powershell
kubectl apply -f .\infrastructure\k8s\frontend\
```

### 7. Scale Down for Resource Management (if needed)
```powershell
# If you have limited resources on t3.small nodes
kubectl scale deployment api-gateway --replicas=1 -n bookstore-dev
kubectl scale deployment frontend --replicas=1 -n bookstore-dev
kubectl scale deployment cart-service --replicas=1 -n bookstore-dev
kubectl scale deployment catalog-service --replicas=1 -n bookstore-dev
kubectl scale deployment order-service --replicas=1 -n bookstore-dev
kubectl scale deployment user-service --replicas=1 -n bookstore-dev
```

### 8. Check Deployment Status
```powershell
# Watch pods until all are running
kubectl get pods -n bookstore-dev -w

# Press Ctrl+C to stop watching

# Check final status
kubectl get pods -n bookstore-dev

# Expected output (all pods should be Running):
# NAME                                      READY   STATUS    RESTARTS   AGE
# admin-service-xxx                         1/1     Running   0          5m
# api-gateway-xxx                           2/2     Running   0          5m
# cart-service-xxx                          2/2     Running   0          5m
# catalog-service-xxx                       2/2     Running   0          5m
# frontend-xxx                              2/2     Running   0          5m
# inventory-service-xxx                     1/1     Running   0          5m
# order-service-xxx                         2/2     Running   0          5m
# postgres-0                                1/1     Running   0          5m
# recommendation-service-xxx                1/1     Running   0          5m
# redis-xxx                                 1/1     Running   0          5m
# review-service-xxx                        1/1     Running   0          5m
# user-service-xxx                          2/2     Running   0          5m
```

**Total Pods:** ~18 pods across 12 services + databases

### 9. Verify Service Health

Check individual service logs if any pod is not Running:

```powershell
# Check specific service logs
kubectl logs -n bookstore-dev -l app=service-name --tail=50

# Examples:
kubectl logs -n bookstore-dev -l app=admin-service --tail=50
kubectl logs -n bookstore-dev -l app=inventory-service --tail=50
kubectl logs -n bookstore-dev -l app=review-service --tail=50
kubectl logs -n bookstore-dev -l app=recommendation-service --tail=50
```

Common issues and fixes:

**admin-service CrashLoopBackOff:**
```powershell
# Check if 'bookstore' database exists
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d postgres -c "\l" | grep bookstore

# If missing, create it:
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d postgres -c "CREATE DATABASE bookstore;"

# Restart admin-service
kubectl delete pod -n bookstore-dev -l app=admin-service
```

**user-service or order-service errors:**
```powershell
# Check if databases and users exist
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d postgres -c "\l" | grep -E "userdb|orderdb"
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d postgres -c "\du" | grep -E "userservice|orderservice|orderuser"

# If missing, the postgres init scripts should have created them
# Check init scripts: infrastructure/k8s/databases/postgres-initdb-configmap.yaml
```

**ImagePullBackOff errors:**
```powershell
# Verify ECR images exist
aws ecr describe-images --repository-name service-name --region us-east-1

# If image missing, rebuild and push:
docker build -t service-name:latest -f services/service-name/Dockerfile services/service-name/
docker tag service-name:latest $ECR_REGISTRY/service-name:latest
docker push $ECR_REGISTRY/service-name:latest
```

---

## New Services Overview

This deployment includes 4 new advanced services that extend the bookstore functionality:

### 1. Admin Service (Go)
**Port:** 8090  
**Purpose:** Business analytics and admin dashboard

**Features:**
- Real-time sales analytics
- Inventory reports and low-stock alerts
- User growth tracking
- Top-selling books analysis
- Aggregated dashboard data

**Database:** `bookstore` (aggregated analytics data)

**Endpoints:**
- `GET /health` - Health check
- `GET /ready` - Readiness probe
- `GET /api/v1/admin/dashboard` - Main dashboard
- `GET /api/v1/admin/analytics/sales` - Sales analytics
- `GET /api/v1/admin/analytics/inventory` - Inventory report
- `GET /api/v1/admin/analytics/users` - User growth
- `GET /api/v1/admin/top-books` - Top-selling books

**Authentication:** Requires JWT token with admin role

---

### 2. Inventory Service (Python/FastAPI)
**Port:** 8086  
**Purpose:** Real-time stock tracking and reservation management

**Features:**
- Stock level tracking per book
- Temporary stock reservations for cart items
- Automatic reservation expiry (background task)
- Stock movement history
- Low stock alerts
- Bulk stock operations

**Database:** `inventory_db`

**Endpoints:**
- `GET /health` - Health check
- `GET /api/v1/inventory` - List all inventory
- `GET /api/v1/inventory/{book_id}` - Get stock level
- `POST /api/v1/inventory/{book_id}/reserve` - Reserve stock
- `POST /api/v1/inventory/{book_id}/release` - Release reservation
- `PUT /api/v1/inventory/{book_id}` - Update stock
- `GET /api/v1/inventory/low-stock` - Get low-stock items

**Background Tasks:**
- Reservation expiry checker (runs every 5 minutes)

---

### 3. Review Service (Python/FastAPI)
**Port:** 8088  
**Purpose:** Book reviews with ML-powered sentiment analysis

**Features:**
- Create, read, update, delete reviews
- Star ratings (1-5)
- Automatic sentiment analysis (NLTK + TextBlob)
- Sentiment classification: Positive/Neutral/Negative
- Review voting (helpful/not helpful)
- Review statistics and aggregations
- Filter by rating, sentiment, date

**Database:** `reviews_db`

**Machine Learning:**
- **NLTK** - Natural Language Toolkit for text processing
- **TextBlob** - Sentiment polarity analysis
- Automatic classification based on text content

**Endpoints:**
- `GET /health` - Health check
- `GET /api/v1/reviews` - List reviews (with filters)
- `GET /api/v1/reviews/{id}` - Get specific review
- `POST /api/v1/reviews` - Create review (with auto-sentiment)
- `PUT /api/v1/reviews/{id}` - Update review
- `DELETE /api/v1/reviews/{id}` - Delete review
- `GET /api/v1/reviews/book/{book_id}` - Reviews for a book
- `GET /api/v1/reviews/book/{book_id}/stats` - Review statistics
- `POST /api/v1/reviews/{id}/vote` - Vote on review

**Sentiment Analysis:**
```python
# Automatic analysis on review creation
polarity = TextBlob(review_text).sentiment.polarity
if polarity > 0.3:
    sentiment = "positive"
elif polarity < -0.3:
    sentiment = "negative"
else:
    sentiment = "neutral"
```

---

### 4. Recommendation Service (Python/FastAPI)
**Port:** 8089  
**Purpose:** Personalized book recommendations

**Features:**
- Tag-based recommendations (similar books)
- Collaborative filtering (users who bought X also bought Y)
- Popularity-based recommendations
- User preference tracking
- Multiple recommendation strategies

**Database:** SQLite (embedded, no postgres needed)

**Endpoints:**
- `GET /api/v1/health` - Health check
- `GET /api/v1/ready` - Readiness check
- `GET /api/v1/recommendations/{user_id}` - Personalized recommendations
- `GET /api/v1/recommendations/book/{book_id}/similar` - Similar books
- `GET /api/v1/recommendations/popular` - Popular books
- `POST /api/v1/recommendations/track` - Track user interaction

**Recommendation Strategies:**
1. **Content-Based:** Recommends books with similar tags/categories
2. **Collaborative Filtering:** Based on user purchase patterns
3. **Popularity:** Trending and best-selling books

---

### Service Dependencies

```
┌─────────────┐
│  Frontend   │
└──────┬──────┘
       │
┌──────▼──────────┐
│  API Gateway    │
└────┬─────┬──────┘
     │     │
     │     ├─────► catalog-service
     │     ├─────► user-service
     │     ├─────► cart-service ──────► inventory-service (check stock)
     │     ├─────► order-service ─────► inventory-service (reserve stock)
     │     ├─────► admin-service ──────► All services (analytics)
     │     ├─────► review-service
     │     └─────► recommendation-service
     │
┌────▼────────┐
│ PostgreSQL  │  ← 6 databases
└─────────────┘
```

---

## Database Configuration

### PostgreSQL Init Scripts

The `postgres-initdb-configmap.yaml` contains SQL scripts that run on first PostgreSQL startup:

**File:** `infrastructure/k8s/databases/postgres-initdb-configmap.yaml`

```yaml
data:
  01-create-databases.sql: |
    CREATE DATABASE reviews_db;
    CREATE DATABASE inventory_db;
    CREATE DATABASE userdb;
    CREATE DATABASE orderdb;
    
  02-create-users.sql: |
    CREATE USER userservice WITH PASSWORD 'dev_password';
    CREATE USER orderservice WITH PASSWORD 'dev_password';
    CREATE USER orderuser WITH PASSWORD 'dev_password';
    
  03-grant-permissions.sql: |
    -- Database-level permissions
    GRANT ALL PRIVILEGES ON DATABASE userdb TO userservice;
    GRANT ALL PRIVILEGES ON DATABASE orderdb TO orderservice, orderuser;
    
    -- Schema-level permissions (userdb)
    \c userdb
    GRANT ALL ON SCHEMA public TO userservice;
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO userservice;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO userservice;
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO userservice;
    
    -- Schema-level permissions (orderdb)
    \c orderdb
    GRANT ALL ON SCHEMA public TO orderservice, orderuser;
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO orderservice, orderuser;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO orderservice, orderuser;
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO orderservice, orderuser;
```

### Manual Database Creation (if needed)

If you add a service that needs a new database:

```powershell
# 1. Create the database
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d postgres -c "CREATE DATABASE new_database;"

# 2. Create a service user
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d postgres -c "CREATE USER new_user WITH PASSWORD 'dev_password';"

# 3. Grant database permissions
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d postgres -c "GRANT ALL PRIVILEGES ON DATABASE new_database TO new_user;"

# 4. Grant schema permissions
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d new_database -c "GRANT ALL ON SCHEMA public TO new_user;"
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d new_database -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO new_user;"
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d new_database -c "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO new_user;"
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d new_database -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO new_user;"
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d new_database -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO new_user;"

# 5. Update the init-db-configmap.yaml to make it permanent
# Add the above SQL commands to the ConfigMap and reapply
```

### Verify Database Setup

```powershell
# List all databases
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d postgres -c "\l"

# List all users
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d postgres -c "\du"

# Check database permissions
kubectl exec -n bookstore-dev postgres-0 -- psql -U bookstore -d postgres -c "\l" | grep -E "userdb|orderdb|inventory_db|reviews_db|bookstore"

# Expected output should show all 6 databases
```

---

## Access Your Application

### 1. Expose Frontend via LoadBalancer
```powershell
# Change frontend service to LoadBalancer
kubectl patch svc frontend -n bookstore-dev -p '{\"spec\":{\"type\":\"LoadBalancer\"}}'

# Wait 2-3 minutes for AWS to provision the load balancer
```

### 2. Get Application URL
```powershell
# Get the LoadBalancer URL
kubectl get svc frontend -n bookstore-dev -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'

# Example output:
# aeb72405221874e0d9c3c643d14d6cf3-1452808211.us-east-1.elb.amazonaws.com
```

### 3. Access Your Application
```powershell
# The URL will be:
# http://[LOADBALANCER-HOSTNAME]

# Open in your browser:
$FRONTEND_URL = (kubectl get svc frontend -n bookstore-dev -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
Start-Process "http://$FRONTEND_URL"
```

### 4. Verify All Services
```powershell
# Check all services
kubectl get svc -n bookstore-dev

# Expected services:
# NAME              TYPE           EXTERNAL-IP
# api-gateway       NodePort       <none>
# cart-service      ClusterIP      <none>
# catalog-service   ClusterIP      <none>
# frontend          LoadBalancer   [AWS-ELB-URL]
# order-service     ClusterIP      <none>
# postgres          ClusterIP      None
# redis             ClusterIP      <none>
# user-service      ClusterIP      <none>
```

---

## Updating Services After Code Changes

This section explains how to update your deployed services when you make code changes to your application.

### Overview of Update Process

When you change code in a service, you need to:
1. **Rebuild** the Docker image with your changes
2. **Push** the new image to ECR
3. **Restart** the Kubernetes pods to use the new image
4. **(Optional)** Clear browser cache if it's the frontend

### Method 1: Quick Update with Image Tag

#### Step 1: Build and Tag New Image
```powershell
# Navigate to the service directory
cd services/[SERVICE-NAME]  # e.g., catalog-service, user-service

# Build with a version tag (recommended)
$VERSION = "v1.1"  # Increment this each time
$AWS_ACCOUNT = "905418472239"  # Your AWS account ID
$SERVICE = "catalog-service"   # Change to your service name

docker build -t ${AWS_ACCOUNT}.dkr.ecr.us-east-1.amazonaws.com/${SERVICE}:${VERSION} .
docker build -t ${AWS_ACCOUNT}.dkr.ecr.us-east-1.amazonaws.com/${SERVICE}:latest .
```

#### Step 2: Push to ECR
```powershell
# Login to ECR (if needed)
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin ${AWS_ACCOUNT}.dkr.ecr.us-east-1.amazonaws.com

# Push both tags
docker push ${AWS_ACCOUNT}.dkr.ecr.us-east-1.amazonaws.com/${SERVICE}:${VERSION}
docker push ${AWS_ACCOUNT}.dkr.ecr.us-east-1.amazonaws.com/${SERVICE}:latest
```

#### Step 3: Restart Pods
```powershell
# Option A: Delete pods (they'll be recreated automatically)
kubectl delete pods -l app=${SERVICE} -n bookstore-dev

# Option B: Rollout restart (recommended)
kubectl rollout restart deployment/${SERVICE} -n bookstore-dev

# Wait for rollout to complete
kubectl rollout status deployment/${SERVICE} -n bookstore-dev
```

#### Step 4: Verify Update
```powershell
# Check pods are running with new image
kubectl get pods -l app=${SERVICE} -n bookstore-dev

# Check image version
kubectl describe pod -l app=${SERVICE} -n bookstore-dev | Select-String "Image:"

# Check logs for startup
kubectl logs -l app=${SERVICE} -n bookstore-dev --tail=50
```

### Method 2: Update with Manifest Change

If you need to change configuration (environment variables, resources, etc.):

#### Step 1: Update Deployment Manifest
```powershell
# Edit the deployment file
code infrastructure/k8s/services/${SERVICE}/deployment.yaml

# Update the image tag:
# OLD: image: 905418472239.dkr.ecr.us-east-1.amazonaws.com/catalog-service:latest
# NEW: image: 905418472239.dkr.ecr.us-east-1.amazonaws.com/catalog-service:v1.1
```

#### Step 2: Apply Updated Manifest
```powershell
# Apply the changes
kubectl apply -f infrastructure/k8s/services/${SERVICE}/

# Kubernetes will perform a rolling update automatically
```

#### Step 3: Monitor Rollout
```powershell
# Watch the rollout
kubectl rollout status deployment/${SERVICE} -n bookstore-dev

# See rollout history
kubectl rollout history deployment/${SERVICE} -n bookstore-dev
```

### Method 3: Force Pull Latest Image

If using `:latest` tag and want to force pull new version:

```powershell
# Set imagePullPolicy to Always (should already be set)
kubectl patch deployment ${SERVICE} -n bookstore-dev -p '{"spec":{"template":{"spec":{"containers":[{"name":"'${SERVICE}'","imagePullPolicy":"Always"}]}}}}'

# Delete pods to force recreation
kubectl delete pods -l app=${SERVICE} -n bookstore-dev
```

### Frontend-Specific Updates

Frontend requires special handling due to browser caching:

#### Step 1: Update Environment Variables (if needed)
```powershell
# Edit frontend config
code infrastructure/k8s/configmaps/frontend-config.yaml

# Apply changes
kubectl apply -f infrastructure/k8s/configmaps/frontend-config.yaml
```

#### Step 2: Rebuild Frontend Image
```powershell
cd frontend/customer-app

# Build with empty API URL (uses nginx proxy)
docker build --no-cache --build-arg VITE_API_URL="" -t ${AWS_ACCOUNT}.dkr.ecr.us-east-1.amazonaws.com/customer-frontend:latest .

# Push to ECR
docker push ${AWS_ACCOUNT}.dkr.ecr.us-east-1.amazonaws.com/customer-frontend:latest
```

#### Step 3: Restart Frontend Pods
```powershell
# Restart deployment
kubectl rollout restart deployment/frontend -n bookstore-dev

# Wait for rollout
kubectl rollout status deployment/frontend -n bookstore-dev
```

#### Step 4: Clear Browser Cache
**IMPORTANT:** Users must clear their browser cache to see changes!

**Windows/Linux:**
- Chrome/Edge: `Ctrl + Shift + R` or `Ctrl + F5`
- Firefox: `Ctrl + Shift + R`

**macOS:**
- Chrome/Edge/Firefox: `Cmd + Shift + R`

**Alternative:** Open in Incognito/Private window

### Update All Services Script

Create a PowerShell script to update all services at once:

```powershell
# update-all-services.ps1
$AWS_ACCOUNT = "905418472239"
$SERVICES = @("catalog-service", "user-service", "cart-service", "order-service", "api-gateway")

foreach ($SERVICE in $SERVICES) {
    Write-Host "Updating $SERVICE..." -ForegroundColor Cyan
    
    # Rollout restart
    kubectl rollout restart deployment/$SERVICE -n bookstore-dev
    
    # Wait for completion
    kubectl rollout status deployment/$SERVICE -n bookstore-dev
    
    Write-Host "✓ $SERVICE updated" -ForegroundColor Green
}

Write-Host "`nAll services updated!" -ForegroundColor Green
```

### Rollback a Failed Update

If an update causes issues:

```powershell
# View rollout history
kubectl rollout history deployment/${SERVICE} -n bookstore-dev

# Rollback to previous version
kubectl rollout undo deployment/${SERVICE} -n bookstore-dev

# Rollback to specific revision
kubectl rollout undo deployment/${SERVICE} -n bookstore-dev --to-revision=2

# Check rollback status
kubectl rollout status deployment/${SERVICE} -n bookstore-dev
```

### Best Practices for Updates

#### 1. Use Version Tags
```powershell
# ✅ GOOD: Use semantic versioning
docker build -t myservice:v1.2.3 .
docker build -t myservice:v1.2 .
docker build -t myservice:latest .

# ❌ BAD: Only use :latest
docker build -t myservice:latest .
```

#### 2. Test Before Deploying
```powershell
# Run image locally first
docker run -p 8080:8080 ${AWS_ACCOUNT}.dkr.ecr.us-east-1.amazonaws.com/${SERVICE}:v1.1

# Test the service
curl http://localhost:8080/health
```

#### 3. Use Rolling Updates
- Kubernetes performs rolling updates by default
- Ensures zero downtime
- Old pods stay running until new ones are healthy

#### 4. Monitor During Updates
```powershell
# Watch pods in real-time
kubectl get pods -n bookstore-dev -w

# Check logs of new pods
kubectl logs -f deployment/${SERVICE} -n bookstore-dev

# Check events for issues
kubectl get events -n bookstore-dev --sort-by='.lastTimestamp'
```

#### 5. Update ConfigMaps and Secrets First
```powershell
# Update config/secrets before deployment
kubectl apply -f infrastructure/k8s/configmaps/${SERVICE}-config.yaml
kubectl apply -f infrastructure/k8s/secrets/

# Then update deployment (pods will restart automatically)
kubectl apply -f infrastructure/k8s/services/${SERVICE}/
```

### Complete Update Example

Here's a complete example updating the catalog-service:

```powershell
# 1. Navigate to service
cd services/catalog-service

# 2. Make your code changes
# ... edit files ...

# 3. Build new image
$VERSION = "v1.1.0"
docker build -t 905418472239.dkr.ecr.us-east-1.amazonaws.com/catalog-service:${VERSION} .
docker build -t 905418472239.dkr.ecr.us-east-1.amazonaws.com/catalog-service:latest .

# 4. Test locally (optional)
docker run --rm -p 8080:8080 905418472239.dkr.ecr.us-east-1.amazonaws.com/catalog-service:latest
# Test: curl http://localhost:8080/health
# Stop: Ctrl+C

# 5. Login to ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 905418472239.dkr.ecr.us-east-1.amazonaws.com

# 6. Push images
docker push 905418472239.dkr.ecr.us-east-1.amazonaws.com/catalog-service:${VERSION}
docker push 905418472239.dkr.ecr.us-east-1.amazonaws.com/catalog-service:latest

# 7. Update deployment (optional - if you want to pin to version)
cd ../../infrastructure/k8s/services/catalog-service
# Edit deployment.yaml to use :v1.1.0 instead of :latest

# 8. Apply changes
kubectl apply -f .

# Or simply restart if using :latest tag
kubectl rollout restart deployment/catalog-service -n bookstore-dev

# 9. Monitor rollout
kubectl rollout status deployment/catalog-service -n bookstore-dev

# 10. Verify
kubectl get pods -l app=catalog-service -n bookstore-dev
kubectl logs -l app=catalog-service -n bookstore-dev --tail=20

# 11. Test the service
$FRONTEND_URL = (kubectl get svc frontend -n bookstore-dev -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
Start-Process "http://$FRONTEND_URL"
```

### Troubleshooting Updates

#### Image Not Updating
```powershell
# Check current image
kubectl describe pod -l app=${SERVICE} -n bookstore-dev | Select-String "Image:"

# Force pull policy
kubectl patch deployment ${SERVICE} -n bookstore-dev -p '{"spec":{"template":{"spec":{"containers":[{"name":"'${SERVICE}'","imagePullPolicy":"Always"}]}}}}'

# Delete and recreate pods
kubectl delete pods -l app=${SERVICE} -n bookstore-dev
```

#### Pods Not Starting After Update
```powershell
# Check pod status
kubectl get pods -l app=${SERVICE} -n bookstore-dev

# Describe failing pod
kubectl describe pod POD_NAME -n bookstore-dev

# Check logs
kubectl logs POD_NAME -n bookstore-dev

# Rollback if needed
kubectl rollout undo deployment/${SERVICE} -n bookstore-dev
```

#### Frontend Changes Not Visible
1. Verify new image was built and pushed
2. Verify pods restarted with new image
3. **Hard refresh browser:** `Ctrl + Shift + R` (Windows) or `Cmd + Shift + R` (Mac)
4. Try incognito/private window
5. Check browser console for cached file timestamps

---

## Troubleshooting

### Common Issues and Solutions

#### 1. Pods Stuck in Pending State
```powershell
# Check pod description
kubectl describe pod POD_NAME -n bookstore-dev

# Common causes:
# - Insufficient resources: Scale down replicas
# - PVC not bound: Check if EBS CSI driver is installed
# - Node not ready: Check node status with `kubectl get nodes`
```

#### 2. PostgreSQL PVC Pending
```powershell
# Verify EBS CSI driver is running
kubectl get pods -n kube-system | Select-String "ebs"

# Should show ebs-csi-controller and ebs-csi-node pods Running

# If not installed:
aws eks create-addon --cluster-name my-bookstore --addon-name aws-ebs-csi-driver --region us-east-1
```

#### 3. Pods in CrashLoopBackOff
```powershell
# Check pod logs
kubectl logs POD_NAME -n bookstore-dev

# Common fixes:
# - Database not ready: Wait for postgres-0 to be Running
# - Wrong secrets: Verify secret keys exist
# - Missing ConfigMaps: Apply configmaps directory
```

#### 4. CreateContainerConfigError
```powershell
# Usually means missing secret keys
# Check secret keys
kubectl get secret postgres-credentials -n bookstore-dev -o yaml

# Ensure these keys exist:
# - catalog_password
# - user_password
# - order_password
# - POSTGRES_PASSWORD

# If missing, update the secret file and reapply
```

#### 5. Image Pull Errors
```powershell
# Verify images exist in ECR
aws ecr describe-images --repository-name catalog-service --region us-east-1

# Verify ECR URLs in manifests match your account
# Format: ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/SERVICE_NAME:latest
```

#### 6. LoadBalancer Stuck in Pending
```powershell
# Check service events
kubectl describe svc frontend -n bookstore-dev

# Verify security groups allow traffic
# Check AWS Console → EC2 → Load Balancers
```

#### 7. Too Many Pods Error
```powershell
# Scale down services to 1 replica each
kubectl scale deployment --all --replicas=1 -n bookstore-dev

# Or increase node count in node group
aws eks update-nodegroup-config --cluster-name my-bookstore --nodegroup-name my-bookstore-workers --scaling-config minSize=2,maxSize=4,desiredSize=3
```

### Useful Debugging Commands
```powershell
# Get all resources in namespace
kubectl get all -n bookstore-dev

# Check pod logs
kubectl logs POD_NAME -n bookstore-dev

# Get recent events
kubectl get events -n bookstore-dev --sort-by='.lastTimestamp' | Select-Object -Last 20

# Describe a problematic pod
kubectl describe pod POD_NAME -n bookstore-dev

# Execute command in pod
kubectl exec -it POD_NAME -n bookstore-dev -- sh

# Port forward to test service locally
kubectl port-forward svc/catalog-service 8081:8081 -n bookstore-dev
```

---

## Important AWS Academy Notes

### Session Limitations
- **Lab sessions expire after 4 hours** - All resources are deleted!
- **Temporary credentials** - Session token required
- **Cannot create IAM roles** - Must use existing LabRole
- **Cannot create VPCs** - Must use existing default VPC
- **Budget limits** - ~$100-200 per session

### Before Lab Expires
```powershell
# Save your work!
# 1. Commit and push code changes to GitHub
git add .
git commit -m "EKS deployment configuration"
git push origin deploy/eks-simple

# 2. Document LoadBalancer URL
kubectl get svc frontend -n bookstore-dev -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' > deployment-url.txt

# 3. Export cluster configuration (optional)
kubectl config view > kubeconfig-backup.yaml
```

### After Lab Restart
You'll need to:
1. Reconfigure AWS credentials (new session token)
2. Recreate EKS cluster from scratch (using this guide)
3. Rebuild and push Docker images
4. Redeploy all Kubernetes resources

---

## Cleanup

### Delete Everything
```powershell
# Delete node group first
aws eks delete-nodegroup --cluster-name my-bookstore --nodegroup-name my-bookstore-workers --region us-east-1

# Wait for node group deletion (5-10 minutes)
aws eks describe-nodegroup --cluster-name my-bookstore --nodegroup-name my-bookstore-workers --region us-east-1 --query "nodegroup.status"

# Delete cluster
eksctl delete cluster --name my-bookstore --region us-east-1

# Delete ECR repositories
aws ecr delete-repository --repository-name catalog-service --force --region us-east-1
aws ecr delete-repository --repository-name user-service --force --region us-east-1
aws ecr delete-repository --repository-name cart-service --force --region us-east-1
aws ecr delete-repository --repository-name order-service --force --region us-east-1
aws ecr delete-repository --repository-name api-gateway --force --region us-east-1
aws ecr delete-repository --repository-name customer-frontend --force --region us-east-1
```

---

## Quick Reference: Common Update Commands

### Quick Service Update
```powershell
# Set your service name
$SERVICE = "catalog-service"  # Change this

# Restart deployment (if image already pushed)
kubectl rollout restart deployment/$SERVICE -n bookstore-dev
kubectl rollout status deployment/$SERVICE -n bookstore-dev
```

### Quick Image Rebuild & Deploy
```powershell
# Set variables
$SERVICE = "catalog-service"
$AWS_ACCOUNT = "905418472239"

# Build, push, and deploy
cd services/$SERVICE
docker build -t ${AWS_ACCOUNT}.dkr.ecr.us-east-1.amazonaws.com/${SERVICE}:latest .
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin ${AWS_ACCOUNT}.dkr.ecr.us-east-1.amazonaws.com
docker push ${AWS_ACCOUNT}.dkr.ecr.us-east-1.amazonaws.com/${SERVICE}:latest
kubectl delete pods -l app=$SERVICE -n bookstore-dev
```

### Frontend Update (with cache clear reminder)
```powershell
# Build and deploy
cd frontend/customer-app
docker build --no-cache --build-arg VITE_API_URL="" -t 905418472239.dkr.ecr.us-east-1.amazonaws.com/customer-frontend:latest .
docker push 905418472239.dkr.ecr.us-east-1.amazonaws.com/customer-frontend:latest
kubectl rollout restart deployment/frontend -n bookstore-dev
kubectl rollout status deployment/frontend -n bookstore-dev

# IMPORTANT: Clear browser cache!
# Windows/Linux: Ctrl + Shift + R
# macOS: Cmd + Shift + R
```

### Check Deployment Status
```powershell
# All pods
kubectl get pods -n bookstore-dev

# Specific service
kubectl get pods -l app=$SERVICE -n bookstore-dev

# Watch in real-time
kubectl get pods -n bookstore-dev -w

# Check logs
kubectl logs -l app=$SERVICE -n bookstore-dev --tail=50
kubectl logs -l app=$SERVICE -n bookstore-dev -f  # Follow mode
```

### Rollback Commands
```powershell
# Rollback to previous version
kubectl rollout undo deployment/$SERVICE -n bookstore-dev

# View history
kubectl rollout history deployment/$SERVICE -n bookstore-dev

# Rollback to specific revision
kubectl rollout undo deployment/$SERVICE -n bookstore-dev --to-revision=2
```

### Update ConfigMaps/Secrets
```powershell
# Update ConfigMap
kubectl apply -f infrastructure/k8s/configmaps/
kubectl rollout restart deployment/$SERVICE -n bookstore-dev  # Restart to pick up changes

# Update Secret
kubectl apply -f infrastructure/k8s/secrets/
kubectl rollout restart deployment/$SERVICE -n bookstore-dev
```

---

## Appendix: Key Values Reference

Replace these values throughout the guide with your own:

| Value | Example | Your Value |
|-------|---------|------------|
| AWS Account ID | `905418472239` | ____________ |
| VPC ID | `vpc-0718312e7f365e4b8` | ____________ |
| Subnet 1 (us-east-1a) | `subnet-0b0a497aca8761715` | ____________ |
| Subnet 2 (us-east-1b) | `subnet-08d684911389ce0b4` | ____________ |
| Subnet 3 (us-east-1d) | `subnet-0229448bb3e7df8fd` | ____________ |
| Frontend URL | `aeb72...elb.amazonaws.com` | ____________ |

---

## Success Checklist

- [ ] AWS Academy lab started and credentials configured
- [ ] eksctl and kubectl installed
- [ ] EKS cluster created successfully
- [ ] 2 worker nodes in Ready state
- [ ] EBS CSI driver addon installed
- [ ] All 6 Docker images built and pushed to ECR
- [ ] Kubernetes manifests updated with ECR URLs
- [ ] All Kubernetes resources deployed
- [ ] PostgreSQL and Redis pods running
- [ ] All service pods running (or most of them)
- [ ] Frontend LoadBalancer provisioned
- [ ] Application accessible via browser

---

**Congratulations! Your Distributed Bookstore is now running on AWS EKS!** 🎉

For issues or questions, refer to the Troubleshooting section above.
