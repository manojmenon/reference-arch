# Enterprise 3-Tier Application

Production-style reference stack: **React + TypeScript** frontend, **Go (Gin)** API with clean architecture, and **three PostgreSQL 16** instances. **postgres-primary-1** is the sole write target; **logical replication** fans out to the other two nodes (same dataset, read-replica style). Vanilla PostgreSQL logical replication does not provide a safe multi-master mesh without conflict resolution—duplicate write paths would cause primary-key violations—so this project implements a **primary + replicas** topology, not bidirectional multi-master replication.

---

## Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Stack](#stack)
- [Prerequisites](#prerequisites)
- [Quick start (Docker Compose)](#quick-start-docker-compose)
- [Run locally (development)](#run-locally-development)
- [Kubernetes](#kubernetes)
- [HTTP API](#http-api)
- [Tests](#tests)
- [Project layout](#project-layout)
- [Environment variables (backend)](#environment-variables-backend)
- [Operations & troubleshooting](#operations--troubleshooting)
- [Specification note](#specification-note)
- [License](#license)

---

## Overview

| Layer    | Technology |
|----------|------------|
| UI       | Vite, React, TypeScript; nginx serves the production build; `/api` is proxied to the backend |
| API      | Gin, Zap logging, pgx connection pool with retries, graceful shutdown, CORS |
| Data     | PostgreSQL 16 with `wal_level=logical`; init SQL under `database/init`; replication bootstrap under `database/replication` |

The sample app exposes **user CRUD** (`/users`) and **`/health`**. The frontend includes a Users page with list/create/edit/delete and basic validation.

---

## Architecture

```mermaid
flowchart LR
  subgraph client [Browser]
    UI[React SPA]
  end
  subgraph compose [Docker Compose / K8s]
    FE[nginx + static assets]
    BE[Go API]
    P[(postgres-primary-1\nwrites)]
    R2[(postgres-primary-2\nreplica)]
    R3[(postgres-primary-3\nreplica)]
  end
  UI -->|"/api/*"| FE
  FE --> BE
  BE --> P
  P -->|logical replication| R2
  P -->|logical replication| R3
```

- **Writes** go to **postgres-primary-1** only (see `DATABASE_URL` in Compose and k8s secrets).
- **replication-bootstrap** is a one-shot job that creates publications/subscriptions after all databases are healthy.
- After changing init SQL or bootstrap logic, **remove volumes** (`docker compose ... down -v`) so old replication slots/subscriptions are not reused incorrectly.

---

## Stack

- **Frontend:** `frontend/` — Vite, React, Vitest + Testing Library.
- **Backend:** `backend/` — `cmd/server`, `internal/` (domain, repository, service, handler), `pkg/` (config, database, logger), `api/routes.go`.
- **Database:** `database/init` — schema and replication role; `database/replication/bootstrap.sh` — Compose replication wiring.
- **Deploy:** `deploy/docker-compose/docker-compose.yml` — full stack; `deploy/k8s/` — StatefulSet, services, backend/frontend Deployments.

---

## Prerequisites

- Docker and Docker Compose v2
- Go 1.23+ (local backend tests / run)
- Node.js 20+ and npm (local frontend dev / tests)

---

## Quick start (Docker Compose)

From **`01/`** (paths in `docker-compose.yml` are relative to this directory):

```bash
cd 01
make up
# or
docker compose -f deploy/docker-compose/docker-compose.yml up --build
```

| Service | URL / endpoint |
|---------|----------------|
| Frontend (nginx + static build) | [http://localhost:3000](http://localhost:3000) — default host port (`FRONTEND_HOST_PORT` overrides) — API proxied under `/api` |
| Backend (direct) | [http://localhost:8080](http://localhost:8080) |
| PostgreSQL instances | `localhost:5432`, `5433`, `5434` → containers **5432** (5432 = write primary; 5433/5434 receive one-way replicated data from 5432) |

**Teardown (including volumes, when replication state must be reset):**

```bash
cd 01
make down
docker compose -f deploy/docker-compose/docker-compose.yml down -v
```

---

## Run locally (development)

**Backend**

```bash
export DATABASE_URL="postgres://appuser:apppass@localhost:5432/appdb?sslmode=disable"
cd 01/backend && go run ./cmd/server
```

**Frontend**

```bash
cd 01/frontend && npm install && npm run dev
```

Vite proxies `/api` to `http://127.0.0.1:8080` by default (override with `VITE_DEV_API`).

---

## Kubernetes

Build images (example tags used in manifests):

```bash
cd 01
docker build -t enterprise-3tier-backend:latest ./backend
docker build -t enterprise-3tier-frontend:latest ./frontend
# kind / minikube / k3s image load steps vary by distro
```

Apply manifests (order: secrets/config, Postgres, then app tier):

```bash
cd 01
kubectl apply -f deploy/k8s/postgres/
kubectl apply -f deploy/k8s/backend/
kubectl apply -f deploy/k8s/frontend/
```

The StatefulSet runs three PostgreSQL replicas with logical WAL settings; the sample `DATABASE_URL` in `deploy/k8s/postgres/secret.yaml` targets `postgres-primary-0`. Peer replication in the Compose file is orchestrated for local dev; for Kubernetes you would add a **Job** similar to `database/replication/bootstrap.sh` using stable DNS names for each pod.

---

## HTTP API

Base URL: `http://localhost:8080` (or `/api` via the frontend nginx proxy, depending on deployment).

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Liveness-style JSON: `{"status":"ok"}` |
| `POST` | `/users` | Create user |
| `GET` | `/users` | List users |
| `GET` | `/users/:id` | Get user by ID |
| `PUT` | `/users/:id` | Update user |
| `DELETE` | `/users/:id` | Delete user |

---

## Tests

```bash
cd 01
make test
```

- **Backend:** `go test ./...` under `backend/`
- **Frontend:** Vitest + Testing Library (`npm test` under `frontend/`)

Individual targets: `make backend-test`, `make frontend-test`.

---

## Project layout

| Path | Purpose |
|------|---------|
| `backend/` | Go API (`/users`, `/health`) |
| `frontend/` | Vite + React UI |
| `database/init` | Schema and replication role (Compose + k8s ConfigMap) |
| `database/replication/bootstrap.sh` | Logical replication wiring for Compose |
| `deploy/docker-compose/` | Full stack |
| `deploy/k8s/` | Kubernetes manifests |
| `scripts/init-project.sh` | Scaffold a fresh copy of the tree |
| `PROJECT_PROMPT.md` | Original specification |

---

## Environment variables (backend)

| Variable | Description |
|----------|-------------|
| `PORT` | HTTP listen port (default `8080`) |
| `DATABASE_URL` | PostgreSQL DSN for `pgxpool` |
| `LOG_LEVEL` | `debug`, `info`, `warn`, `error` |
| `CORS_ORIGIN` | `Access-Control-Allow-Origin` (default `*`) |
| `SHUTDOWN_TIMEOUT` | Graceful shutdown window |
| `DB_MAX_RETRIES` / `DB_RETRY_BACKOFF` | Transient DB retry policy |

---

## Operations & troubleshooting

- **Replication bootstrap failed after schema changes:** Run `docker compose ... down -v`, then `up --build` so volumes and slots match the new init SQL.
- **Backend cannot connect:** Ensure postgres-primary-1 is healthy and `DATABASE_URL` points at the primary, not a replica, for writes.
- **Frontend 502 on `/api`:** Confirm the backend container is up and nginx upstream matches `enterprise-backend:8080` in the frontend image config.

---

## Specification note

`PROJECT_PROMPT.md` originally asked for bidirectional multi-master replication between all three databases. This repository implements **one primary and two logical replicas**, which matches PostgreSQL’s typical replication model and avoids conflicting writes. See [Overview](#overview) for the rationale.

---

## License

Reference / educational use.
