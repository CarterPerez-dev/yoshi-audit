# ©AngelaMos | 2026
# Justfile
# =============================================================================
# yoshi-audit — System monitoring, Docker prune manager, and deep audit TUI
# =============================================================================

set export
set shell := ["bash", "-uc"]

project := file_name(justfile_directory())
version := `git describe --tags --always 2>/dev/null || echo "dev"`

# =============================================================================
# Default
# =============================================================================

default:
    @just --list --unsorted

# =============================================================================
# Linting and Formatting
# =============================================================================

[group('lint')]
lint *ARGS:
    golangci-lint run --timeout=5m {{ARGS}}

[group('lint')]
lint-fix:
    golangci-lint run --timeout=5m --fix

[group('lint')]
format:
    golangci-lint fmt

[group('lint')]
tidy:
    go mod tidy

[group('lint')]
vet:
    go vet ./...

# =============================================================================
# Testing
# =============================================================================

[group('test')]
test *ARGS:
    go test -race ./... {{ARGS}}

[group('test')]
test-v *ARGS:
    go test -race -v ./... {{ARGS}}

[group('test')]
cover:
    go test -race -cover ./...

[group('test')]
cover-html:
    go test -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report: coverage.html"

# =============================================================================
# CI / Quality
# =============================================================================

[group('ci')]
ci: lint test
    @echo "All checks passed."

[group('ci')]
check: lint vet

[group('ci')]
pre-push: tidy vet test
    @echo "Ready to push."

# =============================================================================
# Development
# =============================================================================

[group('dev')]
run:
    go run ./cmd

[group('dev')]
run-debug:
    YOSHI_DEBUG=1 go run ./cmd

# =============================================================================
# Build (Production)
# =============================================================================

[group('prod')]
build:
    go build -ldflags="-s -w" -o bin/yoshi-audit ./cmd
    @echo "Built: bin/yoshi-audit ($(du -h bin/yoshi-audit | cut -f1))"

[group('prod')]
build-debug:
    go build -o bin/yoshi-audit ./cmd

[group('prod')]
install:
    go install ./cmd

# =============================================================================
# Utilities
# =============================================================================

[group('util')]
info:
    @echo "Project:  {{project}}"
    @echo "Version:  {{version}}"
    @echo "Go:       $(go version | cut -d' ' -f3)"
    @echo "OS:       {{os()}} ({{arch()}})"
    @echo "Module:   $(head -1 go.mod | cut -d' ' -f2)"

[group('util')]
update:
    go get -u ./...
    go mod tidy

[group('util')]
clean:
    -rm -rf bin/ coverage.out coverage.html
    @echo "Cleaned build artifacts."
