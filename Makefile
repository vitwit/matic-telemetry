BINARY_NAME := telemetry

.PHONY: all build clean run test

all: build

build:
	@echo "==> Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

run: build
	@echo "==> Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

test:
	@echo "==> Running tests..."
	go test ./...

clean:
	@echo "==> Cleaning up..."
	rm -f $(BINARY_NAME)
