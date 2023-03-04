package router

import (
	"github.com/licheng1013/io-game-go/message"
	"github.com/licheng1013/io-game-go/remote"
)

type Context struct {
	// 具体消息 -> 当这个消息为空时则不返回数据回去
	Message message.Message
	// Rpc服务
	RpcServer remote.RpcServer
}
