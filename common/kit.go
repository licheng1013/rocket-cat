package common

import "github.com/google/uuid"

var CmdKit = &cmdKit{}

type cmdKit struct {
}

// GetMerge 获取路由
func (c cmdKit) GetMerge(cmd, subCmd int64) int64 {
	return (cmd << 16) + subCmd
}

// GetCmd 获取主命令
func (c cmdKit) GetCmd(merge int64) int64 {
	return merge >> 16
}

// GetSubCmd 获取子命令
func (c cmdKit) GetSubCmd(merge int64) int64 {
	return merge & 0xFFFF
}
// 获取uuid
type uuidKit struct {
}

func (k uuidKit) UUID() uint32 {
	return uuid.New().ID()
}

var UuidKit = &uuidKit{}
