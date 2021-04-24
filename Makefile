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

buildq:
	go build -ldflags "-s -w" .
.PHONY:bq

build_linux:
	gox -os "linux" -arch "amd64" -ldflags "-s -w" -output "fspop" .
.PHONY:build_linux
