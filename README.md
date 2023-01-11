# IoGameGo
## 介绍
- 2022/10/14
- 目前还是一个实验性项目

## 描述
- 一个go简单游戏服务器实现，目前是的。
- 打造一个简单的游戏服务器框架，易扩展，易使用。
- 网关与逻辑服通过Grpc进行调试

## 架构图
- ![struct.png](struct.png)

## 功能
### 传输结构
- [x] 支持Json
- [x] 支持Proto

### 连接协议
- [x] 支持Kcp
- [ ] 支持Tcp
- [ ] 支持Websocket

### 负载均衡
- [x] Ncaos已提供

## 注册中心  
### Nacos
- 单机启动： startup.cmd -m standalone
- [Nacos](https://nacos.io/zh-cn/docs/v2/quickstart/quick-start.html)

## 工具类
- [Lancet](https://github.com/duke-git/lancet/blob/main/README_zh-CN.md)

## 如何调试
- 运行 nacos 单机
- 运行 gateway 网关
- 运行 service 服务
- 运行 example 客户端