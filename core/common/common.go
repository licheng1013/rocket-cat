package common

import (
	"errors"
	"reflect"
)

// 断言工具

func AssertErr(err error) {
	if err != nil {
		print(err)
	}
}

func AssertPtr(v interface{}, errInfo string) (err error) {
	pv1 := reflect.ValueOf(v)
	if pv1.Kind() != reflect.Ptr {
		err = errors.New(errInfo)
		return
	}
	return nil
}

const (
	GatewayName  = "gateway-game"
	ServicerName = "service-game"
)

// Closeing 关闭接口
type Closeing interface {
	Close()
}
