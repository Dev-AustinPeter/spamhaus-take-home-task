
.PHONY: build run clean test

BINARY_NAME=bin/spamhaus-take-home-task

build:
	@go build -o $(BINARY_NAME) cmd/main.go

run: build
	@./$(BINARY_NAME)

clean:
	@rm -f $(BINARY_NAME)

test:
	@go test -v ./...
