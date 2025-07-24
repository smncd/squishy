NAME := squishy
VERSION := 0.4.0-beta.1
GOPATH := $(shell go env GOPATH)

define build
	$(info Building $(NAME) v$(VERSION) with arch $1 for $2)
	CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -v -o ./bin/$(NAME)-$(1)-$(2)-$(VERSION) -ldflags "-w -s" ./cmd/main.go
endef

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

build-linux-amd64:
	$(call build,linux,amd64)

build-linux-arm64:
	$(call build,linux,arm64)

build-all: build-linux-arm64 build-linux-amd64

create-release-tag:
	git tag $(NAME)-v$(VERSION) -m "$(NAME) $(VERSION)"

push-release-tag:
	git push origin $(NAME)-v$(VERSION)
