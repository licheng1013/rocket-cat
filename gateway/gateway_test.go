package gateway

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestGateway(t *testing.T) {
	gateway := &Gateway{}
	gateway.Nacos.Register("x", 11)
	time.Sleep(5 * time.Second)
	gateway.Nacos.AllInstances()
	time.Sleep(30 * time.Second)
	gateway.Nacos.Logout()
}

func TestOk(t *testing.T) {
	a, err := A()
	fmt.Println(a)
	fmt.Println(err.Error())
}

func A() (a int, err error) {
	a = 1
	err = errors.New("错误")
	return
}

func B() (int, error) {
	return 1, errors.New("错误")
}
