.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

build: vet
	go build -ldflags "-s -w" .
.PHONY:build

buildquick:
	go build -ldflags "-s -w" .
.PHONY:bq
