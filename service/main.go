package main

import "core/core"

func main() {
	ip := "192.168.101.10"
	service := core.NewService(ip,8999)
	service.Run(ip, 8848)
}
