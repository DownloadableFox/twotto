# GO path
GOPATH := $(shell go env GOPATH)

generate-go-wire:
	@wire ./cmd

generate: generate-go-wire

build: generate
	@mkdir -p bin
	@go build -o bin/bot ./cmd/...

# Use go-migrate to run migrations using the environment variables
migrate-up:
	@if [ -z ${DATABASE_URL} ]; then echo "DATABASE_URL is not set"; exit 1; fi
	@migrate -path ./migrations -database "${DATABASE_URL}" up

migrate-down:
	@if [ -z ${DATABASE_URL} ]; then echo "DATABASE_URL is not set"; exit 1; fi
	@migrate -path ./migrations -database "${DATABASE_URL}" down

run: migrate-up
	@go run ./cmd/...

bootstrap: migrate-up
	@./bin