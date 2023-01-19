NAME = $(shell basename $(PWD))
VERSION = $(shell git describe --tags --always --dirty)
COMMIT = $(shell git rev-parse --short HEAD)
DATE = $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

build:
	go build -o bin/$(NAME) -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

run: build
	./bin/$(NAME)
