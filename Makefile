BIN_LOCATION ?= release/vela-k6
BIN_NAME ?= github.com/go-vela/vela-k6/cmd/vela-k6
MAIN_LOCATION ?= ./cmd/vela-k6

# capture the current date we build the application from
BUILD_DATE = $(shell date +%Y-%m-%dT%H:%M:%SZ)

# check if a git commit sha is already set
ifndef GITHUB_SHA
	# capture the current git commit sha we build the application from
	GITHUB_SHA = $(shell git rev-parse HEAD)
endif

# check if a git tag is already set
ifndef GITHUB_TAG
	# capture the current git tag we build the application from
	GITHUB_TAG = $(shell git describe --tag --abbrev=0)
endif

# check if a go version is already set
ifndef GOLANG_VERSION
	# capture the current go version we build the application from
	GOLANG_VERSION = $(shell go version | awk '{ print $$3 }')
endif

ifeq ($(shell uname -s), Darwin)
	GO_ENVS = GOOS=darwin
else
	GO_ENVS = GOOS=linux GOARCH=amd64
endif

# create a list of linker flags for building the golang application
LD_FLAGS = -X github.com/go-vela/vela-k6/version.Commit=${GITHUB_SHA} -X github.com/go-vela/vela-k6/version.Date=${BUILD_DATE} -X github.com/go-vela/vela-k6/version.Go=${GOLANG_VERSION} -X github.com/go-vela/vela-k6/version.Tag=${GITHUB_TAG}

.PHONY: deps
deps: go-tidy golangci-lint ## Install golang dependencies for the application

.PHONY: check
check: go-tidy check-all lint-all  ## Run all lint checks

.PHONY: clean
clean: clean-all go-tidy ## Clean up the application and test output

.PHONY: build
build: build-all  ## Compile the application

.PHONY: docker
docker: build-linux docker-build  ## Build a Docker image for local testing

.PHONY: test
test: test-all  ## Run all unit tests

.PHONY: coverage
coverage: test-all coverage-all ## Run tests and coverage visualization (only for local dev)

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

.PHONY: lint-all
lint-all: golangci-lint mdl yamllint hadolint shellcheck codespell

.PHONY: clean-all
clean-all:
	@rm -f perf-test.json
	@rm -f ./$(BIN_NAME)
	@rm -f coverage.*
	@rm -f unit-tests.xml
	
.PHONY: test-all
test-all:	
	@go test ./... -coverprofile=coverage.out

.PHONY: coverage-all
coverage-all:	
	@go tool cover -html=coverage.out

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

.PHONY: codespell
codespell:
ifeq ($(strip $(shell which codespell)),)
ifneq ($(strip $(shell which pip)),)
	@pip install codespell
endif
endif
	@codespell --skip node_modules,.git,public,package-lock.json,vendor
	@echo finished running codespell

.PHONY: yamllint
yamllint:
ifeq ($(strip $(shell which yamllint)),)
ifneq ($(strip $(shell which pip)),)
	@pip install yamllint
endif
endif
	@yamllint -f colored -c .yamllint.yml .
	@echo finished running yamllint

.PHONY: hadolint
hadolint:
ifeq ($(strip $(shell which hadolint)),)
ifeq ($(shell uname -s), Darwin)
	@brew install hadolint
endif
endif
	@hadolint --ignore DL3018 --ignore DL3007 ./Dockerfile
	@echo finished running hadolint

.PHONY: shellcheck
shellcheck:
ifeq ($(strip $(shell which shellcheck)),)
ifeq ($(shell uname -s), Darwin)
	@brew install shellcheck
endif
endif
	@find . -name \*.sh ! -path "./vendor/*" | xargs -I {} shellcheck {}
	@echo finished running shellcheck

.PHONY: mdl
mdl:
ifeq ($(strip $(shell which mdl)),)
ifneq ($(strip $(shell which gem)),)
	@gem install mdl
endif
endif
	@find . -name '*.md' ! -name "DOCS.md" ! -path "./vendor/*" | xargs mdl --rules ~MD007,~MD013
# ignore top level header for docs
	@mdl DOCS.md --rules ~MD007,~MD013,~MD002
	@echo finished running mdl

# The `build-all` target is intended to compile
# the Go source code into a binary.
#
# Usage: `make build`
.PHONY: build-all
build-all:
	@echo
	@echo "### Building release/vela-k6 binary"
	${GO_ENVS} CGO_ENABLED=0 go build -a -ldflags '${LD_FLAGS}' -o $(BIN_LOCATION) $(BIN_NAME)

# The `build-static-ci` target is intended to compile
# the Go source code into a statically linked binary
# when used within a CI environment.
#
# Usage: `make build-static-ci`
.PHONY: build-static-ci
build-static-ci:
	@echo
	@echo "### Building CI static release/vela-k6 binary"
	go build -a -ldflags '-s -w -extldflags "-static" ${LD_FLAGS}' -o $(BIN_LOCATION) $(BIN_NAME)

# The `build-linux` target is intended to compile
# the Go source code into a linux-compatible binary.
#
# Usage: `make build-linux`
.PHONY: build-linux
build-linux:
	@echo
	@echo "### Building release/vela-k6 binary for linux"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '${LD_FLAGS}' -o $(BIN_LOCATION) $(BIN_NAME)

# The `docker-build` target is intended to build
# the Docker image for the plugin.
#
# Usage: `make docker-build`
.PHONY: docker-build
docker-build: build-linux
	@echo
	@echo "### Building vela-k6:local image"
	@docker build --no-cache -t vela-k6:local .
