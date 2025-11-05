# Kubernetes Quickstart Guide

## Prerequisites Check

```bash
# Verify installations
docker --version
kubectl version --client
minikube version
```

## Deploy Application

```bash
# Start minikube (if not already running)
minikube start

# Deploy everything
cd infrastructure/k8s
./deploy.sh
```

## Access URLs

After successful deployment, access the application at:

```bash
# Get minikube IP
MINIKUBE_IP=$(minikube ip)

# Frontend
echo "Frontend: http://$MINIKUBE_IP:30000"

# API Gateway
echo "API Gateway: http://$MINIKUBE_IP:30080"
```

**Default URLs (typical minikube IP):**
- Frontend: http://192.168.49.2:30000
- API Gateway: http://192.168.49.2:30080

## Quick Commands

### View Status

```bash
# All pods
kubectl get pods -n bookstore-dev

# All services
kubectl get services -n bookstore-dev

# Check deployments
kubectl get deployments -n bookstore-dev
```

### View Logs

```bash
# API Gateway
kubectl logs -f deployment/api-gateway -n bookstore-dev

# Catalog Service
kubectl logs -f deployment/catalog-service -n bookstore-dev

# Frontend
kubectl logs -f deployment/frontend -n bookstore-dev

# PostgreSQL
kubectl logs -f postgres-0 -n bookstore-dev
```

### Troubleshooting

```bash
# Describe a pod (detailed info)
kubectl describe pod <pod-name> -n bookstore-dev

# Get pod events
kubectl get events -n bookstore-dev --sort-by='.lastTimestamp'

# Check resource usage
kubectl top pods -n bookstore-dev
```

### Clean Up

```bash
# Remove all deployments
cd infrastructure/k8s
./undeploy.sh

# Or delete everything including namespace
kubectl delete namespace bookstore-dev
```

## Expected Pod Status

After successful deployment:

```
NAME                               READY   STATUS    RESTARTS   AGE
api-gateway-xxxxxxxxxx-xxxxx       1/1     Running   0          1m
api-gateway-xxxxxxxxxx-xxxxx       1/1     Running   0          1m
catalog-service-xxxxxxxxxx-xxxxx   1/1     Running   0          2m
catalog-service-xxxxxxxxxx-xxxxx   1/1     Running   0          2m
frontend-xxxxxxxxxx-xxxxx          1/1     Running   0          1m
frontend-xxxxxxxxxx-xxxxx          1/1     Running   0          1m
postgres-0                         1/1     Running   0          3m
```

All pods should be `1/1 Running`.

## Architecture

```
┌─────────────────┐
│    Frontend     │ (NodePort 30000)
│   (React SPA)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  API Gateway    │ (NodePort 30080)
│  (Golang 1.24)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Catalog Service │ (ClusterIP 8081/50051)
│  (Golang 1.21)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   PostgreSQL    │ (Headless Service)
│  (StatefulSet)  │
└─────────────────┘
```

## Testing Endpoints

```bash
# Get minikube IP
MINIKUBE_IP=$(minikube ip)

# Test API Gateway health
curl http://$MINIKUBE_IP:30080/health

# Test catalog endpoints
curl http://$MINIKUBE_IP:30080/api/v1/books

# Access frontend in browser
open http://$MINIKUBE_IP:30000
```

## Rebuilding After Code Changes

```bash
# Configure Docker to use minikube
eval $(minikube docker-env)

# Rebuild specific service
cd services/catalog-service
docker build -t catalog-service:latest .

# Delete pods to force restart with new image
kubectl delete pods -l app=catalog-service -n bookstore-dev
```

## Common Issues

### Pods in ImagePullBackOff
**Solution:** Make sure you built images with minikube's Docker daemon
```bash
eval $(minikube docker-env)
./deploy.sh
```

### Pods in CrashLoopBackOff
**Solution:** Check logs for errors
```bash
kubectl logs <pod-name> -n bookstore-dev
kubectl describe pod <pod-name> -n bookstore-dev
```

### Can't access services
**Solution:** Verify minikube IP and ports
```bash
minikube ip
kubectl get services -n bookstore-dev
```

## Next Steps

1. Add more services (User Service, Payment Service, etc.)
2. Set up Ingress Controller
3. Add monitoring (Prometheus + Grafana)
4. Add distributed tracing (Jaeger)
5. Implement HorizontalPodAutoscaler
6. Add network policies for security

## Support

For detailed documentation, see [README.md](./README.md)
