.PHONY: help
help:
	@echo "Make targets"
	@echo "lint:        run golangci-lint"
	@echo "test:        run tests"

.PHONY: all
all: lint test

.PHONY: lint
lint:
	golangci-lint run -E golint ./...

.PHONY: test
test:
	go test -coverprofile coverage.out -race -v ./...

.PHONY: coverage
coverage: test
	go tool cover -html=coverage.out
