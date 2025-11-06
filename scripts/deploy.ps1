# EKS Deployment Script - Complete

# === CLUSTER CREATION ===

# Create EKS cluster
eksctl create cluster -f cluster-config.yaml

# Install EBS CSI Driver addon
aws eks create-addon --cluster-name my-bookstore --addon-name aws-ebs-csi-driver --region us-east-1

# Create node group (replace with YOUR subnets and account ID)
aws eks create-nodegroup `
  --cluster-name my-bookstore `
  --nodegroup-name my-bookstore-workers `
  --node-role arn:aws:iam::905418472239:role/LabRole `
  --subnets subnet-0b0a497aca8761715 subnet-08d684911389ce0b4 subnet-0229448bb3e7df8fd `
  --instance-types t3.medium `
  --scaling-config minSize=2,maxSize=3,desiredSize=3 `
  --disk-size 20 `
  --region us-east-1

# Configure kubectl
aws eks update-kubeconfig --name my-bookstore --region us-east-1

# Verify nodes
kubectl get nodes




# === DEPLOYMENT ===

# Set variables
$ACCOUNT_ID = (aws sts get-caller-identity --query Account --output text)
$ECR_REGISTRY = "$ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com"

# Create namespace
kubectl apply -f .\infrastructure\k8s\namespaces\development.yaml

# Create secrets
kubectl apply -f .\infrastructure\k8s\secrets\

# Deploy databases
kubectl apply -f .\infrastructure\k8s\databases\

# Wait for postgres
kubectl wait --for=condition=ready pod/postgres-0 -n bookstore-dev --timeout=300s

# Deploy ConfigMaps
kubectl apply -f .\infrastructure\k8s\configmaps\

# Deploy services
kubectl apply -f .\infrastructure\k8s\services\catalog-service\
kubectl apply -f .\infrastructure\k8s\services\user-service\
kubectl apply -f .\infrastructure\k8s\services\cart-service\
kubectl apply -f .\infrastructure\k8s\services\order-service\
kubectl apply -f .\infrastructure\k8s\services\inventory-service\
kubectl apply -f .\infrastructure\k8s\services\review-service\
kubectl apply -f .\infrastructure\k8s\services\recommendation-service\
kubectl apply -f .\infrastructure\k8s\services\admin-service\

# Deploy API Gateway
kubectl apply -f .\infrastructure\k8s\services\api-gateway\

# Deploy Frontend
kubectl apply -f .\infrastructure\k8s\frontend\

# Check status
kubectl get pods -n bookstore-dev

# Expose frontend 
kubectl patch svc frontend -n bookstore-dev -p '{\"spec\":{\"type\":\"LoadBalancer\"}}'

# Get frontend URL
kubectl get svc frontend -n bookstore-dev -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
