# Get the git commit
GIT_COMMIT := $(git rev-parse HEAD)
GIT_DIRTY := $(test -n `git status --porcelain`" && echo "+CHANGES" || true)
#GIT_DIRTY := ""

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
.PHONY:bootstrap

buildq:
	go build -ldflags "-s -w" .
.PHONY:buildq

buildlq:
	gox -osarch "linux/amd64" -ldflags "-s -w" -tags "fspop" -gocmd go -output fspop .
.PHONY:buildlq

build: vet
	gox -osarch "darwin/amd64 linux/386 linux/amd64 linux/arm linux/arm64 windows/386 windows/amd64" \
	-ldflags "-s -w" \
	-tags "fspop"    \
	-gocmd go        \
	-output "pkg/{{.OS}}_{{.Arch}}/fspop"  \
	.
.PHONY:build
