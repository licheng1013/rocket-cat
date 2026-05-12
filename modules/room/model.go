package room

// QueryReq 表示房间查询请求。
type QueryReq struct {
	RoomId int64 `json:"roomId"` // 房间 ID
}

// RoomResp 表示房间信息响应。
type RoomResp struct {
	RoomId int64  `json:"roomId"` // 房间 ID
	Name   string `json:"name"`   // 房间名称
	Owner  int64  `json:"owner"`  // 房主用户 ID
}

// CreateReq 表示房间创建请求。
type CreateReq struct {
	Name string `json:"name"` // 房间名称
}

// ListReq 表示房间列表请求。
type ListReq struct {
	Limit int `json:"limit"` // 返回数量上限
}

// ListResp 表示房间列表响应。
type ListResp struct {
	Rooms []*RoomResp `json:"rooms"` // 房间列表
}
