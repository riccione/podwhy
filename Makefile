BINARY_NAME=podwhy
BIN_DIR=bin
VERSION?=dev

.PHONY: build build-all clean test fmt vet

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/podwhy

build-all:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/podwhy
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/podwhy
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/podwhy
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/podwhy

clean:
	rm -rf $(BIN_DIR)

test:
	go test ./... -v

fmt:
	go fmt ./...

vet:
	go vet ./...
