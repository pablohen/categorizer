.PHONY: run build test clean

# Binary name
BINARY_NAME=categorizer

# Build the application
build:
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME) cmd/categorizer/main.go
	@echo "Build complete: bin/$(BINARY_NAME)"

# Run the application
run:
	@go run cmd/categorizer/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin
	@echo "Clean complete"
