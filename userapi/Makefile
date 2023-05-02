BIN_DIR := ${PWD}/bin
MIGRATION_DIR := ${PWD}/migrations
DSN := social-network-srv:1qaz@/social-network?parseTime=true

.PHONY: install-tools
install-tools:
	GOBIN=$(BIN_DIR) go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4
	GOBIN=$(BIN_DIR) go install github.com/pressly/goose/v3/cmd/goose@v3.10.0

.PHONY: generate
generate:
	@echo "User API"
	@$(BIN_DIR)/oapi-codegen --config api/openapi/userapi.config.yaml api/openapi/userapi.yaml


.PHONY: migration-create
migration-create:
	$(BIN_DIR)/goose -dir=$(MIGRATION_DIR) create ${NAME} sql
	$(BIN_DIR)/goose -dir=$(MIGRATION_DIR) fix

.PHONY: migration-up
migration-up:
	$(BIN_DIR)/goose -dir=$(MIGRATION_DIR) mysql "$(DSN)" up

.PHONY: migration-down
migration-down:
	$(BIN_DIR)/goose -dir=$(MIGRATION_DIR) mysql "$(DSN)" down


.PHONY: run
run:
	DSN=$(DSN) go run ./cmd/userapi

.PHONY: build
build:
	mkdir -p build
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -ldflags="\
		-w -extldflags '-static'"\
		-o build ./cmd/userapi
