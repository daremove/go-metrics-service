.PHONY: lint format test generate

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
