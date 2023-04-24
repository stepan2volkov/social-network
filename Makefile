BIN_DIR := ${PWD}/bin

.PHONY: install-tools
install-tools:
	GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

.PHONY: lint
lint:
	$(BIN_DIR)/golangci-lint run -c golangci-lint.yaml profile/...

