NAME := squishy
VERSION := $(shell cat .version)
GOPATH := $(shell go env GOPATH)

define build
	$(info Building $(NAME) v$(VERSION) with arch $1 for $2)
	CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -v -o ./bin/$(NAME)-$(1)-$(2)-$(VERSION) -ldflags "-w -s" ./cmd/main.go
endef

install-dev-deps:
	$(info Installing dependencies...)
	go mod download
	go install github.com/air-verse/air@latest

generate:
	go generate ./...

dev: generate
	$(info Starting air dev process)
	$(GOPATH)/bin/air

run:
	$(info Running $(NAME))
	go run cmd/main.go

cleanup:
	rm ./bin/* -rf

build-linux-amd64:
	$(call build,linux,amd64)

build-linux-arm64:
	$(call build,linux,arm64)

build-all: generate cleanup build-linux-arm64 build-linux-amd64

create-release: build-all
	glab release create v$(VERSION)
	glab release upload v$(VERSION) ./bin/*
