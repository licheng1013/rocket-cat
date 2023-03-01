package router

import (
	"github.com/io-game-go/message"
	"github.com/io-game-go/remote"
)

type Context struct {
	// 具体消息
	Message message.Message
	// Rpc服务
	RpcServer remote.RpcServer
}
