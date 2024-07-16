build:
	@go build -o ./bin/fs

run: build
	@./bin/fs

test:
	@go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := run