.DEFAULT:
	build

.PHONY: fmt vet clean build build-amd64 build-proto build-sql build-hs clean test

build: build-proto build-hs

build-amd64:
	GOOS=linux GOARCH=amd64 make build

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/commands/protocol/pb/commands.proto

build-sql:
	sqlc generate

build-hs: vet
	go build -o ./bin/hslog ./cmd/hslog

clean:
	rm -rf ./bin/*

test: 
	go test -v ./...

