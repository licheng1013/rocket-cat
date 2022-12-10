package main

import (
	"core"
	"log"
)

func main() {
	gateway := core.NewGateway()
	gateway.Run("192.168.101.10", 8848)
	gateway.Server.Stop(func() {
		gateway.Nacos.Logout()
		log.Println("关机操作!")
	})
}
