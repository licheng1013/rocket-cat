package room

type QueryReq struct {
	RoomId int64 `json:"roomId"`
}

type RoomResp struct {
	RoomId int64  `json:"roomId"`
	Name   string `json:"name"`
	Owner  int64  `json:"owner"`
}

type CreateReq struct {
	Name string `json:"name"`
}

type ListReq struct {
	Limit int `json:"limit"`
}

type ListResp struct {
	Rooms []*RoomResp `json:"rooms"`
}
