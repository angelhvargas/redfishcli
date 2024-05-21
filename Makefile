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

# Default target executed when no arguments are given to make.
default: build

# Builds the project.
build:
	$(GOBUILD) -o $(BINARY_NAME) .

# Cleans the project.
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Runs the unit tests.
test:
	$(GOTEST) -v ./...

# Runs the integration tests.
integration-test:
	CONFIG_FILE=$(CONFIG_FILE) $(GOTEST_INTEGRATION) -v ./...

# Formats the code.
fmt:
	$(GOFMT) ./...

# Runs the linter.
vet:
	$(GOVET) ./...

# Runs the full test suite (unit and integration tests).
test-all: fmt vet test integration-test

# Installs the application.
install:
	$(GOBUILD) -o /usr/local/bin/$(BINARY_NAME) .

# Lists all available make targets.
list:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Help information
help: list ## Displays this help information.


# CI/CD tooling

# Targets with descriptions for the list command.
build: ## Builds the project.
clean: ## Cleans the project.
test: ## Runs the unit tests.
integration-test: ## Runs the integration tests.
fmt: ## Formats the code.
vet: ## Runs the linter.
test-all: ## Runs the full test suite (unit and integration tests).
install: ## Installs the application.
list: ## Lists all available make targets.
help: ## Displays this help information.
