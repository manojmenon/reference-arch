echo "update the GITHUB_TOEKN and then only proceed"
read x
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=manojmenon \
  --docker-password=GITHUB \
  --docker-email=manojmenon.menon@email.com
