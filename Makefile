.PHONY: help test test-unit test-integration test-s3 test-db test-api lint build clean setup-local teardown-local

# Default target
help:
	@echo "Available targets:"
	@echo "  make test           - Run all unit tests"
	@echo "  make test-unit      - Run unit tests only"
	@echo "  make test-integration - Run all integration tests (requires setup)"
	@echo "  make test-s3        - Run S3 integration tests"
	@echo "  make test-db        - Run database integration tests"
	@echo "  make test-api       - Run API integration tests"
	@echo "  make lint           - Run linters"
	@echo "  make build          - Build the application"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make setup-local    - Setup local testing environment"
	@echo "  make teardown-local - Teardown local testing environment"

# Run all unit tests
test: test-unit

test-unit:
	go test -v -race ./tests

# Run all integration tests
test-integration: test-s3 test-db test-api

test-s3:
	go test -v -tags=integration ./tests -s3

test-db:
	go test -v -tags=integration ./tests -database

test-api:
	go test -v -tags=integration ./tests -api

# Run linters
lint:
	@echo "Running gofmt..."
	@gofmt -l -w .
	@echo "Running go vet..."
	@go vet ./...
	@echo "Running goimports..."
	@which goimports > /dev/null || go install golang.org/x/tools/cmd/goimports@latest
	@goimports -local github.com/digitaldrywood/github-integration-testing-demo -w .

# Build the application
build:
	go build -o app ./src

# Clean build artifacts
clean:
	rm -f app
	rm -rf dist/
	rm -rf test-results/
	rm -f *.log
	rm -f coverage.txt

# Setup local testing environment
setup-local:
	./scripts/setup-local-testing.sh

# Teardown local testing environment
teardown-local:
	@echo "Stopping containers..."
	@docker stop minio-test postgres-test 2>/dev/null || true
	@docker rm minio-test postgres-test 2>/dev/null || true
	@echo "Cleanup complete"

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run tests with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"