# VERSION=$(shell git describe --tags --always)
VERSION=1.0.0
SERVICE_NAME ?= $(shell basename $(CURDIR))

.PHONY: build
# build
build:
	mkdir -p bin/
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./cmd/...

.PHONY: generate
# generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY: lint
# lint
lint:
	golangci-lint run -v --timeout 10m

.PHONY: test
# test
test:
	go test -v ./... -count=1

.PHONY: migrate-create
# migrate-create: create migration SQL file with name={file_name}
migrate-create:
	goose -dir=./migrations create $(name) sql


# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help