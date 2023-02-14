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
	Ip   string
	Port uint16
	Name string
}
