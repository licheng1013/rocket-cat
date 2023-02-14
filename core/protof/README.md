# 如何编译消息
- 如果你需要进行的自己的设计消息处理那么你需要实现该接口
- **编译后需要重新实现: Message 接口 **

## 编译消息
- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

-  格式： protoc --go_out 目标路径 文件路径
- protoc --go_out=. --go_opt=paths=source_relative message.proto


## Grpc
- 编译
- protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc.proto