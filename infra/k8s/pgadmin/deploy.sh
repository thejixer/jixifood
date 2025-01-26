#!/bin/bash
kubectl apply -f pgadmin-secret.yaml
kubectl apply -f pgadmin-deployment.yaml
kubectl apply -f pgadmin-service.yaml