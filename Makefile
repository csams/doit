.DEFAULT_GOAL := build

.PHONY: default
default: build ;

.PHONY: build
build: require-go ## build the thing
	go build -o ./bin/doit main.go

.PHONY: clean
clean: ## clean out the binaries
	rm -rf ./bin/*

.PHONY: require-%
require-%:
	@if ! command -v $* 1> /dev/null 2>&1; then echo "$* not found in \$$PATH"; exit 1; fi

.PHONY: help
help: ## Show this help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
