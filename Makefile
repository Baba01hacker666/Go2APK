GO ?= go
BINARY := go2apk
BIN_DIR := bin

.PHONY: all
all: verify build

## build: compile the CLI into bin/
.PHONY: build
build:
	$(GO) build -o $(BIN_DIR)/$(BINARY) ./cmd/$(BINARY)

## install: install the CLI into GOBIN
.PHONY: install
install:
	$(GO) install ./cmd/$(BINARY)

## test: run the test suite with the race detector
.PHONY: test
test:
	$(GO) test -race ./...

## cover: run the tests and report coverage per package
.PHONY: cover
cover:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out | tail -1

## verify: everything CI enforces - formatting, vet and tests
.PHONY: verify
verify: fmt-check vet test

## fmt: rewrite sources with gofmt
.PHONY: fmt
fmt:
	$(GO) fmt ./...

## fmt-check: fail if any file is not gofmt'd
.PHONY: fmt-check
fmt-check:
	@unformatted="$$($(GO)fmt -l .)"; \
	if [ -n "$$unformatted" ]; then \
		echo "These files are not gofmt'd:" >&2; \
		echo "$$unformatted" >&2; \
		exit 1; \
	fi

## vet: run go vet
.PHONY: vet
vet:
	$(GO) vet ./...

## generate: regenerate the demo app's Android project from go2apk.yaml
.PHONY: generate
generate:
	$(GO) run ./cmd/$(BINARY) generate

## preview: render the demo app's UI preview into preview.html
.PHONY: preview
preview:
	$(GO) run ./cmd/$(BINARY) preview

## doctor: report which build prerequisites are present
.PHONY: doctor
doctor:
	$(GO) run ./cmd/$(BINARY) doctor

## clean: remove build outputs and generated artifacts
.PHONY: clean
clean:
	rm -rf $(BIN_DIR) dist coverage.out preview.html
	$(GO) run ./cmd/$(BINARY) clean

## help: list the available targets
.PHONY: help
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## //' | sort | awk -F': ' '{printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
