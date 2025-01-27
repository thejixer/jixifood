#!/bin/bash
# Install NGINX Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml

kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

# Apply the Ingress resource
kubectl apply -f ingress.yaml