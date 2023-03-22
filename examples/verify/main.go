package main

import (
	"fmt"

	"github.com/licheng1013/rocket-cat/common"
)

func main() {
	sm := common.NewSafeMap()
	sm.Set(1, 10)
	sm.Set(2, "hello")
	sm.Set(3, true)
	fmt.Println(sm.Keys())
	fmt.Println(sm.Get(2))
	fmt.Println(sm.Values())
}
