package router

type Proxy interface {
	// InvokeFunc 调用函数
	InvokeFunc(ctx Context) []byte
	// SetProxy 默认目标对象是A 你编写了B代理 -> 那么A将会传递下去,所以你在代理类需要调用代理对象的方法
	SetProxy(proxy Proxy)
}

// ProxyFunc 代理模型
type ProxyFunc struct {
	proxy Proxy
}

func (p *ProxyFunc) InvokeFunc(ctx Context) []byte {
	return p.proxy.InvokeFunc(ctx)
}

func (p *ProxyFunc) SetProxy(proxy Proxy) {
	p.proxy = proxy
}
