include ./config/config-local.env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DOCKER
# ==================================================================================== #

## up: Build and start the service container along with dependencies (Postgres, Kafka)
.PHONY: up
up: postgres/ready kafka/ready $(SERVICE_NAME)/start
	@echo "$(SERVICE_NAME) is up and running."

## build: Build Docker images for the service
.PHONY: build
build:
	@docker-compose build $(SERVICE_NAME)
	@echo "Docker image for $(SERVICE_NAME) built."

## postgres/start: Start the PostgreSQL container if not running
.PHONY: postgres/start
postgres/start:
	@docker-compose up -d postgres
	@echo "Starting PostgreSQL container..."

## postgres/ready: Wait until PostgreSQL is ready
.PHONY: postgres/ready
postgres/ready: postgres/start
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker-compose exec postgres pg_isready -U $$POSTGRES_USER -d $$POSTGRESQL_DB; do \
		echo "Waiting for PostgreSQL database $(POSTGRESQL_DB)..."; \
		sleep 2; \
	done
	@echo "PostgreSQL database $(POSTGRESQL_DB) is ready!"

# ==================================================================================== #
# KAFKA TASKS
# ==================================================================================== #

## kafka/start: Start the Kafka container if not running
.PHONY: kafka/start
kafka/start:
	@docker-compose up -d kafka
	@echo "Starting Kafka container..."

## kafka/ready: Wait until Kafka is ready
.PHONY: kafka/ready
kafka/ready: kafka/start
	@echo "Waiting for Kafka to be ready..."
	@until docker-compose ps kafka | grep -q "Up"; do \
		echo "Waiting for Kafka to be up..."; \
		sleep 2; \
	done
	@echo "Kafka is ready!"

# ==================================================================================== #
# MIGRATION TASKS
# ==================================================================================== #

## migrate/up: Apply all migrations to the database
.PHONY: migrate/up
migrate/up:
	@echo "Applying migrations for $(SERVICE_NAME)..."
	@migrate -path ./migrations -database postgres://$(POSTGRESQL_USER):$(POSTGRESQL_PASSWORD)@localhost:$(POSTGRESQL_PORT)/$(POSTGRESQL_DB)?sslmode=$(POSTGRESQL_SSLMODE) up
	@echo "Migrations applied for $(SERVICE_NAME)."

## migrate/down: Rollback the last migration
.PHONY: migrate/down
migrate/down:
	@echo "Rolling back the last migration for $(SERVICE_NAME)..."
	@migrate -path ./migrations -database postgres://$(POSTGRESQL_USER):$(POSTGRESQL_PASSWORD)@localhost:$(POSTGRESQL_PORT)/$(POSTGRESQL_DB)?sslmode=$(POSTGRESQL_SSLMODE) down
	@echo "Last migration rolled back for $(SERVICE_NAME)."

## migrate/new name=$(name): Create a new migration
.PHONY: migrate/new
migrate/new:
	@echo "Creating new migration file for $(name)..."
	@migrate create -seq -ext sql -dir ./migrations $(name)
	@echo "New migration file created in ./migrations/"
	
# ==================================================================================== #
# SERVICE TASKS
# ==================================================================================== #

## $(SERVICE_NAME)/start: Start the service container
.PHONY: $(SERVICE_NAME)/start
$(SERVICE_NAME)/start:
	@docker-compose up -d $(SERVICE_NAME)
	@echo "$(SERVICE_NAME) is started."

## run: Fully start the service (Build, wait for DB + Kafka, then start the service)
.PHONY: run
run: build up
	@echo "$(SERVICE_NAME) is fully up and running."

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## gen: generates tests mocks
.PHONY: gen 
gen: 
	mockgen -source=internal/auth/pg_repository.go -destination=internal/auth/mock/pg_repository_mock.go
	mockgen -source=internal/auth/usecase.go -destination=internal/auth/mock/usecase_mock.go

## swag: generates swagger documentation
.PHONY: swag
swag:
	swag init -g cmd/api/main.go

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify

	@echo 'Formatting code...'
	go fmt ./...
	swag fmt ./...

	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...

	@echo 'Running tests...'
	go test -race -vet=off ./...