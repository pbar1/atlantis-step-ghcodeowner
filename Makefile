export CGO_ENABLED     := 0
export DOCKER_BUILDKIT := 1

BIN     := $(shell basename $(PWD))
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"
IMAGE   := ghcr.io/pbar1/$(BIN)

build: clean
	GOOS=linux   GOARCH=arm64 go build -o bin/$(BIN)_linux_arm64  $(LDFLAGS) main.go
	GOOS=linux   GOARCH=amd64 go build -o bin/$(BIN)_linux_amd64  $(LDFLAGS) main.go
	GOOS=darwin  GOARCH=amd64 go build -o bin/$(BIN)_darwin_amd64 $(LDFLAGS) main.go
	du -sh bin/*

image: build
	docker build . -t $(IMAGE):$(VERSION) -t $(IMAGE):latest

image-push: image
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

clean:
	rm -rf bin
