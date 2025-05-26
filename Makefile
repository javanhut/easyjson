# EasyJSON Makefile

.PHONY: all test test-verbose test-coverage benchmark clean build example install lint fmt vet

# Default target
all: fmt vet test

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Run the example
example:
	@echo "Running example..."
	cd example && go run main.go

# Build the example
build:
	@echo "Building example..."
	cd example && go build -o easyjson-example main.go

# Install dependencies and tools
install:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -f coverage.out coverage.html
	rm -f example/easyjson-example
	go clean ./...

# Initialize module (run once when setting up)
init:
	@echo "Initializing Go module..."
	go mod init github.com/yourname/easyjson
	go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Format, vet, and test (default)"
	@echo "  test         - Run tests"
	@echo "  test-verbose - Run tests with verbose output and race detection"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  benchmark    - Run benchmarks"
	@echo "  example      - Run the example program"
	@echo "  build        - Build the example program"
	@echo "  install      - Install dependencies"
	@echo "  lint         - Run linter (requires golangci-lint)"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  clean        - Clean build artifacts"
	@echo "  init         - Initialize Go module"
	@echo "  help         - Show this help message"

# Development workflow
dev: fmt vet test example

# CI workflow
ci: fmt vet test-coverage benchmark
