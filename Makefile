BINARY_NAME := online_checker
WIN_BINARY_NAME := $(BINARY_NAME).exe
VERSION_PATH := github.com/artkescha/$(BINARY_NAME)
BINARY_VERSION := $(shell git describe --tags --always)
BINARY_BUILD_DATE := $(shell date +%d.%m.%Y)
BUILD_FOLDER := .build

PRINTF_FORMAT := "\033[35m%-18s\033[33m %s\033[0m\n"

.PHONY: all build windows linux etc-windows etc-linux vendor lint test clean help

all: lint test vendor build

build: windows linux ## Default: build for windows and linux

windows: etc-windows $(BUILD_FOLDER)/windows/config.yaml vendor ## Build artifacts for windows
	@printf $(PRINTF_FORMAT) BINARY_NAME: $(WIN_BINARY_NAME)
	@printf $(PRINTF_FORMAT) BINARY_VERSION: $(BINARY_VERSION)
	@printf $(PRINTF_FORMAT) BINARY_BUILD_DATE: $(BINARY_BUILD_DATE)
	@printf $(PRINTF_FORMAT) VERSION_PATH: $(VERSION_PATH)
	mkdir -p $(BUILD_FOLDER)/windows
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc  go build  -ldflags "-s -w -X $(VERSION_PATH).Version=$(BINARY_VERSION) -X $(VERSION_PATH).BuildDate=$(BINARY_BUILD_DATE)" -o $(BUILD_FOLDER)/windows/$(WIN_BINARY_NAME) .

etc-windows:
	mkdir -p $(BUILD_FOLDER)/windows
	cp -r etc/* $(BUILD_FOLDER)/windows/

$(BUILD_FOLDER)/windows/config.yaml:
	@mkdir -p $(BUILD_FOLDER)/windows
	generate-logconfig -service-name $(BINARY_NAME) -os windows > $(BUILD_FOLDER)/windows/config.yaml

linux: etc-linux $(BUILD_FOLDER)/linux/config.yaml vendor ## Build artifacts for linux
	@printf $(PRINTF_FORMAT) BINARY_NAME: $(BINARY_NAME)
	@printf $(PRINTF_FORMAT) BINARY_VERSION: $(BINARY_VERSION)
	@printf $(PRINTF_FORMAT) BINARY_BUILD_DATE: $(BINARY_BUILD_DATE)
	@printf $(PRINTF_FORMAT) VERSION_PATH: $(VERSION_PATH)
	mkdir -p $(BUILD_FOLDER)/linux
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w -X $(VERSION_PATH).Version=$(BINARY_VERSION) -X $(VERSION_PATH).BuildDate=$(BINARY_BUILD_DATE)" -o $(BUILD_FOLDER)/linux/$(BINARY_NAME) .

etc-linux:
	mkdir -p $(BUILD_FOLDER)/linux
	cp -r etc/* $(BUILD_FOLDER)/linux/

vendor: ## Get dependencies according to go.mod
	env GO111MODULE=auto go mod vendor

test: vendor ## Start unit-tests
	go test ./...

lint: vendor ## Start static code analysis
	hash golangci-lint 2>/dev/null || go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --timeout=5m

clean: ## Remove vendor and artifacts
	rm -rf vendor
	rm -rf $(BUILD_FOLDER)

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'