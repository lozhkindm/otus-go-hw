APP_BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/scheduler"
SENDER_BIN := "./bin/sender"
MIGRATIONS_BIN := "./bin/goose"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(APP_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(SCHEDULER_BIN) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(SENDER_BIN) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(APP_BIN) -config ./configs/.env

run-scheduler: build
	$(SCHEDULER_BIN) -config ./configs/.env.scheduler

run-sender: build
	$(SENDER_BIN) -config ./configs/.env.sender

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

generate:
	protoc proto/calendar.proto --go_out=./internal/server/grpc/pb --go-grpc_out=./internal/server/grpc/pb

.PHONY: tools build run build-img run-img version test lint
