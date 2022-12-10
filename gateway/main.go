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
	gateway.App.EnableMessageLog = true
	gateway.Run("192.168.101.10", 8848)
}
