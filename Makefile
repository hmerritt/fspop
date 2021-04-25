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

bootstrap:
    go get github.com/mitchellh/gox
	go generate -tags tools tools/tools.go
.PHONY: bootstrap

buildq:
	go build -ldflags "-s -w" .
.PHONY:bq

build: vet
	gox -osarch "linux/386 linux/amd64 linux/arm linux/arm64 windows/386 windows/amd64" \
	-ldflags "-s -w" \
	-tags "fspop"    \
	-gocmd go        \
	-output "pkg/{{.OS}}_{{.Arch}}/fspop"  \
	.
.PHONY:build
