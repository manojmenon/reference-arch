# Enterprise 3-Tier Application

Production-style reference stack: React + TypeScript frontend, Go (Gin) API with clean architecture, and three PostgreSQL 16 instances where **postgres-primary-1** is the write target and **logical replication fans out** to nodes 2 and 3 (same data, read replicas). A full multi-master mesh is not supported by vanilla PostgreSQL logical replication without conflict handling—duplicate paths cause primary-key violations.

## Prerequisites

- Docker and Docker Compose v2
- Go 1.23+ (local backend tests / run)
- Node.js 20+ and npm (local frontend dev / tests)

## Run with Docker Compose

From the repository root:

```bash
make up
# or
docker compose -f deploy/docker-compose/docker-compose.yml up --build
```

- Frontend (nginx + static build): [http://localhost:3000](http://localhost:3000) by default (`FRONTEND_HOST_PORT` in Compose overrides the host port) — API is proxied under `/api` to the backend.
- Backend JSON API: [http://localhost:8080](http://localhost:8080) (direct)
- PostgreSQL: `localhost:5432`, `5433`, `5434` mapped to the three instances (5432 = write primary; 5433/5434 receive one-way replicated data from 5432).

The `replication-bootstrap` one-shot container configures publications and subscriptions after all databases are healthy. If you change init SQL or replication bootstrap logic, reset volumes: `docker compose ... down -v` so old subscriptions/slots are not reused.

## Run locally (development)

**Backend**

```bash
export DATABASE_URL="postgres://appuser:apppass@localhost:5432/appdb?sslmode=disable"
cd backend && go run ./cmd/server
```

**Frontend**

```bash
cd frontend && npm install && npm run dev
```

Vite proxies `/api` to `http://127.0.0.1:8080` by default (override with `VITE_DEV_API`).

## Run with Kubernetes

Build images locally and load them into your cluster (example tags used in manifests):

```bash
docker build -t enterprise-3tier-backend:latest ./backend
docker build -t enterprise-3tier-frontend:latest ./frontend
# kind / minikube / k3s image load steps vary by distro
```

Apply manifests (order: secrets/config, Postgres, then app tier):

```bash
kubectl apply -f deploy/k8s/postgres/
kubectl apply -f deploy/k8s/backend/
kubectl apply -f deploy/k8s/frontend/
```

The StatefulSet runs three PostgreSQL replicas with logical WAL settings; the sample `DATABASE_URL` in `deploy/k8s/postgres/secret.yaml` targets `postgres-primary-0`. Peer replication between all three pods is orchestrated in Docker Compose; for Kubernetes you would add a Job similar to `database/replication/bootstrap.sh` with stable DNS names for each pod.

## Tests

```bash
make test
```

- Backend: `go test ./...` under `backend/`
- Frontend: Vitest + Testing Library (`npm test` under `frontend/`)

## Project layout

- `backend/` — Go API (`/users`, `/health`)
- `frontend/` — Vite + React UI
- `database/init` — schema and replication role (used by Compose and k8s ConfigMap)
- `database/replication/bootstrap.sh` — logical replication wiring for Compose
- `deploy/docker-compose/` — full stack
- `deploy/k8s/` — Kubernetes manifests
- `scripts/init-project.sh` — scaffold a fresh copy of the tree
- `PROJECT_PROMPT.md` — original specification

## Environment variables (backend)

| Variable | Description |
|----------|-------------|
| `PORT` | HTTP listen port (default `8080`) |
| `DATABASE_URL` | PostgreSQL DSN for `pgxpool` |
| `LOG_LEVEL` | `debug`, `info`, `warn`, `error` |
| `CORS_ORIGIN` | `Access-Control-Allow-Origin` (default `*`) |
| `SHUTDOWN_TIMEOUT` | Graceful shutdown window |
| `DB_MAX_RETRIES` / `DB_RETRY_BACKOFF` | Transient DB retry policy |

## License

Reference / educational use.
