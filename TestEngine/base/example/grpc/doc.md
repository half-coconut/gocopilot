### Protobuf 文档

https://protobuf.dev/programming-guides/proto3/

- 安装 protobuf

```shell
brew install protobuf
protoc --version
```

- install the protocol compiler plugins for Go using the following commands:

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

```

生成命令

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user.proto
```

- 生成 `user.pb.go` -> go 代码
- 生成 `user_grpc.pb.go` -> grpc 代码，客户端和服务端


使用 buf 工具，方便编译
```shell
npm install @bufbuild/buf

npx buf --version

➜  TestEngine git:(main) ✗ npx buf generate api/proto
```