kubectl create namespace wegigo
kubectl create serviceaccount rabbitmq
kubectl create -f rabbitmq-rbac.yaml
kubectl create -f rabbitmq-service.yaml
kubectl create -f rabbitmq-statefulset.yaml
