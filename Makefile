.PHONY: db_create db_drop db_reset db_migrate db_seed build run

include .env

# Commands
db_create:
	@echo "Creating database..."
	mysql -u $(DB_USER) -p -e "CREATE DATABASE IF NOT EXISTS $(DB_NAME);"

db_drop:
	@echo "Dropping database..."
	mysql -u $(DB_USER) -p -e "DROP DATABASE IF EXISTS $(DB_NAME);"

db_reset: db_drop db_create

db_migrate:
	@echo "Running database migrations..."
	for file in rlnotif/mysqldb/migration/*.sql; do \
		echo "Running migration: $$file"; \
		mysql -u $(DB_USER) -p $(DB_NAME) < "$$file"; \
	done

db_seed:
	@echo "Seeding database..."
	@go build -o bin/seed ./rlnotif/cmd/seed/main.go
	@seed

build:
	@go build -o bin/rlnotif ./rlnotif/cmd/api-server/main.go

run: build
	@./bin/rlnotif
