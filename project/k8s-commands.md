https://minikube.sigs.k8s.io/docs/start/
https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/
kubectl apply -f k8s
kubectl get pods
kubectl get svc
kubectl get deployments
delete deployments broker-service mongo rabbitmq
kubectl delete svc broker-service mongo rabbitmq
To reach the services from outside, frontend in this case also runs from outside, use port 8000 then. 
expose service as load balancer -> kubectl expose deployment broker-service --type=LoadBalancer --port=8080 --target-port=8080
minikube tunnel
