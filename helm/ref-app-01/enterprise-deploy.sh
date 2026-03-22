echo "update the GITHUB_TOEKN and then only proceed"
read x

 kubectl delete ns ent-3t-app-01
kubectl create ns ent-3t-app-01
kubectl -n ent-3t-app-01 create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=manojmenon \
  --docker-password=${GITHUB} \
  --docker-email=manojmenon.menon@email.com

kubectl -n ent-3t-app-01 apply -f deploy/k8s/postgres/
kubectl -n ent-3t-app-01 apply -f deploy/k8s/backend/
kubectl -n ent-3t-app-01 apply -f deploy/k8s/frontend/
