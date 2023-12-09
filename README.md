# gRPC Reddit

A minimal gRPC application that mimics the functionality of Reddit.

This is the "A3 - gRPC" assignment for the 17-625 API Design course at Carnegie Mellon University.

https://github.com/tomy0000000/grpc-reddit/assets/23290356/be9165e5-d3e5-4b14-a785-40f837773a3a

## Commands

- Generate gRPC code

```shell
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    Reddit/reddit.proto
```

- Run server

```shell
go run ./server
```

- Run client

```shell
go run ./client
```

- Run tests

```shell
go test ./...
```
