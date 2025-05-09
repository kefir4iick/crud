APP_NAME = go
BIN_DIR = bin
TEST_PKGS = ./...

.PHONY: test deps cover lint build run clean install-lint migrate-up migrateup1 migratedown migratedown1 run-db stop-db

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

migrateup:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -verbose down 1

run-db:
	docker run --name go-postgres -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_NAME) -p $(DB_PORT):5432 -d postgres:13-alpine

stop-db:
	docker stop go-postgres
	docker rm go-postgres
