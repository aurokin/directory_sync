SHELL := /usr/bin/env bash

BIN_DIR := $(CURDIR)/bin
GOLANGCI_LINT_VERSION := v1.64.8

.PHONY: help
help:
	@printf "%s\n" "Targets:" \
		"  make tools       Install dev tools into ./bin" \
		"  make fmt         Format Go code (gofmt)" \
		"  make fmt-check   Fail if code is not gofmt'd" \
		"  make lint        Run golangci-lint" \
		"  make test        Run unit tests" \
		"  make check       Run fmt-check + lint + test" \
		"  make build       Build ./bin/dsync"

.PHONY: tools
tools: $(BIN_DIR)/golangci-lint

$(BIN_DIR)/golangci-lint:
	@mkdir -p "$(BIN_DIR)"
	GOBIN="$(BIN_DIR)" go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: fmt-check
fmt-check:
	@set -euo pipefail; \
	dirs="$$(go list -f '{{.Dir}}' ./... | sort -u)"; \
	if [ -z "$$dirs" ]; then exit 0; fi; \
	unformatted="$$(gofmt -l $$dirs)"; \
	if [ -n "$$unformatted" ]; then \
		echo "These files are not gofmt'd:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

.PHONY: lint
lint: tools
	"$(BIN_DIR)/golangci-lint" run

.PHONY: test
test:
	@go test ./...

.PHONY: check
check: fmt-check lint test

.PHONY: build
build:
	@mkdir -p "$(BIN_DIR)"
	@go build -o "$(BIN_DIR)/dsync" ./cmd/dsync
