.DEFAULT:
	build

.PHONY: fmt vet clean build build-amd64 build-proto build-hisight-hook build-hisight-hook build-hisight-log-server build-hisight-server build-bash-hook test

build: build-proto build-hisight-hook build-hisight-log-server build-hisight-server build-bash-hook

build-amd64:
	GOOS=linux GOARCH=amd64 make build

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build-proto:
	protoc --go_out=./internal/commands --go_opt=paths=source_relative --go-grpc_out=./internal/commands --go-grpc_opt=paths=source_relative proto/history.proto

build-hisight-hook: vet
	go build -o ./bin/hisight-hook ./cmd/client-hook

build-hisight-log-server: vet
	go build -o ./bin/hisight-log-server ./cmd/log-server

build-hisight-server: vet
	go build -o ./bin/hisight-server ./cmd/hisight-server

build-bash-hook: vet
	go build -o ./bin/bashHook ./cmd/shellhook


clean:
	rm -rf ./bin/*


test: 
	go test -v ./...

