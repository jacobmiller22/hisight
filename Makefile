.DEFAULT:
	build

.PHONY: fmt vet clean build build-proto build-hisight-hook build-hisight-hook build-hisight-log-server build-hisight-server test

build: build-proto build-hisight-hook build-hisight-log-server build-hisight-server

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/history.proto

build-hisight-hook: vet
	go build -o ./bin/hisight-hook ./cmd/client-hook

build-hisight-log-server: vet
	go build -o ./bin/hisight-log-server ./cmd/log-server

build-hisight-server: vet
	go build -o ./bin/hisight-server ./cmd/hisight-server

clean:
	rm -rf ./bin/*


test: 
	go test -v ./...

