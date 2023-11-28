
.PHONY: compile
compile:
	mkdir -p ./bin
	go build -mod=readonly -o ./bin/fat-controller ./cmd/

.PHONY: test
test:
	go test ./...
