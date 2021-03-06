APP_BIN := "./bin/calendar"
MIGRATIONS_BIN := "./bin/goose"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(APP_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(APP_BIN) -config ./configs/.env

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(APP_BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

build_goose:
	go build -o $(MIGRATIONS_BIN) ./cmd/goose/*.go

migrate: build_goose
	./bin/goose --dir=migrations --config=configs/.env up

rollback: build_goose
	./bin/goose --dir=migrations --config=configs/.env down

.PHONY: tools build run build-img run-img version test lint
