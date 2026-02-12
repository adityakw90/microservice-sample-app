.PHONY: test-api test-api-cover test-api-clean api-build api-run help verbose

# Default target
.DEFAULT_GOAL := api-help

# Verbose target (used as: make test verbose)
verbose:
	@:

# API Test (usage: make test-api or make test-api verbose)
test-api:
	@if echo "$(MAKECMDGOALS)" | grep -q "verbose"; then \
		echo "Running API tests with verbose output..."; \
		cd api-gateway && go test -v ./...; \
	else \
		echo "Running API tests..."; \
		cd api-gateway && go test ./...; \
	fi

# API Test with coverage
test-api-cover:
	@echo "Running API tests with coverage..."
	@if echo "$(MAKECMDGOALS)" | grep -q "verbose"; then \
		cd api-gateway && go test -v -race -covermode=atomic -coverprofile=coverage.txt ./...; \
	else \
		cd api-gateway && go test -race -covermode=atomic -coverprofile=coverage.txt ./...; \
	fi

# Clean API test cache and coverage files
test-api-clean:
	@echo "Cleaning API test cache and coverage files..."
	@cd api-gateway && go clean -testcache
	@cd api-gateway && rm -f coverage.txt
	@echo "Clean complete"

build-api:
	@echo "Building the API application..."
	@cd api-gateway && GOOS=linux GOARCH=amd64 go build -o bin/app cmd/main.go

lint-api:
	@echo "Linting the API application..."
	@cd api-gateway && golangci-lint run

run-api:
	@echo "Running the API application..."
	@cd api-gateway && go run cmd/main.go

# Help target
help:
	@echo "Available targets:"
	@echo "  test [verbose]       - Run all tests (add verbose for verbose output)"
	@echo "  test-cover [verbose] - Run tests with coverage (add verbose for verbose output)"
	@echo "  test-clean           - Clean test cache and coverage files"
	@echo "  bench [verbose]      - Run all benchmarks (add verbose for verbose output)"
	@echo "  build                - Build the application"
	@echo "  run                  - Run the application"
	@echo "  help                 - Show this help message"
