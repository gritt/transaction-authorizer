APP_NAME=authorizer

.SILENT:
.DEFAULT_GOAL := help

.PHONY: help
help:
	$(info Available Commands:)
	$(info -> install                 install dependencies)
	$(info -> test                    run all tests with coverage)
	$(info -> build                   generates a build)
	$(info -> build-linux             generates a build for linux)

.PHONY: install
install:
	go mod tidy

.PHONY: test
test: install
	go test -race ./... -cover -v

.PHONY: build
build: install
	go build -o ${APP_NAME} ./cmd

.PHONY: build-linux
build-linux: install
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${APP_NAME}-tux ./cmd
