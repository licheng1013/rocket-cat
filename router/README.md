# 路由处理
## 介绍
- 2022/11/28

## 描述
- 此包是对路由的处理

## 自定义代理示例
- 必须实现 Proxy 接口

```go
type C struct {
	proxy Proxy
}

func (p *C) InvokeFunc(msg message.Message) []byte {
	router.FileLogger().Println("C执行")
	invokeFunc := p.proxy.InvokeFunc(msg)
	router.FileLogger().Println("C之后")
	return invokeFunc
}
func (p *C) SetProxy(proxy Proxy) {
	p.proxy = proxy
}

```