package ws

import (
	"encoding/json"
	"errors"
)

var ErrEmptyPacket = errors.New("empty packet")

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
