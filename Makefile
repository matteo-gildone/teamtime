.PHONY: build test clean install lint

# Build the binary
build:
	go build -o bin/teamtime main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install the binary
install:
	go install main.go

# Lint the code
lint:
	go vet ./...
	go fmt ./...

# Run all checks
check: lint test