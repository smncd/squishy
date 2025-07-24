NAME := squishy
GOPATH := $(shell go env GOPATH)

define build
	$(info Building $(NAME) with arch $1 for $2)
	CGO_ENABLED=0 GOOS=$2 GOARCH=$1 go build -v -o ./bin/$(NAME)-$(1)-$(2) -ldflags "-w -s" ./cmd/main.go
endef

install-deps:
	$(info Installing dependencies...)
	go mod download
	go install github.com/air-verse/air@latest

dev:
	$(info Starting air dev process)
	$(GOPATH)/bin/air

run:
	$(info Running $(NAME))
	go run cmd/main.go

build-amd64-linux:
	$(call build,amd64,linux)

build-arm64-linux:
	$(call build,arm64,linux)

build-all: build-amd64-linux build-arm64-linux
