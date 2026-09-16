.PHONY: build test test-race lint lint-docs tidy vet fmt run clean install release

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

lint-docs:
	@echo "🐷 Linting markdown..."
	@if command -v npx >/dev/null 2>&1; then \
		npx --yes markdownlint-cli2 "**/*.md" "#node_modules"; \
	else \
		echo "npx not installed — skipping markdown lint"; \
	fi

tidy:
	@echo "🐷 Checking go.mod is tidy..."
	@$(GO) mod tidy -diff

vet:
	@echo "🐷 Running go vet..."
	@$(GO) vet ./...

fmt:
	@echo "🐷 Formatting code..."
	@gofmt -s -w .
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w -local github.com/abraira85/pork .; \
	else \
		echo "goimports not installed — skipping import grouping"; \
	fi

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
