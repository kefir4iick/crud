APP_NAME = crud
BIN_DIR = bin
TEST_PKGS = ./...

.PHONY: test deps cover lint build run clean install-lint migrate-up run-db stop-db

deps:
	go mod tidy
	go mod download

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

test: deps
	go test -v $(TEST_PKGS) || true

cover: deps
	go test -coverprofile=coverage.out $(TEST_PKGS)
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint: deps
	golangci-lint run $(TEST_PKGS)

build: deps
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd

run: build
	$(BIN_DIR)/$(APP_NAME)

clean:
	rm -rf $(BIN_DIR) coverage.out coverage.html

migrate-up:
	psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -c "CREATE TABLE IF NOT EXISTS cars (id VARCHAR(36) PRIMARY KEY, make VARCHAR(255) NOT NULL, model VARCHAR(255) NOT NULL, year INTEGER NOT NULL, price INTEGER NOT NULL);"

run-db:
	docker run --name go-postgres -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_NAME) -p 5432:5432 -d postgres:13-alpine

stop-db:
	docker stop go-postgres
	docker rm go-postgres
