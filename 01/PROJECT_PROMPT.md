# 🚀 Enterprise 3-Tier Application Prompt

Create a production-grade, enterprise-ready 3-layer architecture application with the following requirements:

---

# 🧱 ARCHITECTURE

## Layers:
1. Frontend (React + TypeScript)
2. Backend (Go - Gin or Fiber)
3. Database (3 PostgreSQL instances with logical replication)

---

# 🗄️ DATABASE LAYER

## Setup:
- 3 PostgreSQL instances:
  - postgres-primary-1
  - postgres-primary-2
  - postgres-primary-3

## Requirements:
- Enable logical replication
- Configure bidirectional (multi-master style) replication
- Each DB should replicate to the other two
- Use `wal_level=logical`

## Schema:
Create a `users` table:

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

- BACKEND (GO)

- Framework:
	•	Use Gin (or Fiber)

- Architecture:
  Follow Clean Architecture:

backend/
  cmd/
    server/
      main.go
  internal/
    domain/
      user.go
    repository/
      user_repository.go
    service/
      user_service.go
    handler/
      user_handler.go
  pkg/
    logger/
    config/
    database/
  api/
    routes.go

    Features:
	•	CRUD APIs for users:
	•	POST /users
	•	GET /users
	•	GET /users/:id
	•	PUT /users/:id
	•	DELETE /users/:id

Logging:
	•	Structured logging (Zap or Logrus)
	•	Log:
	•	requests
	•	responses
	•	errors

DB Access:
	•	Use pgx or GORM
	•	Connection pooling
	•	Retry logic

Config:
	•	Environment-based config

⸻

🎨 FRONTEND (REACT)

Stack:
	•	React + TypeScript
	•	Axios for API calls
	•	React Query (optional)

frontend/
  src/
    components/
    pages/
      Users/
    services/
      api.ts
    hooks/
    types/
    utils/

    Features:
	•	Users page:
	•	List users
	•	Create user
	•	Edit user
	•	Delete user

UI:
	•	Simple but clean enterprise UI
	•	Form validation
	•	Error handling

⸻

🐳 DOCKERIZATION

Each service must have:
	•	Dockerfile
	•	Optimized multi-stage builds

Services:
	•	frontend
	•	backend
	•	postgres1
	•	postgres2
	•	postgres3

⸻

🧪 TESTING

Backend:
	•	Unit tests for:
	•	services
	•	repository

Frontend:
	•	Basic component tests (Jest + React Testing Library)

Structure:

tests/
  backend/
  frontend/

k8s support

deploy/k8s/
  frontend/
    deployment.yaml
    service.yaml
  backend/
    deployment.yaml
    service.yaml
  postgres/
    statefulset.yaml
    service.yaml

    Requirements:
	•	Use StatefulSet for PostgreSQL
	•	PersistentVolumeClaims
	•	ConfigMaps for config
	•	Secrets for DB credentials


  Docker compose support

  deploy/docker-compose/
  docker-compose.yml


  Shell script to create the fullp project structure

  scripts/init-project.sh

It should:
	•	Create full directory structure
	•	Initialize git repo
	•	Add .gitignore
	•	Create README.md
	•	Create base Dockerfiles
	•	Create base configs

Final Project structure

project-root/
  frontend/
  backend/
  database/
    init/
    replication/
  deploy/
    docker-compose/
    k8s/
  scripts/
  tests/
  .gitignore
  README.md

  📦 EXTRA REQUIREMENTS
	•	Use environment variables everywhere
	•	Add Makefile for common commands
	•	Include health checks:
	•	/health (backend)
	•	Graceful shutdown in Go
	•	CORS enabled

🎯 OUTPUT EXPECTATION

Generate:
	•	Full working codebase
	•	Ready to run with:
	•	docker-compose up
	•	kubectl apply -f deploy/k8s


⚠️ NOTES
	•	Keep code production-grade
	•	Follow best practices
	•	Add comments where needed
	•	Ensure services can scale independently

---

# 🧰 `scripts/init-project.sh`

```bash
#!/bin/bash

set -e

PROJECT_NAME="enterprise-3tier-app"

echo "Creating project structure..."

mkdir -p $PROJECT_NAME/{frontend,backend,database/init,database/replication}
mkdir -p $PROJECT_NAME/deploy/{docker-compose,k8s}
mkdir -p $PROJECT_NAME/scripts
mkdir -p $PROJECT_NAME/tests/{backend,frontend}

cd $PROJECT_NAME

echo "Initializing git..."
git init

echo "Creating .gitignore..."
cat <<EOF > .gitignore
node_modules/
.env
*.log
dist/
build/
coverage/
*.db
EOF

echo "Creating README..."
cat <<EOF > README.md
# Enterprise 3-Tier Application

## Run with Docker Compose
\`\`\`
docker-compose up --build
\`\`\`

## Run with Kubernetes
\`\`\`
kubectl apply -f deploy/k8s/
\`\`\`
EOF

echo "Creating Makefile..."
cat <<EOF > Makefile
up:
	docker-compose -f deploy/docker-compose/docker-compose.yml up --build

down:
	docker-compose -f deploy/docker-compose/docker-compose.yml down

test:
	go test ./... && npm test --prefix frontend
EOF

echo "Project initialized successfully!"
```

Save the prompt as PROJECT_PROMPT.md
