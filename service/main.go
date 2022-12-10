package main

import (
	"core/core"
	"log"
)

func main() {
	gateway := core.NewService()
	gateway.App.SetProt(10000)
	gateway.App.SetBeforeFunc(func() {
		log.SetFlags(log.LstdFlags + log.Lshortfile)
	})
	gateway.Run("192.168.101.10", 8848)
}
