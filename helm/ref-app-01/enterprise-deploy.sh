echo "update the GITHUB_TOKEN and then only proceed"
echo GITHUB_TOKEN=$GITHUB_TOKEN
echo ""
echo "Press CR if the above is correct or Ctrl-C to quit"
read x

kubectl delete ns ent-3t-app-01 --ignore-not-found

kubectl create ns ent-3t-app-01
kubectl -n ent-3t-app-01 create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=manojmenon \
  --docker-password=${GITHUB_TOKEN} \
  --docker-email=manojmenon.menon@email.com

kubectl -n ent-3t-app-01 apply -f deploy/k8s/postgres/
kubectl -n ent-3t-app-01 apply -f deploy/k8s/backend/
kubectl -n ent-3t-app-01 apply -f deploy/k8s/frontend/


kubectl rollout status -n ent-3t-app-01 deployment/enterprise-frontend --timeout=300s
kubectl rollout status -n ent-3t-app-01 deployment/enterprise-backend --timeout=300s

sudo ufw allow 8080/tcp
kubectl port-forward -n ent-3t-app-01 svc/enterprise-frontend 8080:80 --address 0.0.0.0
