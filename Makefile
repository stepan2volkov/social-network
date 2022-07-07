BUILD_COMMIT := $(shell git log --format="%H" -n 1)
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%S')
FLAGS = GOOS=linux GOARCH=amd64 CGO_ENABLED=0

PROJECT = github.com/stepan2volkov/social-network
CMD:= $(PROJECT)/cmd/social-network
BIN:= ${PWD}/bin

OPENAPI_GEN:=${PWD}/scripts/generate-openapi

.PHONY: install-tools
install-tools:
	$(info Installing oapi-codegen into ./bin folder)
	@mkdir -p ./bin
	@GOBIN=${BIN} go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.11.0
	@GOBIN=${BIN} go install github.com/pressly/goose/v3/cmd/goose@v3.6.1

.PHONY: generate
generate:
	$(info Generating openapi spec)
	@./bin/oapi-codegen -config ${OPENAPI_GEN}/generate.openapi.auth.yaml ./api/openapi-spec/auth.openapi.yaml
	@./bin/oapi-codegen -config ${OPENAPI_GEN}/generate.openapi.profiles.yaml ./api/openapi-spec/profiles.openapi.yaml
	@./bin/oapi-codegen -config ${OPENAPI_GEN}/generate.openapi.friends.yaml ./api/openapi-spec/friends.openapi.yaml

.PHONY: local-migrate
local-migrate:
	@./bin/goose -dir ./migrations mysql "social-network-srv:1qaz@/social-network?parseTime=true" up

.PHONY: build
build:
	$(FLAGS) go build -a -tags netgo -ldflags="\
		-w -extldflags '-static'\
		-X '$(PROJECT)/internal/config.BuildCommit=$(BUILD_COMMIT)'\
		-X '${PROJECT}/internal/config.BuildTime=${BUILD_TIME}'"\
		-o bin $(CMD)

run:
	go run cmd/social-network/main.go --local