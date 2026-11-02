.PHONY: help build run test clean migrate

help:
	@echo "Available commands:"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make migrate   - Run database migrations"
	@echo "  make deps      - Download dependencies"

build:
	@echo "Building application..."
	@go build -o bin/server ./cmd/server

run:
	@echo "Running application..."
	@go run ./cmd/server/main.go

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@go clean

migrate:
	@echo "Running database migrations..."
	@cd migrations && go run migrate.go

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

install-deps:
	@echo "Installing required dependencies..."
	@go get github.com/go-sql-driver/mysql
	@go get golang.org/x/crypto/bcrypt
	@go mod tidy
