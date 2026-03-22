#kubectl create secret generic repo-creds   -n argocd   --from-literal=url=https://github.com/manojmenon/argocd-test.git   --from-literal=username=manojmenon   --from-literal=password=GH_TOKEN
argocd repo add https://github.com/manojmenon/argocd-test.git \
  --username manojmenon \
  --password GH_TOKEN
