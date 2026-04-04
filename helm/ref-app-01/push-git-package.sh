echo "Do a docker login first"
 docker login ghcr.io -u manojmenon 

docker build -t enterprise-3tier-backend:latest ./backend
docker build -t enterprise-3tier-frontend:latest ./frontend

docker tag enterprise-3tier-backend:latest ghcr.io/manojmenon/enterprise-3tier-backend:1.0
docker tag enterprise-3tier-frontend:latest ghcr.io/manojmenon/enterprise-3tier-frontend:1.0

docker push ghcr.io/manojmenon/enterprise-3tier-backend:1.0
docker push ghcr.io/manojmenon/enterprise-3tier-frontend:1.0
