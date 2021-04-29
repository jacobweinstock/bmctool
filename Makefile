BINARY:=bmctool
OSFLAG:= $(shell go env GOHOSTOS)
GIT_COMMIT:=$(shell git rev-parse --short HEAD)
BUILD_ARGS:=GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags '-s -w -extldflags "-static" -X ${REPO}/server/httpsvr.GitRev=${GIT_COMMIT}'

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run unit tests
	go test -gcflags=-l -v -covermode=count ./...

.PHONY: test-ci
test-ci: ## Run unit tests for CI
	go test -coverprofile=cover.out ./... > test.output

.PHONY: cover
cover: ## Run unit tests with coverage report
	go test -gcflags=-l -coverprofile=cover.out ./...
	go tool cover -func=cover.out
	rm -rf cover.out

.PHONY: lint
lint:  ## Run linting
	@echo be sure golangci-lint is installed: https://golangci-lint.run/usage/install/
	golangci-lint run

.PHONY: goimports
goimports: ## run goimports updating files in place
	@echo be sure goimports is installed
	goimports -w .

.PHONY: goimports-check
goimports-check: ## run goimports displaying diffs
	@echo be sure goimports is installed
	goimports -d . | (! grep .)

.PHONY: all-checks
all-checks: cover lint goimports ## run all tests and formatters
	go vet ./...

.PHONY: darwin
darwin: ## compile for darwin
	GOOS=darwin ${BUILD_ARGS} -o bin/${BINARY}-darwin-amd64 main.go

.PHONY: linux
linux: ## compile for linux
	GOOS=linux ${BUILD_ARGS} -o bin/${BINARY}-linux-amd64 main.go

.PHONY: build
build: ## compile the binary for the native OS
ifeq (${OSFLAG},linux)
	@$(MAKE) linux
else
	@$(MAKE) darwin
endif

.PHONY: image
image: ## make the Container Image
	docker build -t ${BINARY}:local . 
