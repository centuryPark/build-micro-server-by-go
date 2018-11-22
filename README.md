## build micro serve by golang part 1

### 更改摘要
- 使用 go-micro 代替 grpc
- 把server 和 client 都容器化

### 安装依赖
```bash
 go get -u github.com/micro/protobuf/proto
 go get -u github.com/micro/protobuf/protoc-gen-go
 brew install protobuf
```
### 运行程序
```bash
 make build
 go run main.go // server side
 go run cli.go     //  client side
```