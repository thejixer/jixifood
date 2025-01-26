at this stage, since I'm not a devops engineer, I don't put much effort in CI/CD 
and the deployment is semi-manual

first we need to create the namespace
```
cd k8s
kubectl apply -f namespace.yaml
```

and then run the necessary deployments one by one

```
cd postgres
./deploy.sh
cd ..
```
```
cd pgadmin
./deploy.sh
cd ..
```
```
cd redis
./deploy.sh
cd ..
```
