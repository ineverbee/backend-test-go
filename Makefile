.PHONY: build
build:
	go build -v ./cmd/auto

.PHONY: run
run: build
	./auto

.DEFAULT_GOAL := build
