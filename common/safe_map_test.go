package common

import (
	"fmt"
	"testing"
)

func TestSyncMap(t *testing.T) {
	sm := NewSafeMap()
	sm.Set(1, 10)
	sm.Set(2, "hello")
	sm.Set(3, true)
	fmt.Println(sm.Keys())
	fmt.Println(sm.Get(2))
	fmt.Println(sm.Values())
}
