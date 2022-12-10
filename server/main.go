package main

import (
	"core/core"
	"log"
)

func main() {
	gateway := core.NewGateway()
	gateway.App.SetBeforeFunc(func() {
		log.SetFlags(log.LstdFlags + log.Lshortfile)
	})
	gateway.Run("192.168.101.10", 8848)
	gateway.App.Stop(func() {
		gateway.Nacos.Logout()
		log.Println("关机操作!")
	})
}
