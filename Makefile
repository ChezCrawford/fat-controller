
BINARY_NAME := fat-controller

.PHONY: compile
compile: compile_darwin

.PHONY: compile_all
compile_all: compile_linux compile_linux_arm compile_darwin

.PHONY: compile_linux
compile_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/linux_amd64/${BINARY_NAME} ./cmd/

.PHONY: compile_linux_arm
compile_linux_arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/linux_arm/${BINARY_NAME} ./cmd/

.PHONY: compile_darwin
compile_darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin_amd64/${BINARY_NAME} ./cmd/

.PHONY: test
test:
	go test ./...
