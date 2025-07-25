NAME := squishy
VERSION := $(shell cat .version)
GOPATH := $(shell go env GOPATH)
LDFLAGS := "-w -s -X main.Version=$(VERSION)"

define build-bin
	$(info Building $(NAME) v$(VERSION) with arch $1 for $2)
	CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -v -o ./bin/$(NAME)-$(1)-$(2)-$(VERSION) -ldflags $(LDFLAGS) ./cmd/main.go
endef

define run-bin
	$(info Running $(NAME) v$(VERSION) with arch $1 for $2)
	./bin/$(NAME)-$(1)-$(2)-$(VERSION)
endef

air:
	@go build -o ./.air/main -ldflags $(LDFLAGS) ./cmd/main.go

install-dev-deps:
	$(info Installing dependencies...)
	go mod download
	go install github.com/air-verse/air@latest

dev:
	$(info Starting air dev process)
	$(GOPATH)/bin/air

run:
	$(info Running $(NAME))
	go run cmd/main.go

cleanup:
	rm ./bin/* -rf

build-linux-amd64:
	$(call build-bin,linux,amd64)

build-linux-arm64:
	$(call build-bin,linux,arm64)

build-all: cleanup build-linux-arm64 build-linux-amd64

run-linux-amd64:
	$(call run-bin,linux,amd64)

run-linux-arm64:
	$(call run-bin,linux,arm64)

create-release: build-all
	glab release create v$(VERSION)
	git fetch --tags origin
	glab release upload v$(VERSION) ./bin/*
