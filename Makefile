.PHONY: build test fmt lint release run install

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/gobp ./cmd/gobp

run: build
	./bin/gobp

test:
	go test ./...

fmt:
	gofmt -s -w .

lint:
	go vet ./...

install:
	go install $(LDFLAGS) ./cmd/gobp

release:
	@if [ -z "$(TAG)" ]; then echo "Usage: make release TAG=v0.1.0"; exit 1; fi
	git tag -a $(TAG) -m "Release $(TAG)"
	git push origin $(TAG)
