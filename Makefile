GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOTIDY=$(GOCMD) mod tidy
BINARY_NAME=web
DEPLOY_ASSETS=web_unix ui resume.conf
DEPLOY_TARGET=www.dennis-schoen.de:dennis-schoen-de/

.PHONY: all
all: tidy build

GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_ARGS=-ldflags "-X main.version=$(GIT_COMMIT)"

.PHONY: build
build:
	$(GOBUILD) $(BUILD_ARGS) -o $(BINARY_NAME) -v cmd/web/*

# Cross compilation
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_ARGS) -o  $(BINARY_NAME)_unix -v cmd/web/*

.PHONY: test
test:
	$(GOTEST) -v -race ./...

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: tidy
tidy:
	$(GOTIDY)

.PHONE: deploy
deploy:
	rsync -r $(DEPLOY_ASSETS) $(DEPLOY_TARGET)
