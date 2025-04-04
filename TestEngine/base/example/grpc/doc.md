
protobuf 文档

https://protobuf.dev/programming-guides/proto3/

```shell
brew install protobuf
protoc --version
```

- install the protocol compiler plugins for Go using the following commands:
```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

```

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user.proto

```