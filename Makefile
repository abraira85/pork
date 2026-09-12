.PHONY: build test lint run clean install

GO ?= $(shell command -v go 2> /dev/null || echo $$HOME/.local/go/bin/go)

build:
	@echo "🐷 Building pork..."
	@$(GO) build -o bin/pork main.go

test:
	@echo "🐷 Running tests..."
	@$(GO) test ./...

lint:
	@echo "🐷 Running linter..."
	@golangci-lint run

run: build
	@./bin/pork

clean:
	@echo "🐷 Cleaning up..."
	@rm -rf bin/

install:
	@echo "🐷 Installing pork..."
	@go install
