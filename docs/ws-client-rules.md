# WS 测试页面规则

适用文件：

```text
examples/ws-client.html
```

## 目标

测试页面用于模拟 Rocket Cat 的 WebSocket JSON 协议包，方便开发和验证 `cmd + subCmd` 路由、参数绑定、统一响应。

## 技术约束

- 使用单文件 HTML。
- 使用 Vue CDN。
- 不引入构建工具。
- 不引入 npm 依赖。
- CSS、HTML、JavaScript 保持内嵌。
- 页面可直接用浏览器打开。

## 默认连接

默认 WebSocket 地址：

```text
ws://localhost:10100/ws
```

服务启动命令：

```bash
go run .
```

## 发包协议

页面发送 JSON 包：

```json
{
  "cmd": 1001,
  "subCmd": 1,
  "data": {}
}
```

字段规则：

- `cmd`：主命令。
- `subCmd`：子命令。
- `data`：业务请求体，必须是合法 JSON 对象。

## 返回协议

服务端返回 JSON 包：

```json
{
  "cmd": 1001,
  "subCmd": 1,
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

字段规则：

- `code = 0` 表示成功。
- `code != 0` 表示失败。
- `msg` 用于展示状态或错误信息。
- `data` 为业务响应体。

## 内置预设

测试页面必须保留以下预设，除非对应 Go 模块命令发生变化：

| 名称 | cmd | subCmd | data |
| --- | ---: | ---: | --- |
| 用户登录 | 1001 | 1 | `{"account":"demo","password":"123456"}` |
| 用户退出 | 1001 | 2 | `{"token":"token-demo"}` |
| 用户信息 | 1001 | 3 | `{"uid":10001}` |
| 发送聊天 | 2001 | 2 | `{"toUid":10002,"content":"hello"}` |
| 聊天列表 | 2001 | 5 | `{"limit":20}` |
| 创建房间 | 3001 | 2 | `{"name":"Lobby"}` |
| 房间列表 | 3001 | 5 | `{"limit":20}` |

## 页面功能

测试页面应支持：

- 连接 WebSocket。
- 断开 WebSocket。
- 显示连接状态。
- 选择发包预设。
- 编辑 `cmd`。
- 编辑 `subCmd`。
- 编辑 `data JSON`。
- 格式化 `data JSON`。
- 发送协议包。
- 显示发送日志。
- 显示接收日志。
- 显示错误日志。
- 清空日志。

## 日志规则

日志按最新在前展示。

日志类型：

- `open`：连接成功。
- `send`：发送请求包。
- `recv`：收到响应包。
- `error`：连接错误或 JSON 错误。
- `close`：连接关闭。

收到的 JSON 字符串应尽量格式化展示；如果不是合法 JSON，则原样展示。

## 交互规则

- 未连接时禁止发送。
- 已连接时禁止重复连接。
- 未连接时禁止断开。
- `data JSON` 非法时不发送请求，并写入错误日志。
- 切换预设时同步更新 `cmd`、`subCmd`、`data JSON`。

## 样式规则

- 页面保持轻量、清晰、偏工具型。
- 不做营销落地页。
- 不使用复杂动画。
- 控件要适合频繁调试。
- 移动端需要可用，窄屏时改为单列布局。

## 后续维护规则

修改测试页面前，先读取本文档。

当新增 Go 模块或命令时，需要同步更新：

- `examples/ws-client.html` 的 `presets`。
- 本文档的“内置预设”表格。
- README 中必要的示例说明。
