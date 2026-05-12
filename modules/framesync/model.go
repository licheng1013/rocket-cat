package framesync

type CheckReq struct {
	ClientId string `json:"clientId"` // Client id.
}

type CheckResp struct {
	Exists  bool     `json:"exists"`            // Whether a sync room already exists.
	RoomId  string   `json:"roomId,omitempty"`  // Sync room id.
	Players []string `json:"players,omitempty"` // Room member ids.
	Frames  []*Frame `json:"frames,omitempty"`  // Frames from the beginning.
}

type JoinMatchReq struct {
	ClientId string `json:"clientId"` // Client id.
}

type JoinMatchResp struct {
	Matched bool     `json:"matched"`           // Whether a room has been matched.
	RoomId  string   `json:"roomId,omitempty"`  // Sync room id.
	Players []string `json:"players,omitempty"` // Room member ids.
}

type ExitMatchReq struct {
	ClientId string `json:"clientId"` // Client id.
}

type ExitMatchResp struct {
	Exited bool `json:"exited"` // Whether the matching queue was exited.
}

type SubmitReq struct {
	ClientId string         `json:"clientId"` // Client id.
	RoomId   string         `json:"roomId"`   // Sync room id.
	Input    map[string]any `json:"input"`    // Client input snapshot.
}

type SubmitResp struct {
	Frame *Frame `json:"frame"` // Generated frame.
}

type Frame struct {
	Index  int64                  `json:"index"`  // Frame index.
	Inputs map[string]InputRecord `json:"inputs"` // Inputs by client id.
}

type InputRecord struct {
	ClientId string         `json:"clientId"` // Client id.
	Input    map[string]any `json:"input"`    // Client input snapshot.
}

type PushEvent struct {
	Type    string   `json:"type"`              // Event type.
	RoomId  string   `json:"roomId,omitempty"`  // Sync room id.
	Players []string `json:"players,omitempty"` // Room member ids.
	Frame   *Frame   `json:"frame,omitempty"`   // Broadcast frame.
}
