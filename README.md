## build micro serve by golang part 1

### 安装依赖
```bash
 go get -u google.golang.org/grpc
 go get -u github.com/golang/protobuf/protoc-gen-go
 brew install protobuf
```
### 运行程序
```bash
 make build
 go run main.go // server side
 go run cli.go     //  client side
```