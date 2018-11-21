# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /bin/bash

.PHONY: all
all: fmt vet install

.PHONY: install
install:
	go install ./...

.PHONY: fmt
fmt:
	goimports -w .
	gofmt -s -w .

.PHONY: vet
vet:
	go vet .

.PHONY: test
test:
	go test -v ./...

.PHONE: mod
mod:
	go mod init github.com/sadlil/failsafe