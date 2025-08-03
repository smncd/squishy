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

.PHONY: air
air:
	@go build -o ./.air/main -ldflags $(LDFLAGS) ./cmd/main.go

.PHONY: install-dev-deps
install-dev-deps:
	$(info Installing dependencies...)
	go mod download
	go install github.com/air-verse/air@latest

.PHONY: dev
dev:
	$(info Starting air dev process)
	@$(GOPATH)/bin/air

.PHONY: run
run:
	$(info Running $(NAME))
	go run cmd/main.go

.PHONY: cleanup
cleanup:
	rm ./bin/* -rf

.PHONY: build-linux-amd64
build-linux-amd64:
	$(call build-bin,linux,amd64)

.PHONY: build-linux-arm64
build-linux-arm64:
	$(call build-bin,linux,arm64)

.PHONY: build-all
build-all: cleanup build-linux-arm64 build-linux-amd64

.PHONY: run-linux-amd64
run-linux-amd64:
	$(call run-bin,linux,amd64)

.PHONY: run-linux-arm64
run-linux-arm64:
	$(call run-bin,linux,arm64)

.PHONY: build-docker-image
publish-docker-image:
	docker build --platform linux/amd64,linux/arm64 . -t registry.gitlab.com/smncd/squishy:$(VERSION)
	docker push registry.gitlab.com/smncd/squishy

.PHONY: create-release
create-release: build-all
	glab release create v$(VERSION)
	git fetch --tags origin
	glab release upload v$(VERSION) ./bin/*
