BINARY_NAME := driver-scanner
VERSION ?= dev
REGISTRY ?= docker.io
IMAGE_NAME ?= $(REGISTRY)/$(BINARY_NAME)
BIN_DIR := .bin
COVERAGE_DIR := .coverage
REPORT_DIR := .reports

.PHONY: help build vet test test-ctrf lint clean docker-build docker-push

## help: show this help message
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## //' | column -t -s ':'

## build: compile the binary with version info injected via ldflags
build:
	@mkdir -p $(BIN_DIR)
	go build -ldflags "-X main.version=$(VERSION)" -o $(BIN_DIR)/$(BINARY_NAME) ./cmd

## vet: run go vet static analysis
vet:
	go vet ./...

## test: run tests with coverage report (depends on build)
test: build
	@mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@echo ""
	@echo "=== Coverage ==="
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1

## test-ctrf: run tests with CTRF output for GitHub reporting
test-ctrf:
	@mkdir -p $(REPORT_DIR)
	@echo "🧪 Running tests with CTRF output..."
	@if ! command -v gotestsum >/dev/null 2>&1; then \
		echo "📦 Installing gotestsum..."; \
		go install gotest.tools/gotestsum@latest; \
	fi
	@if ! command -v go-ctrf-json-reporter >/dev/null 2>&1; then \
		echo "📦 Installing go-ctrf-json-reporter..."; \
		go install github.com/ctrf-io/go-ctrf-json-reporter/cmd/go-ctrf-json-reporter@latest; \
	fi
	@echo "🔄 Running tests with CTRF reporter..."
	gotestsum --jsonfile $(REPORT_DIR)/gotestsum.json && go-ctrf-json-reporter -o $(REPORT_DIR)/ctrf-report.json < $(REPORT_DIR)/gotestsum.json

## lint: run golangci-lint v2 via go run (no install required)
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0 run ./...

## clean: remove built binary and coverage artifacts
clean:
	rm -rf $(BIN_DIR) $(COVERAGE_DIR)

## docker-build: build docker image with multi-stage (depends on test)
docker-build: test
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE_NAME):$(VERSION) .

## docker-push: push docker image to registry
docker-push:
	docker push $(IMAGE_NAME):$(VERSION)
