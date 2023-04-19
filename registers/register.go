package registers

import "fmt"

// IRegister 注册中心必须实现的接口
type IRegister interface {
	// GetIp 获取单个ip
	GetIp() (ClientInfo, error)
	// ListIp 获取所有ip
	ListIp(serverName string) ([]ClientInfo, error)
	// Close 当服务器关闭时调用
	Close()
	// ClientInfo  注册信息
	ClientInfo() ClientInfo
	// Run 此处用于启动注册中心，如etcd，consul等，你可以在这里设置心跳等操作。
	Run()
	// RegisterServer Register 注册
	RegisterServer(server ServerInfo)
	// RegisterClient 注册客户端信息
	RegisterClient(client ClientInfo)
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
