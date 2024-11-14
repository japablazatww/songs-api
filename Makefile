## Makefile

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!".

## stop_searcher: stops the searcher container
stop_searcher:
	@echo "Stopping searcher container..."
	docker-compose stop searcher

## start_searcher: starts the searcher container
start_searcher:
	@echo "Starting searcher container..."
	docker-compose start searcher

## restart_postgres: restarts the postgres container
restart_postgres:
	@echo "Restarting PostgreSQL container..."
	docker-compose restart postgresdb

## create_db: creates the database
create_db:
	@echo "Creating database 'songs'..."
	docker-compose exec -T postgres createdb -U postgres songs

## drop_db: drops the database, make sure to stop the searcher container first
drop_db:
	@echo "Dropping database 'songs'..."
	docker-compose exec -T postgres dropdb -U postgres songs

## migrate_db: migrates the database
migrate_db_up:
	@echo "Migrating database..."
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/songs?sslmode=disable" -verbose up

## migrate_db: migrates the database to a previous version
migrate_db_down:
	@echo "Migrating database..."
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/songs?sslmode=disable" -verbose down

# copy_env: copies the .env file to the searcher container
copy_env:
	@echo "Copying .env.example and renaming it to .env..."
	cp .env.example .env

.PHONY: up up_build down