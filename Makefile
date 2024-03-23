.PHONY: build
build:
	docker compose build

.PHONY: start
start:
	docker compose up -d

.PHONY: down
down:
	docker compose down -v

.PHONY: test
test:
	go test -race ./...

.PHONY: lint
lint: golangci-lint nilaway-lint

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run

.PHONY: nilaway-lint
nilaway-lint:
	nilaway ./...