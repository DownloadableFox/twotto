# GO path
GOPATH := $(shell go env GOPATH)

generate-go-wire:
	@echo "Running Google wire"
	@go generate ./...