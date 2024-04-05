# GO path
GOPATH := $(shell go env GOPATH)

generate-go-wire:
	@wire ./cmd

generate: generate-go-wire

run: 
	@go run ./cmd/...