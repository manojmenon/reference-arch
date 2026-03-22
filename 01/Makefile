.PHONY: up down test backend-test frontend-test docker-build lint

up:
	docker compose -f deploy/docker-compose/docker-compose.yml up --build

down:
	docker compose -f deploy/docker-compose/docker-compose.yml down

docker-build:
	docker compose -f deploy/docker-compose/docker-compose.yml build

test: backend-test frontend-test

backend-test:
	cd backend && go test ./...

frontend-test:
	cd frontend && npm test

lint:
	cd backend && go vet ./...
