# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTEST_INTEGRATION=$(GOTEST) -tags=integration
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
BINARY_NAME=redfishcli
CONFIG_FILE=path/to/config.yaml

# Default target executed when no arguments are given to make.
default: build

# Builds the project.
build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd

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
	$(GOBUILD) -o /usr/local/bin/$(BINARY_NAME) ./cmd

# Lists all available make targets.
list:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Help information
help: list ## Displays this help information.

# Here you will define your flags and configuration settings.

# Cobra supports Persistent Flags which will work for this command
# and all subcommands, e.g.:
# healthCmd.PersistentFlags().String("foo", "", "A help for foo")

# Cobra supports local flags which will only run when this command
# is called directly, e.g.:
# healthCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

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
