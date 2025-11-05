# Kubernetes Deployment Guide

This guide explains how to deploy the Distributed Bookstore application to a local Kubernetes cluster using Minikube.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Docker**: Container runtime
- **kubectl**: Kubernetes command-line tool
- **minikube**: Local Kubernetes cluster

### Verify Prerequisites

```bash
# Check Docker
docker --version

# Check kubectl
kubectl version --client

# Check minikube
minikube version
```

## Quick Start

### 1. Start Minikube

```bash
minikube start
```

### 2. Deploy the Application

From the `infrastructure/k8s` directory, run:

```bash
./deploy.sh
```

This script will:

1. Configure Docker to use Minikube's Docker daemon
2. Build all Docker images (Catalog, User, Cart, Order, Recommendation, API Gateway, Frontend)
3. Create the Kubernetes namespace
4. Deploy secrets and configmaps
5. Deploy PostgreSQL and Redis dependencies
6. Deploy all microservices
7. Display access information

### 3. Access the Application

After deployment completes, the script will display URLs to access the application:

```
Frontend:    http://<MINIKUBE_IP>:30000
API Gateway: http://<MINIKUBE_IP>:30080
```

You can also get the Minikube IP manually:

```bash
minikube ip
```

Then access:

- **Frontend**: `http://<MINIKUBE_IP>:30000`
- **API Gateway**: `http://<MINIKUBE_IP>:30080`

### 4. Verify Deployment

Check the status of all pods:

```bash
kubectl get pods -n bookstore-dev
```

Expected output (example):

```
NAME                            READY   STATUS    RESTARTS   AGE
api-gateway-xxxxxxxxxx-xxxxx    1/1     Running   0          2m
cart-service-xxxxxxxxxx-xxxxx   1/1     Running   0          2m
catalog-service-xxxxxxxxxx-xxxx 1/1     Running   0          3m
frontend-xxxxxxxxxx-xxxxx       1/1     Running   0          1m
order-service-xxxxxxxxxx-xxxxx  1/1     Running   0          2m
recommendation-service-xxxxx    1/1     Running   0          2m
user-service-xxxxxxxxxx-xxxxx   1/1     Running   0          2m
postgres-0                      1/1     Running   0          4m
redis-xxxxxxxxxx-xxxxx          1/1     Running   0          2m
```

Check services:

```bash
kubectl get services -n bookstore-dev
```

### 5. Run Smoke Tests (optional)

After the deployment is up, you can execute the automated smoke checks against the API Gateway:

```bash
./smoke-tests.sh          # Uses the Minikube IP automatically
# or
./smoke-tests.sh http://<gateway-host>:<port>
```

These tests hit the health endpoints, exercise user registration, cart and order flows, and verify recommendation routes.

## Manual Deployment Steps

If you prefer to deploy manually instead of using the script:

### Step 1: Configure Docker Environment

```bash
# Point your shell to minikube's docker daemon
eval $(minikube docker-env)
```

### Step 2: Build Docker Images

```bash
# Build Catalog Service
cd ../../services/catalog-service
docker build -t catalog-service:latest .

# Build API Gateway
cd ../api-gateway
docker build -t api-gateway:latest .

# Build Frontend
cd ../../frontend/customer-app
docker build -t frontend:latest .

# Return to k8s directory
cd ../../infrastructure/k8s
```

### Step 3: Create Namespace

```bash
kubectl apply -f namespaces/development.yaml
```

### Step 4: Create Secrets and ConfigMaps

```bash
kubectl apply -f secrets/
kubectl apply -f configmaps/
```

### Step 5: Deploy PostgreSQL

```bash
kubectl apply -f databases/
```

Wait for PostgreSQL to be ready:

```bash
kubectl wait --for=condition=ready pod -l app=postgres -n bookstore-dev --timeout=300s
```

### Step 6: Deploy Catalog Service

```bash
kubectl apply -f services/catalog-service/
```

Wait for Catalog Service to be ready:

```bash
kubectl wait --for=condition=ready pod -l app=catalog-service -n bookstore-dev --timeout=300s
```

### Step 7: Deploy API Gateway

```bash
kubectl apply -f services/api-gateway/
```

Wait for API Gateway to be ready:

```bash
kubectl wait --for=condition=ready pod -l app=api-gateway -n bookstore-dev --timeout=300s
```

### Step 8: Deploy Frontend

```bash
kubectl apply -f frontend/
```

## Useful Commands

### View Logs

```bash
# View logs for a specific pod
kubectl logs -f <pod-name> -n bookstore-dev

# View logs for all pods in a deployment
kubectl logs -f deployment/catalog-service -n bookstore-dev

# View logs for all containers in the namespace
kubectl logs -f -l app=api-gateway -n bookstore-dev
```

### Describe Resources

```bash
# Describe a pod (useful for debugging)
kubectl describe pod <pod-name> -n bookstore-dev

# Describe a service
kubectl describe service catalog-service -n bookstore-dev

# Describe a deployment
kubectl describe deployment api-gateway -n bookstore-dev
```

### Scale Deployments

```bash
# Scale catalog service to 3 replicas
kubectl scale deployment catalog-service --replicas=3 -n bookstore-dev

# Scale API Gateway to 1 replica
kubectl scale deployment api-gateway --replicas=1 -n bookstore-dev
```

### Execute Commands in Pods

```bash
# Open a shell in a pod
kubectl exec -it <pod-name> -n bookstore-dev -- /bin/sh

# Run a specific command
kubectl exec <pod-name> -n bookstore-dev -- wget -O- http://localhost:8081/health
```

### Port Forwarding

If you want to access services without using NodePort:

```bash
# Forward API Gateway to local port 8080
kubectl port-forward service/api-gateway 8080:8080 -n bookstore-dev

# Forward Frontend to local port 3000
kubectl port-forward service/frontend 3000:80 -n bookstore-dev

# Forward PostgreSQL to local port 5432
kubectl port-forward service/postgres 5432:5432 -n bookstore-dev
```

### View Resource Usage

```bash
# View resource usage by pods
kubectl top pods -n bookstore-dev

# View resource usage by nodes
kubectl top nodes
```

## Undeploying the Application

To remove all deployed resources:

```bash
./undeploy.sh
```

Or manually:

```bash
kubectl delete -f frontend/
kubectl delete -f services/api-gateway/
kubectl delete -f services/catalog-service/
kubectl delete -f databases/
kubectl delete -f configmaps/
kubectl delete -f secrets/
kubectl delete namespace bookstore-dev
```

## Troubleshooting

### Pods Not Starting

Check pod status and events:

```bash
kubectl get pods -n bookstore-dev
kubectl describe pod <pod-name> -n bookstore-dev
```

Common issues:

- **ImagePullBackOff**: Image not found. Make sure you've built the images using minikube's Docker daemon (`eval $(minikube docker-env)`)
- **CrashLoopBackOff**: Container is crashing. Check logs with `kubectl logs <pod-name> -n bookstore-dev`

### Database Connection Issues

Check if PostgreSQL is running:

```bash
kubectl get pods -l app=postgres -n bookstore-dev
```

Check PostgreSQL logs:

```bash
kubectl logs -f postgres-0 -n bookstore-dev
```

Test database connection from catalog service:

```bash
kubectl exec -it <catalog-service-pod> -n bookstore-dev -- /bin/sh
# Inside the pod
wget -O- http://localhost:8081/health
```

### Service Not Accessible

Check service endpoints:

```bash
kubectl get endpoints -n bookstore-dev
```

Check if pods are ready:

```bash
kubectl get pods -n bookstore-dev
```

### Rebuilding Images

If you make changes to the code, rebuild the images:

```bash
# Make sure you're using minikube's Docker daemon
eval $(minikube docker-env)

# Rebuild the image
cd services/catalog-service
docker build -t catalog-service:latest .

# Delete the pods to force them to restart with the new image
kubectl delete pods -l app=catalog-service -n bookstore-dev
```

### Resetting Everything

If you want to start fresh:

```bash
# Delete everything
./undeploy.sh
kubectl delete namespace bookstore-dev

# Or reset minikube completely
minikube delete
minikube start
```

## Architecture Overview

The deployment consists of:

1. **PostgreSQL StatefulSet** (1 replica)

   - Persistent storage for catalog data
   - Service: `postgres:5432`

2. **Catalog Service Deployment** (2 replicas)

   - HTTP API on port 8081
   - gRPC API on port 50051
   - Service: `catalog-service:8081`, `catalog-service:50051`

3. **API Gateway Deployment** (2 replicas)

   - HTTP API on port 8080
   - Exposed via NodePort 30080
   - Service: `api-gateway:8080`

4. **Frontend Deployment** (2 replicas)
   - Static React app served by Nginx
   - Exposed via NodePort 30000
   - Service: `frontend:80`

## Configuration

### Secrets

Located in `secrets/`:

- `postgres-credentials.yaml`: Database credentials
- `jwt-secret.yaml`: JWT signing secret

**Note**: In production, use proper secret management (e.g., Sealed Secrets, Vault)

### ConfigMaps

Located in `configmaps/`:

- `api-gateway-config.yaml`: API Gateway configuration
- `catalog-service-config.yaml`: Catalog Service configuration
- `frontend-config.yaml`: Frontend configuration

### Resource Limits

Current resource limits:

**PostgreSQL**:

- Requests: 256Mi memory, 250m CPU
- Limits: 512Mi memory, 500m CPU

**Catalog Service**:

- Requests: 128Mi memory, 100m CPU
- Limits: 512Mi memory, 500m CPU

**API Gateway**:

- Requests: 128Mi memory, 100m CPU
- Limits: 512Mi memory, 500m CPU

**Frontend**:

- Requests: 64Mi memory, 50m CPU
- Limits: 256Mi memory, 200m CPU

You can adjust these in the respective deployment YAML files.

## Next Steps

1. **Add Ingress**: Set up Nginx Ingress Controller for better routing
2. **Add Redis**: Deploy Redis for caching
3. **Add RabbitMQ**: Deploy RabbitMQ for message queuing
4. **Add Monitoring**: Deploy Prometheus and Grafana
5. **Add Tracing**: Deploy Jaeger for distributed tracing
6. **Add HPA**: Configure Horizontal Pod Autoscaling
7. **Add Network Policies**: Implement security policies

## Production Considerations

When deploying to production:

1. **Use proper secrets management** (not plaintext in YAML)
2. **Configure resource requests and limits** based on actual usage
3. **Set up persistent volume backups** for PostgreSQL
4. **Configure ingress with SSL/TLS** certificates
5. **Implement network policies** for security
6. **Set up monitoring and alerting**
7. **Configure health checks** appropriately
8. **Use a proper image registry** (not local images)
9. **Implement CI/CD** for automated deployments
10. **Configure pod disruption budgets** for high availability
