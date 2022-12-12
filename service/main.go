package main

import "core/core"

func main() {
	service := core.NewService(8999)
	service.Run("192.168.101.10", 8848)
}
