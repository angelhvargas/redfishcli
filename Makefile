# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTEST_INTEGRATION=$(GOTEST) -tags=integration
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
BINARY_NAME=redfishcli
CONFIG_FILE=""
DOCKER_IMAGE=redfishcli:latest

# Default target executed when no arguments are given to make.
default: build

# Builds the project.
build: ## Builds the project.
	$(GOBUILD) -o $(BINARY_NAME) .

# Cleans the project.
clean: ## Cleans the project.
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Runs the unit tests.
test: ## Runs the unit tests.
	$(GOTEST) -v ./...

# Runs the integration tests.
integration-test: ## Runs the integration tests.
	CONFIG_FILE=$(CONFIG_FILE) $(GOTEST_INTEGRATION) -v ./...

# Formats the code.
fmt: ## Formats the code.
	$(GOFMT) ./...

# Runs the linter.
vet: ## Runs the linter.
	$(GOVET) ./...

# Runs the full test suite (unit and integration tests).
test-all: ## Runs the full test suite (unit and integration tests).
	$(MAKE) fmt
	$(MAKE) vet
	$(MAKE) test
	$(MAKE) integration-test

# Installs the application.
install: ## Installs the application.
	$(GOBUILD) -o /usr/local/bin/$(BINARY_NAME) .

# Docker build
docker-build: ## Builds the Docker image.
	docker build -t $(DOCKER_IMAGE) .

# Docker run example (runs help)
docker-run: ## Runs the Docker container (help command).
	docker run --rm $(DOCKER_IMAGE) --help

# Lists all available make targets.
help: ## Displays this help information.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
