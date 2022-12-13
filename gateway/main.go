package main

import (
	"core/core"
	"log"
)

func main() {
	ip := "192.168.101.10"
	gateway := core.NewGateway(ip,8001)
	gateway.App.SetBeforeFunc(func() {
		log.SetFlags(log.LstdFlags + log.Lshortfile)
	})
	gateway.Run(ip, 8848)
}
