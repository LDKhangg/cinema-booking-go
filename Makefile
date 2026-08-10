GOPATH := $(shell go env GOPATH)
SWAG   := $(GOPATH)/bin/swag

.PHONY: docs run build

docs:
	$(SWAG) init -g cmd/main.go -o docs

run: docs
	go run ./cmd/main.go

build:
	go build -o bin/app ./cmd/main.go
