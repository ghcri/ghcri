SHELL := /bin/bash

.PHONY: all check format vet build test generate tidy release

-include Makefile.env

VERSION := v0.0.1

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GO_BUILD := go build -ldflags "-X main.Version=${VERSION}"

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check               to do static check"
	@echo "  build               to create bin directory and build"
	@echo "  test                to run test"

format:
	go fmt ./...

vet:
	go vet ./...

generate:
	go generate ./...

build: tidy generate format vet
	${GO_BUILD} -o bin/ghcri ./cmd/ghcri

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...
	go tool cover -html="coverage.txt" -o "coverage.html"

tidy:
	go mod tidy
	go mod verify
