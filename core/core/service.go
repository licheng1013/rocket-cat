package core

// Service 构建服务
type Service struct {
	Nacos Nacos
}

func (n *Service) Register(port uint64) {
	n.Nacos.SetServerConfig("192.168.101.10", 8848)
	n.Nacos.Register("192.168.101.10", port, ServerName)
}
