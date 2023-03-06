package registers

import "fmt"

// Register 注册中心必须实现的接口！
type Register interface {
	// Register 注册
	Register(RegisterInfo)
	// RegisterClient 注册客户端信息
	RegisterClient(RegisterInfo)
	// GetIp 获取1个ip
	GetIp() (RegisterInfo, error)
	// ListIp 获取所有ip
	ListIp() ([]RegisterInfo, error)
	// Close 用于关机等操作
	Close()
	// RegisterInfo 注册信息
	RegisterInfo() RegisterInfo
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

func (v RegisterInfo) Addr() string {
	return v.Ip + ":" + fmt.Sprint(v.Port)
}
