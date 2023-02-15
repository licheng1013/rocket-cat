# 如何编译消息
- 如果你需要进行的自己的设计消息处理那么你需要实现该接口
- **编译后需要重新实现: Message 接口 **

## 编译消息
- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

-  格式： protoc --go_out 目标路径 文件路径
- protoc --go_out=. --go_opt=paths=source_relative message.proto


## RpcInfo
- 编译
- protoc --go_out=. --go_opt=paths=source_relative rpc_info.proto

## RpcService
- 使用rpc的时需要设置包目录,与上面的不同还需要设置service输出
- protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative rpc_service.proto
