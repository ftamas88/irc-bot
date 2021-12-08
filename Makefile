# This file is the development Makefile for the project in this repository.
# All variables listed here are used as substitution in these Makefile targets.
SERVICE_NAME = bot

################################################################################

.PHONY: setup
setup: ## Downloads and install various libs for development.
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	go get golang.org/x/tools/cmd/goimports

.PHONY: build
build: lint ## Builds project binary.
	go build -ldflags -race -o ./bin/$(SERVICE_NAME).exe -v ./cmd/app/app.go

.PHONY: test
test: lint ## Runs the test suite. Some projects might rely on a local development infrasructure to run tests. See `infra-up`.
	go test -v -race -bench=./... -benchmem -timeout=120s -cover -coverprofile=./test/coverage.txt ./...

.PHONY: test-quick ## Run the quick test suite
test-quick:
	go test -short -failfast

.PHONY: run
run: build .env ## Builds project binary and executes it.
	bin/$(SERVICE_NAME)

.PHONY: full
full: clean build fmt lint test ## Cleans up, builds the service, reformats code, lints and runs the test suite.

################################################################################

.PHONY: lint
lint: ## Runs linter against the service codebase.
	golangci-lint run
	@echo "[✔] Linter OK"


.PHONY: fmt
fmt: ## Runs gofmt against the service codebase.
	gofmt -w -s .
	goimports -w .
	go clean ./...

.PHONY: tidy
tidy: ## Runs go mod tidy against the service codebase.
	go mod tidy

.PHONY: clean
clean: ## Removes temporary files and deletes the service binary.
	go clean ./...
	rm -f bin/$(SERVICE_NAME)

.env:
	cp .env.dist .env

.PHONY: version
version: ## Displays the current version of the Go toolchain.
	go version

################################################################################
