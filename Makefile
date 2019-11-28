GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOTIDY=$(GOCMD) mod tidy
BINARY_NAME=web

.PHONY: all
all: tidy test build

GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_ARGS=-ldflags "-X main.version=$(GIT_COMMIT)"

.PHONY: build
build:
	$(GOBUILD) $(BUILD_ARGS) -o $(BINARY_NAME) -v cmd/web/*

.PHONY: test
test:
	$(GOTEST) -v -race ./...

pkged.go:
	pkger

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: tidy
tidy:
	$(GOTIDY)
