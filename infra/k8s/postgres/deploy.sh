#!/bin/bash
kubectl apply -f postgres-secret.yaml
kubectl apply -f postgres-init-sql.yaml
kubectl apply -f postgres-pvc.yaml
kubectl apply -f deployment.yaml
kubectl apply -f postgres-service.yaml