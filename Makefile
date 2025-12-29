.PHONY:
	build test lint clean install fmt
install:
	go mod tidy
build:
	go build -o bin/gendiff ./cmd/gendiff
clean:
	rm -rf bin/
run:
	go run cmd/gendiff/main.go
lint:
	golangci-lint run ./...
test:
	go test -coverprofile=coverage.out ./...

