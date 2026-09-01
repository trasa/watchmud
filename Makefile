GO      ?= go
BIN_DIR := bin
BINARY  := watchmud

.DEFAULT_GOAL := build

## build: compile server into bin/watchmud
.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BINARY) ./cmd/watchmud

## test: run all tests
.PHONY: test
test:
	$(GO) test ./...

## vet: run static analysis
.PHONY: vet
vet:
	$(GO) vet ./...

## fmt: rewrite source files to canonical formatting
.PHONY: fmt
fmt:
	$(GO) fmt ./...

## fmt-check: fail if any files need formatting (CI)
.PHONY: fmt-check
fmt-check:
	@files=$$(gofmt -files .); \
	if [ -n "$$files" ]; then \
	  echo "not gofmt'd:"; echo "$$files"; \
	  exit 1; \
	fi

## generate: regenerate string output, etc. (run manually after editing enums)
.PHONY: generate
generate:
	$(GO) generate ./...

## tidy: prune and sync go.mod / go.sum
.PHONY: tidy
tidy:
	$(GO) mod tidy

## check: everything CI should enforce
.PHONY: check
check: fmt-check vet test

## all: the full local workflow
.PHONY: all
all: check build

## run: build and start the server with example config
.PHONY: run
run: build
	$(BIN_DIR)/$(BINARY)

## clean: remove build output and cached test results
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
	$(GO) clean -testcache

## help: list available targets
.PHONY: help
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## / /'
