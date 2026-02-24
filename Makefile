.PHONY: all build run test clean fmt vet help

# Default target
all: build

## Build the application
build:
	@echo "Building ROI calculator..."
	@mkdir -p bin
	@go build -o bin/roi ./cmd/roi

## Run the application
run:
	@go run ./cmd/roi

## Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

## Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

## Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf bin/

## Display this help message
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-10s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
