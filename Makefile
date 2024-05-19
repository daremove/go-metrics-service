.PHONY: build lint format test generate

SERVER_SOURCE_PATH=./cmd/server
AGENT_SOURCE_PATH=./cmd/agent

VERSION_PATH=github.com/daremove/go-metrics-service/cmd/buildversion
VERSION=0.0.1

build:
	@echo "Building the project..."
	@go build -ldflags "-X '$(VERSION_PATH).BuildVersion=$(VERSION)' -X '$(VERSION_PATH).BuildDate=$(shell date +'%Y/%m/%d')' -X '$(VERSION_PATH).BuildCommit=$(shell git rev-parse HEAD)'" \
 		-o $(SERVER_SOURCE_PATH)/server $(SERVER_SOURCE_PATH)/*.go
	@go build -ldflags "-X '$(VERSION_PATH).BuildVersion=$(VERSION)' -X '$(VERSION_PATH).BuildDate=$(shell date +'%Y/%m/%d')' -X '$(VERSION_PATH).BuildCommit=$(shell git rev-parse HEAD)'" \
	 	-o $(AGENT_SOURCE_PATH)/agent $(AGENT_SOURCE_PATH)/*.go

lint:
	@go build -o ./cmd/staticlint/analyzer ./cmd/staticlint/. && ./cmd/staticlint/analyzer ./...

format:
	@goimports -l -w  .

test:
	@go test -count=1 -cover ./...

test-coverage:
	@cd scripts && ./test_coverage.sh

generate:
	@go generate ./...
