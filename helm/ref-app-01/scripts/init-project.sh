#!/usr/bin/env bash
set -euo pipefail

PROJECT_NAME="${1:-enterprise-3tier-app}"

echo "Creating project structure under ${PROJECT_NAME}..."

mkdir -p "${PROJECT_NAME}"/{frontend,backend,database/init,database/replication}
mkdir -p "${PROJECT_NAME}"/deploy/{docker-compose,k8s/{frontend,backend,postgres}}
mkdir -p "${PROJECT_NAME}"/scripts
mkdir -p "${PROJECT_NAME}"/tests/{backend,frontend}

cd "${PROJECT_NAME}"

if [[ ! -d .git ]]; then
  echo "Initializing git..."
  git init
fi

echo "Creating .gitignore..."
cat <<'EOF' > .gitignore
node_modules/
.env
.env.local
*.log
dist/
build/
coverage/
*.db
backend/bin/
.idea/
.vscode/
EOF

echo "Creating README..."
cat <<'EOF' > README.md
# Enterprise 3-Tier Application

See PROJECT_PROMPT.md for full requirements.

## Docker Compose

```bash
docker compose -f deploy/docker-compose/docker-compose.yml up --build
```

## Kubernetes

```bash
kubectl apply -f deploy/k8s/postgres/
kubectl apply -f deploy/k8s/backend/
kubectl apply -f deploy/k8s/frontend/
```
EOF

echo "Creating Makefile..."
cat <<'EOF' > Makefile
.PHONY: up down test

up:
	docker compose -f deploy/docker-compose/docker-compose.yml up --build

down:
	docker compose -f deploy/docker-compose/docker-compose.yml down

test:
	cd backend && go test ./...
	cd frontend && npm test
EOF

echo "Creating base Dockerfiles (placeholder)..."
mkdir -p backend frontend
cat <<'EOF' > backend/Dockerfile
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /server ./cmd/server
FROM alpine:3.20
COPY --from=build /server /server
ENTRYPOINT ["/server"]
EOF

cat <<'EOF' > frontend/Dockerfile
FROM node:20-alpine AS build
WORKDIR /app
COPY package.json ./
RUN npm install
COPY . .
RUN npm run build
FROM nginx:1.27-alpine
COPY --from=build /app/dist /usr/share/nginx/html
EOF

echo "Project initialized. Copy source from the reference repository or expand these stubs."
echo "Done."
