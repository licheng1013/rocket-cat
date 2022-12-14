# 如何编译消息

## 安装
- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

-  格式： protoc --go_out 目标路径 文件路径
- protoc --go_out ./  proto_message.proto
