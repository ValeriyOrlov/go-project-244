.PHONY:
	build test lint clean install fmt
install:
	go mod tidy
build:
	go build -o bin/gendiff ./cmd/main.go
clean:
	rm -rf bin/
run:
	go run cmd/main.go
lint:
	golangci-lint run ./...
test:
	go test -coverprofile=bin/coverage.out ./...
fmt:
	gofmt -w .

