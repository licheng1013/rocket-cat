package router

import (
	"github.com/licheng1013/rocket-cat/common"
	"runtime/debug"
)

type ErrProxy struct {
	ProxyFunc
}

func (e *ErrProxy) InvokeFunc(ctx *Context) {
	// 捕获异常
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case *ServiceError:
				errInfo := err.(*ServiceError)
				ctx.Message.SetMessage(errInfo.Message)
				ctx.Message.SetBody(errInfo.Message)
				common.Logger().Println("业务异常 -> ", errInfo.Message)
				break
			default:
				ctx.Message = nil
				common.Logger().Println("系统异常 -> ", err)
			}
			debug.PrintStack()
		}
	}()
	e.Proxy.InvokeFunc(ctx)
}

type ServiceError struct {
	error
	Code    int
	Message string
}

func NewServiceError(code int, message string) *ServiceError {
	return &ServiceError{Code: code, Message: message}
}
