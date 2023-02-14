package registers

// Register 注册中心必须实现的接口！
type Register interface {
	// Register 注册
	Register(RegisterInfo)
	// GetIp 获取1个ip
	GetIp() RegisterInfo
	// ListIp 获取所有ip
	ListIp() []RegisterInfo
}

type RegisterInfo struct {
	// ip
	Ip string
	// 端口
	Port uint16
	// 当前注册服务名
	ServiceName string
	// 远程注册服务名 -> 在service则是gateway的注册名，在gateway则相反
	RemoteName string
}
