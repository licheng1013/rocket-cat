package registers

import "fmt"

// IRegister 注册中心必须实现的接口
type IRegister interface {
	GetIp() (ClientInfo, error)
	ListIp(serverName string) ([]ClientInfo, error)
}

// Register 注册中心必须实现的接口！
type Register interface {
	// Register 注册
	Register(info ServerInfo)
	// RegisterClient 注册客户端信息
	RegisterClient(ClientInfo)
	// GetIp 获取1个ip
	GetIp() (ClientInfo, error)
	// ListIp 获取所有ip
	ListIp(serverName string) ([]ClientInfo, error)
	// Close 用于关机等操作
	Close()
	// ClientInfo  注册信息
	ClientInfo() ClientInfo
}

// ClientInfo 客户端注册的数据
type ClientInfo struct {
	// ip
	Ip string
	// 端口
	Port uint16
	// 当前注册服务名
	ServiceName string
	// 远程注册服务名 -> 在service则是gateway的注册名，在gateway则相反
	RemoteName string
}

func (v ClientInfo) Addr() string {
	return v.Ip + ":" + fmt.Sprint(v.Port)
}

// ServerInfo 注册中心注册的数据
type ServerInfo struct {
	// ip
	Ip string
	// 端口
	Port uint16
}

func (v ServerInfo) Addr() string {
	return v.Ip + ":" + fmt.Sprint(v.Port)
}
