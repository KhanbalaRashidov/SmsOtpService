# SMS OTP Service Makefile

.PHONY: help install build run test clean docker swagger dev

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	go mod download
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@latest

swagger: ## Generate swagger documentation
	swag init -g cmd/api/main.go --output docs

build: swagger ## Build the application
	go build -o bin/sms-otp-service cmd/api/main.go

run: swagger ## Run the application locally
	go run cmd/api/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf docs/
	rm -f coverage.out coverage.html
	go clean

format: ## Format code
	go fmt ./...
	goimports -w .

lint: ## Run linter
	golangci-lint run

docker-build: ## Build Docker image
	docker build -t sms-otp-service:latest .

docker-run: ## Run with Docker Compose
	docker-compose up -d

docker-stop: ## Stop Docker containers
	docker-compose down

docker-logs: ## Show Docker logs
	docker-compose logs -f sms-otp-service

docker-clean: ## Clean Docker containers and images
	docker-compose down -v
	docker rmi sms-otp-service:latest

dev: ## Start development environment
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5
	make swagger
	go run cmd/api/main.go

dev-with-pgadmin: ## Start development environment with PgAdmin
	docker-compose --profile dev up -d

migrate-up: ## Run database migrations (if you add migrate tool)
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/sms_otp_db?sslmode=disable" up

migrate-down: ## Rollback database migrations
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/sms_otp_db?sslmode=disable" down

api-test: ## Test API endpoints
	@echo "Testing API endpoints..."
	@echo "1. Health check:"
	curl -s http://localhost:8080/health | jq .
	@echo "\n2. Send OTP:"
	curl -s -X POST http://localhost:8080/api/v1/otp/send \
		-H "Content-Type: application/json" \
		-d '{"phone_number": "+994501234567", "purpose": "verification"}' | jq .

docs: swagger ## Open Swagger documentation
	@echo "Swagger docs available at: http://localhost:8080/docs/"
	@command -v open >/dev/null 2>&1 && open http://localhost:8080/docs/ || \
	command -v xdg-open >/dev/null 2>&1 && xdg-open http://localhost:8080/docs/ || \
	echo "Please open http://localhost:8080/docs/ in your browser"

setup: install swagger ## Initial project setup
	@echo "Setting up SMS OTP Service..."
	@echo "1. Installing dependencies..."
	@make install
	@echo "2. Generating swagger docs..."
	@make swagger
	@echo "3. Creating .env file..."
	@cp .env.example .env 2>/dev/null || true
	@echo "Setup complete! Run 'make dev' to start development server."