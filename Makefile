.PHONY: build test lint vet fmt run clean install release

GO ?= $(shell command -v go 2> /dev/null || echo $$HOME/.local/go/bin/go)

build:
	@echo "🐷 Building pork..."
	@$(GO) build -o bin/pork main.go

test:
	@echo "🐷 Running tests..."
	@$(GO) test ./...

test-race:
	@echo "🐷 Running tests with race detector..."
	@$(GO) test ./... -race

lint:
	@echo "🐷 Running linter..."
	@golangci-lint run

vet:
	@echo "🐷 Running go vet..."
	@$(GO) vet ./...

fmt:
	@echo "🐷 Formatting code..."
	@gofmt -s -w .

run: build
	@./bin/pork

clean:
	@echo "🐷 Cleaning up..."
	@rm -rf bin/ dist/

install:
	@echo "🐷 Installing pork..."
	@$(GO) install

release:
	@echo "🐷 Tagging a release with GoReleaser..."
	@goreleaser release --clean