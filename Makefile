# Copyright (c) 2023 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

BIN_NAME ?= vela-k6
BIN_LOCATION ?= ./
MAIN_LOCATION ?= .

.PHONY: deps
deps: go-tidy golangci-lint ## Install golang dependencies for the application

.PHONY: check
check: go-tidy check-all golangci-lint  ## Run all lint checks

.PHONY: clean
clean: clean-all go-tidy ## Clean up the application and test output

.PHONY: build
build: build-all  ## Compile the application

.PHONY: build-docker
build-docker:  ## Compile the application for testing locally with Docker
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_NAME)

.PHONY: test
test: test-all  ## Run all unit tests

.PHONY: help
help: ## Show all valid options
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: start
start: ## Run application
	@go run $(MAIN_LOCATION)

.PHONY: check-all
check-all:	
	@go vet ./...
	@go fmt ./...

.PHONY: clean-all
clean-all:
	@rm -f perf-test.json
	@rm -f ./$(BIN_NAME)
	@rm -f coverage.*
	@rm -f unit-tests.xml
	
.PHONY: test-all
test-all:	
	@go test ./... -coverprofile=coverage.out

.PHONY: go-tidy
go-tidy:	
	@go mod tidy

.PHONY: golangci-lint
golangci-lint:	
ifeq ($(strip $(shell which golangci-lint)),)
ifeq ($(shell uname -s), Darwin)
	@brew install golangci-lint
endif
endif
	@golangci-lint run ./...
	@echo finished running golangci-lint

.PHONY: build-all
build-all:	
ifeq ($(shell uname -s), Darwin)
	GOOS=darwin
endif
	CGO_ENABLED=0
	GOOS=linux
	GOARCH=amd64
	@go build -a -ldflags '-extldflags "-static"' -o $(BIN_NAME) $(BIN_LOCATION)
	@echo finished building golangci-lint
