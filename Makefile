GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOTIDY=$(GOCMD) mod tidy
BINARY_NAME=web
DEPLOY_ASSETS=$(BINARY_NAME)_unix ui resume.conf
DEPLOY_TARGET=www.dennis-schoen.de:dennis-schoen-de/

.PHONY: all
all: tidy build

GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_ARGS=-ldflags "-X main.version=$(GIT_COMMIT)"

.PHONY: build
## build: build the application
build:
	${GOBUILD} ${BUILD_ARGS} -o ${BINARY_NAME} -v cmd/web/*.go

.PHONY: build-linux
## build-linux: cross-compile application for linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GOBUILD} ${BUILD_ARGS} -o  ${BINARY_NAME}_unix -v cmd/web/*

.PHONY: test
## test: runs go test with default values
test:
	${GOTEST} -v -race ./...

.PHONY: clean
## clean: cleans the binary
clean:
	rm -f ${BINARY_NAME}

.PHONY: tidy
## tidy: tidy go modules
tidy:
	${GOTIDY}

.PHONE: deploy
## deploy: deploys the application
deploy: build-linux
	rsync -r ${DEPLOY_ASSETS} ${DEPLOY_TARGET}

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
