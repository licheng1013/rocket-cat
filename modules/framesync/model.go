package framesync

// CheckReq 表示检查已有同步房间请求。
type CheckReq struct {
	ClientId string `json:"clientId"` // Client id.
}

// CheckResp 表示检查已有同步房间响应。
type CheckResp struct {
	Exists  bool     `json:"exists"`            // Whether a sync room already exists.
	RoomId  string   `json:"roomId,omitempty"`  // Sync room id.
	Players []string `json:"players,omitempty"` // Room member ids.
	Frames  []*Frame `json:"frames,omitempty"`  // Frames from the beginning.
}

// JoinMatchReq 表示加入匹配请求。
type JoinMatchReq struct {
	ClientId string `json:"clientId"` // Client id.
}

// JoinMatchResp 表示加入匹配响应。
type JoinMatchResp struct {
	Matched bool     `json:"matched"`           // Whether a room has been matched.
	RoomId  string   `json:"roomId,omitempty"`  // Sync room id.
	Players []string `json:"players,omitempty"` // Room member ids.
}

// ExitMatchReq 表示退出匹配请求。
type ExitMatchReq struct {
	ClientId string `json:"clientId"` // Client id.
}

// ExitMatchResp 表示退出匹配响应。
type ExitMatchResp struct {
	Exited bool `json:"exited"` // Whether the matching queue was exited.
}

// SubmitReq 表示提交输入快照请求。
type SubmitReq struct {
	ClientId string         `json:"clientId"` // Client id.
	RoomId   string         `json:"roomId"`   // Sync room id.
	Input    map[string]any `json:"input"`    // Client input snapshot.
}

// SubmitResp 表示提交输入快照响应。
type SubmitResp struct {
	Frame *Frame `json:"frame"` // Generated frame.
}

// Frame 表示服务端生成的同步帧。
type Frame struct {
	Index  int64                  `json:"index"`  // Frame index.
	Inputs map[string]InputRecord `json:"inputs"` // Inputs by client id.
}

// InputRecord 表示某个客户端在一帧中的输入。
type InputRecord struct {
	ClientId string         `json:"clientId"` // Client id.
	Input    map[string]any `json:"input"`    // Client input snapshot.
}

// PushEvent 表示服务端主动推送事件。
type PushEvent struct {
	Type    string   `json:"type"`              // Event type.
	RoomId  string   `json:"roomId,omitempty"`  // Sync room id.
	Players []string `json:"players,omitempty"` // Room member ids.
	Frame   *Frame   `json:"frame,omitempty"`   // Broadcast frame.
}
