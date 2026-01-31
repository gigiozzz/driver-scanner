BINARY_NAME := driver-scanner
VERSION ?= dev
REGISTRY ?= docker.io
IMAGE_NAME ?= $(REGISTRY)/$(BINARY_NAME)
BIN_DIR := .bin
COVERAGE_DIR := .coverage

.PHONY: help build test lint clean docker-build docker-push

## help: show this help message
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## //' | column -t -s ':'

## build: compile the binary with version info injected via ldflags
build:
	@mkdir -p $(BIN_DIR)
	go build -ldflags "-X main.version=$(VERSION)" -o $(BIN_DIR)/$(BINARY_NAME) ./cmd

## test: run tests with coverage report (depends on build)
test: build
	@mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@echo ""
	@echo "=== Coverage ==="
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1

## lint: run golangci-lint
lint:
	golangci-lint run ./...

## clean: remove built binary and coverage artifacts
clean:
	rm -rf $(BIN_DIR) $(COVERAGE_DIR)

## docker-build: build docker image with multi-stage (depends on test)
docker-build: test
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE_NAME):$(VERSION) .

## docker-push: push docker image to registry
docker-push:
	docker push $(IMAGE_NAME):$(VERSION)
