package ws

import (
	"encoding/json"
	"errors"
)

var ErrEmptyPacket = errors.New("empty packet")

// Bind 将请求包中的 data 字段绑定到指定请求结构。
func Bind[T any](ctx *Context) (*T, error) {
	var req T
	if ctx == nil || ctx.Packet == nil {
		return nil, ErrEmptyPacket
	}
	if len(ctx.Packet.Data) == 0 {
		return &req, nil
	}
	if err := json.Unmarshal(ctx.Packet.Data, &req); err != nil {
		Fail(ctx, 400, err.Error())
		return nil, err
	}
	return &req, nil
}
