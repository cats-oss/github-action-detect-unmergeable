GO_CMD := go
GO_BUILD := $(GO_CMD) build

DIST_APP_BIN_NAME := app

all: help

help:
	@echo "Specify the task"
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@exit 1


.PHONY: clean
clean: ## Clean up generated artifacts
	rm -rf $(DIST_APP_BIN_NAME)


.PHONY: build
build: clean ## Build the application
	GO111MODULE=on $(GO_BUILD) -o $(DIST_APP_BIN_NAME)