include .env

dev:
	go run cmd/main.go
	
# Testing commands
test:
	@if [ -z "$(func)" ]; then \
		go test $(shell go list ./pkg/... | grep -v /mocks | grep -v "_test.go") --cover; \
	else \
		go test $(shell go list ./pkg/... | grep -v /mocks | grep -v "_test.go") -run $(func) -v; \
	fi
	
test-coverage:
	go test $(shell go list ./pkg/... | grep -v /mocks | grep -v "_test.go") -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	open coverage.html

# Mock generation commands
# Mockery using package configuration, check file .mockery.yaml for more details
# Once add new package, check .mockery.yaml and specify the package if needed
gen-mocks:
	@make remove-mocks
	mockery

remove-mocks:
	find ./ -name "mock_*.go" -type f -exec rm -f {} \;

# BOT DB
initdb:
	initdb --username=${DB_USER}
	@echo "Configuring PostgreSQL port..."
	@if [ -f ${PGDATA}/postgresql.conf ]; then \
		sed -i 's/^#port = .*/port = 5433/' ${PGDATA}/postgresql.conf; \
		echo "Port configured successfully."; \
	else \
		echo "Error: postgresql.conf not found in the specified location."; \
		exit 1; \
	fi

createdb:
	@echo "Creating database..."
	@createdb -p 5433 --no-password ${DB_NAME} -U ${DB_USER} || (echo "Failed to create database. Make sure PostgreSQL is running and the specified port is correct." && exit 1)
	@echo "Database '${DB_NAME}' created successfully."

startdb:
	pg_ctl -D ${PGDATA} start -o "-k ${PGHOST}"

stopdb:
	@pg_ctl -D ${PGDATA} stop -m fast || true

removedb:
	rm -rf ${PGDATA}

migrate-new:
	sql-migrate new -env=local ${name}

migrate-up:
	sql-migrate up -env=local

migrate-down:
	sql-migrate down -env=local

# Seed data from SQL files
seed-data:
	@echo "Seeding data from ./migrations/seeds/*.sql"
	@for file in ./migrations/seeds/*.sql; do \
		echo "Executing $$file..."; \
		psql -U ${DB_USER} -d ${DB_NAME} -p ${DB_PORT} -f $$file || exit 1; \
	done
	@echo "Data seeding completed successfully."
