GO      ?= go
BIN_DIR := bin
BINARY  := watchmud

# where `make test-db` looks for the mongo from docker-compose.yml
TEST_MONGO_URI ?= mongodb://localhost:27018

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
	@files=$$(gofmt -l .); \
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

## db-up: start the local mongo in docker
.PHONY: db-up
db-up:
	docker compose up -d mongo

## db-down: stop the local mongo, keeping its data
.PHONY: db-down
db-down:
	docker compose down

## db-reset: stop the local mongo and delete every character in it
.PHONY: db-reset
db-reset:
	docker compose down -v

## db-shell: open a mongosh against the local mongo
.PHONY: db-shell
db-shell:
	docker compose exec mongo mongosh watchmud

## test-db: run the tests that need a real mongo (make db-up first)
.PHONY: test-db
test-db:
	WATCHMUD_TEST_MONGO_URI=$(TEST_MONGO_URI) $(GO) test ./mongostore/ -count=1 -v

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
