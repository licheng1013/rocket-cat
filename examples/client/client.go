package main

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/router"
)

func main() {
	// Create SprintXxx functions to mix strings with other non-colorized strings:
	router.LogFunc(1, Hello)
	fmt.Println("客户端")
	router.StartLogo()
}

func Hello(ctx *router.Context) {

}
