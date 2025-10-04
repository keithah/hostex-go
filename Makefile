.PHONY: test test-integration test-unit test-all fmt vet lint clean build example

# Run unit tests only (no API key needed)
test-unit:
	go test -v -race -coverprofile=coverage.out -short ./...

# Run integration tests (requires HOSTEX_API_KEY)
test-integration:
	@if [ -z "$$HOSTEX_API_KEY" ]; then \
		echo "Error: HOSTEX_API_KEY environment variable is not set"; \
		echo "Either:"; \
		echo "  1. Set it: export HOSTEX_API_KEY=your_key"; \
		echo "  2. Or source .env: source .env && make test-integration"; \
		exit 1; \
	fi
	go test -v -run TestIntegration ./...

# Run all tests (unit + integration)
test-all:
	@if [ -z "$$HOSTEX_API_KEY" ]; then \
		echo "Error: HOSTEX_API_KEY environment variable is not set"; \
		exit 1; \
	fi
	go test -v -race -coverprofile=coverage.out ./...

# Default test (unit only, safe for CI without secrets)
test: test-unit

# Format code
fmt:
	gofmt -s -w .

# Check formatting
fmt-check:
	@if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "Code is not formatted:"; \
		gofmt -s -d .; \
		exit 1; \
	fi

# Run go vet
vet:
	go vet ./...

# Run linters
lint: fmt-check vet
	@command -v staticcheck > /dev/null || go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

# Clean build artifacts
clean:
	rm -f coverage.out
	rm -f example/example

# Build the library
build:
	go build ./...

# Build and run example
example:
	@if [ -z "$$HOSTEX_API_KEY" ]; then \
		echo "Error: HOSTEX_API_KEY not set. Run: source .env && make example"; \
		exit 1; \
	fi
	cd example && go run main.go

# Show coverage in browser
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run all checks (format, vet, test)
check: lint test-unit

# Show help
help:
	@echo "Available targets:"
	@echo "  test            - Run unit tests (default, no API key needed)"
	@echo "  test-unit       - Run unit tests only"
	@echo "  test-integration - Run integration tests (requires HOSTEX_API_KEY)"
	@echo "  test-all        - Run all tests with coverage"
	@echo "  fmt             - Format code"
	@echo "  fmt-check       - Check code formatting"
	@echo "  vet             - Run go vet"
	@echo "  lint            - Run all linters"
	@echo "  build           - Build the library"
	@echo "  example         - Build and run example program"
	@echo "  coverage        - Generate and view coverage report"
	@echo "  check           - Run all checks (lint + test-unit)"
	@echo "  clean           - Clean build artifacts"
	@echo "  help            - Show this help message"
