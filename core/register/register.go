package register

// Register 注册中心必须实现的接口！
type Register interface {
	// RequestUrl 请求地址
	RequestUrl() RequestInfo
}

type RequestInfo struct {
	Ip   string
	Port uint64
}
