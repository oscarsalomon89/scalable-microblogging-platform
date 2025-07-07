SHELL := /bin/bash

VERSION := 1.0

DB_URL=postgres://admin:admin@localhost:5432/ponti_api_db?sslmode=disable
MIGRATIONS_DIR := ./migrations
MIGRATE := migrate -path $(MIGRATIONS_DIR) -database $(DB_URL)
MIGRATIONS_NAME=$(name)

.PHONY: migrate-up migrate-down migrate-force migrate-version

# Migrations
migrate-up: 
	@echo "Running migrations..."
	@$(MIGRATE) up

migrate-down: 
	@echo "Running migrations down..."
	@$(MIGRATE) down

migrate-create:
	@echo "Creating migration..."
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(MIGRATIONS_NAME)

migrate-force:
	@echo "Forcing migration..."
	@$(MIGRATE) force -1

migrate-version:
	@echo "Getting migration version..."
	@$(MIGRATE) version

# Core commands
run:
	@echo "Starting API..."
	@eval $$(egrep -v '^#' .env | xargs) APP_PATH=$$PWD go run ./cmd/api/*.go

test:
	@echo "Running tests..."
	@go test ./...

compose-up:
	@echo "Starting services in dev mode..."
	docker compose up --build -d

compose-down:
	@echo "Stopping services in dev mode..."
	docker compose down --remove-orphans

compose-dev-up:
	@echo "Starting services in dev mode..."
	docker compose -f docker-compose.dev.yml up --build -d

compose-dev-down:
	@echo "Stopping services in dev mode..."
	docker compose -f docker-compose.dev.yml down --remove-orphans

compose-logs:
	@echo "Fetching logs for dev services..."
	docker compose logs -f
