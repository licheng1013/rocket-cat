package room

import "errors"

// QueryService 查询单个房间信息。
func QueryService(req *QueryReq) (*RoomResp, error) {
	if req == nil || req.RoomId <= 0 {
		return nil, errors.New("roomId required")
	}

	return &RoomResp{
		RoomId: req.RoomId,
		Name:   "Lobby",
		Owner:  10001,
	}, nil
}

// CreateService 创建房间。
func CreateService(req *CreateReq, owner int64) (*RoomResp, error) {
	if req == nil || req.Name == "" {
		return nil, errors.New("name required")
	}

	return &RoomResp{
		RoomId: 1,
		Name:   req.Name,
		Owner:  owner,
	}, nil
}

// ListService 查询房间列表。
func ListService(req *ListReq) (*ListResp, error) {
	limit := 20
	if req != nil && req.Limit > 0 {
		limit = req.Limit
	}
	if limit > 100 {
		limit = 100
	}

	return &ListResp{
		Rooms: []*RoomResp{
			{RoomId: 1, Name: "Lobby", Owner: 10001},
		},
	}, nil
}
