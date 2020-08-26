.PHONY: help
help:
	@echo "Make targets"
	@echo "lint:        run golangci-lint"
	@echo "test:        run tests"

.PHONY: lint
lint:
	golangci-lint run -E golint ./...

.PHONY: test
test:
	go test -race -v ./...
