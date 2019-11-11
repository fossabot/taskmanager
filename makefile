.PHONY: all fmt lint test bench help

all: help

fmt: ## gofmt all project
	@gofmt -l -s -w .

lint: ## Lint the files
	@golangci-lint run

test: ## run unit tests
	@go test -race -short -coverprofile=coverage.txt

bench: ## run tests and benchmarks
	@go test  -bench=. -benchmem

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'