.PHONY: build
build:
	docker compose build

.PHONY: start
start: build
	docker compose up -d

.PHONY: down
down:
	docker compose down -v

.PHONY: test
test:
	go test -race ./...

.PHONY: lint
lint: golangci-lint

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run
