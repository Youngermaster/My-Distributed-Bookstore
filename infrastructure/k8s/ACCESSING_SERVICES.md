# Accessing Services on Minikube (macOS with Docker Driver)

## Problem

When using Minikube with the Docker driver on macOS, the cluster IP (e.g., `192.168.49.2`) is **NOT accessible** from your host machine or local network. This is because the Minikube cluster runs inside a Docker container.

## Your Services Status ✅

All services are running perfectly:
- **Frontend**: Running (2 replicas) - NodePort 30000
- **API Gateway**: Running (2 replicas) - NodePort 30080
- **All backend services**: Healthy and responding

The issue is **only** with network access from your machine.

---

## Solution 1: Use `minikube service` (RECOMMENDED)

**Best for**: Quick testing and development

### Commands

Open **two separate terminal windows**:

```bash
# Terminal 1 - Frontend
minikube service frontend -n bookstore-dev

# Terminal 2 - API Gateway
minikube service api-gateway -n bookstore-dev
```

### What happens
- Creates a temporary tunnel to the service
- Automatically opens your browser
- Shows accessible URL (e.g., `http://127.0.0.1:54321`)
- Keep terminals open while using the app

### Pros
- Easiest solution
- No configuration needed
- Works immediately

### Cons
- Need to keep terminal windows open
- URLs change each time
- One terminal per service

---

## Solution 2: Use `kubectl port-forward` (SIMPLE)

**Best for**: Consistent localhost URLs

### Commands

Open **two separate terminal windows**:

```bash
# Terminal 1 - Forward Frontend to localhost:3000
kubectl port-forward -n bookstore-dev service/frontend 3000:80

# Terminal 2 - Forward API Gateway to localhost:8080
kubectl port-forward -n bookstore-dev service/api-gateway 8080:8080
```

### Access
- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080

### Pros
- Consistent URLs (always localhost)
- Easy to remember ports
- No sudo password needed

### Cons
- Need to keep terminal windows open
- One terminal per service
- Only accessible from localhost (not from other devices on network)

---

## Solution 3: Use `minikube tunnel` (PERMANENT)

**Best for**: Long-term development, accessing from other devices on network

### Command

Open **one terminal window**:

```bash
minikube tunnel
# Enter your sudo password when prompted
```

**Keep this terminal open!**

### Access

Once tunnel is running:
- **Frontend**: http://192.168.49.2:30000
- **API Gateway**: http://192.168.49.2:30080

### For LoadBalancer Support

To make this work better, update your service types from `NodePort` to `LoadBalancer`:

```bash
# Edit frontend service
kubectl patch service frontend -n bookstore-dev -p '{"spec":{"type":"LoadBalancer"}}'

# Edit API gateway service
kubectl patch service api-gateway -n bookstore-dev -p '{"spec":{"type":"LoadBalancer"}}'

# Get LoadBalancer IPs (may take a moment)
kubectl get services -n bookstore-dev
```

### Pros
- Works like production LoadBalancer
- Can access from other devices on local network
- Most realistic for testing

### Cons
- Requires sudo password
- Must keep terminal open
- May have issues if network changes

---

## Solution 4: Change Service Types to LoadBalancer (PRODUCTION-LIKE)

**Best for**: Testing production-like setup

### Steps

1. Update service definitions:

```bash
# Patch frontend
kubectl patch service frontend -n bookstore-dev -p '{"spec":{"type":"LoadBalancer"}}'

# Patch API gateway
kubectl patch service api-gateway -n bookstore-dev -p '{"spec":{"type":"LoadBalancer"}}'
```

2. Start tunnel (in separate terminal):

```bash
minikube tunnel
```

3. Get external IPs:

```bash
kubectl get services -n bookstore-dev
```

4. Access services via the EXTERNAL-IP shown

### To make permanent, update your service manifests

Edit `infrastructure/k8s/services/frontend.yaml`:

```yaml
spec:
  type: LoadBalancer  # Change from NodePort
  ports:
  - port: 80
    targetPort: 80
    # Remove nodePort: 30000
```

Edit `infrastructure/k8s/services/api-gateway.yaml`:

```yaml
spec:
  type: LoadBalancer  # Change from NodePort
  ports:
  - port: 8080
    targetPort: 8080
    # Remove nodePort: 30080
```

Then redeploy:

```bash
kubectl apply -f infrastructure/k8s/services/frontend.yaml
kubectl apply -f infrastructure/k8s/services/api-gateway.yaml
```

---

## Recommended Workflow

### For Quick Testing (Right Now)

Use **Solution 2** (port-forward):

```bash
# Terminal 1
kubectl port-forward -n bookstore-dev service/frontend 3000:80

# Terminal 2
kubectl port-forward -n bookstore-dev service/api-gateway 8080:8080

# Access at http://localhost:3000
```

### For Regular Development

Use **Solution 3** (minikube tunnel):

```bash
# Start once and leave running
minikube tunnel

# Access at http://192.168.49.2:30000
```

### For Testing from Phone/Other Devices

Use **Solution 4** (LoadBalancer + tunnel):

```bash
# Patch services to LoadBalancer
kubectl patch service frontend -n bookstore-dev -p '{"spec":{"type":"LoadBalancer"}}'
kubectl patch service api-gateway -n bookstore-dev -p '{"spec":{"type":"LoadBalancer"}}'

# Start tunnel
minikube tunnel

# Get the external IP
kubectl get services -n bookstore-dev

# Access from any device on your network using the EXTERNAL-IP
```

---

## Troubleshooting

### "Connection refused" or "This site can't be reached"

**Check if tunnel/port-forward is running:**

```bash
# Check if port-forward is active
lsof -i :3000
lsof -i :8080

# Check if tunnel is active
ps aux | grep "minikube tunnel"
```

### "Tunnel is already running"

```bash
# Find and kill existing tunnel
ps aux | grep "minikube tunnel" | grep -v grep | awk '{print $2}' | xargs kill

# Start new tunnel
minikube tunnel
```

### Services not responding

```bash
# Check pod status
kubectl get pods -n bookstore-dev

# Check pod logs
kubectl logs -n bookstore-dev deployment/frontend --tail=50
kubectl logs -n bookstore-dev deployment/api-gateway --tail=50

# Restart pods if needed
kubectl rollout restart deployment/frontend -n bookstore-dev
kubectl rollout restart deployment/api-gateway -n bookstore-dev
```

### Need to access from specific port

```bash
# Forward to any port you want
kubectl port-forward -n bookstore-dev service/frontend 8888:80

# Access at http://localhost:8888
```

---

## Summary

| Solution | Ease | Persistent | Network Access | Terminals Needed |
|----------|------|------------|----------------|------------------|
| minikube service | ⭐⭐⭐ | ❌ Random URL | ❌ Localhost only | 2 (one per service) |
| port-forward | ⭐⭐⭐ | ✅ Same URL | ❌ Localhost only | 2 (one per service) |
| minikube tunnel | ⭐⭐ | ✅ Same IP | ✅ Local network | 1 (sudo needed) |
| LoadBalancer + tunnel | ⭐ | ✅ External IP | ✅ Local network | 1 (sudo needed) |

**Our recommendation**: Start with **port-forward** (Solution 2) for immediate testing, then move to **minikube tunnel** (Solution 3) for regular development.

---

## Quick Start (Copy-Paste)

```bash
# Option A: Port-forward (easiest, immediate access)
# Terminal 1:
kubectl port-forward -n bookstore-dev service/frontend 3000:80

# Terminal 2:
kubectl port-forward -n bookstore-dev service/api-gateway 8080:8080

# Access: http://localhost:3000

# Option B: Minikube tunnel (better for development)
# Terminal 1:
minikube tunnel
# Enter password when prompted

# Access: http://192.168.49.2:30000
```

---

## Why This Happens

- **Docker Driver Limitation**: Minikube with Docker driver creates an isolated network
- **macOS Networking**: Docker on Mac doesn't expose container networks to host
- **NodePort Limitation**: NodePort services in Minikube need special access methods
- **Production vs Development**: In production K8s (EKS, GKE), LoadBalancer type works automatically

This is a **known limitation** and affects all Minikube users on macOS with Docker driver. The solutions above are the standard workarounds recommended by Minikube documentation.
