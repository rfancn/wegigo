kubectl -n wegigo delete all --all
kubectl -n wegigo delete role --all
kubectl -n wegigo delete rolebinding --all
kubectl -n wegigo delete serviceaccount --all
kubectl delete serviceaccount rabbitmq --ignore-not-found=true
kubectl delete namespace wegigo --ignore-not-found=true

