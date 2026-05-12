package framesync

import (
	"errors"
	"fmt"
	"sync"
)

var state = newSyncState()

type syncState struct {
	mu          sync.Mutex
	waiting     []string
	rooms       map[string]*syncRoom
	clientRooms map[string]string
}

type syncRoom struct {
	id      string
	players []string
	frames  []*Frame
}

// CheckService checks whether the client already has a sync room.
func CheckService(req *CheckReq) (*CheckResp, error) {
	if req == nil || req.ClientId == "" {
		return nil, errors.New("clientId required")
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	room := state.roomByClient(req.ClientId)
	if room == nil {
		return &CheckResp{Exists: false}, nil
	}

	return &CheckResp{
		Exists:  true,
		RoomId:  room.id,
		Players: append([]string(nil), room.players...),
		Frames:  cloneFrames(room.frames),
	}, nil
}

// JoinMatchService joins the client into the matching queue.
func JoinMatchService(req *JoinMatchReq) (*JoinMatchResp, *PushEvent, error) {
	if req == nil || req.ClientId == "" {
		return nil, nil, errors.New("clientId required")
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	if room := state.roomByClient(req.ClientId); room != nil {
		return nil, nil, errors.New("sync room already exists")
	}
	if state.isWaiting(req.ClientId) {
		return &JoinMatchResp{Matched: false}, nil, nil
	}

	state.waiting = append(state.waiting, req.ClientId)
	if len(state.waiting) < 2 {
		return &JoinMatchResp{Matched: false}, nil, nil
	}

	players := append([]string(nil), state.waiting[:2]...)
	state.waiting = state.waiting[2:]

	room := &syncRoom{
		id:      fmt.Sprintf("sync-%s-%s", players[0], players[1]),
		players: players,
		frames:  make([]*Frame, 0, 128),
	}
	state.rooms[room.id] = room
	for _, player := range players {
		state.clientRooms[player] = room.id
	}

	resp := &JoinMatchResp{
		Matched: true,
		RoomId:  room.id,
		Players: append([]string(nil), players...),
	}
	event := &PushEvent{
		Type:    "started",
		RoomId:  room.id,
		Players: append([]string(nil), players...),
	}
	return resp, event, nil
}

// ExitMatchService removes the client from the matching queue.
func ExitMatchService(req *ExitMatchReq) (*ExitMatchResp, error) {
	if req == nil || req.ClientId == "" {
		return nil, errors.New("clientId required")
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	for i, clientId := range state.waiting {
		if clientId == req.ClientId {
			state.waiting = append(state.waiting[:i], state.waiting[i+1:]...)
			return &ExitMatchResp{Exited: true}, nil
		}
	}

	return &ExitMatchResp{Exited: false}, nil
}

// SubmitService appends a client input and creates the next sync frame.
func SubmitService(req *SubmitReq) (*SubmitResp, *PushEvent, error) {
	if req == nil || req.ClientId == "" || req.RoomId == "" {
		return nil, nil, errors.New("clientId and roomId required")
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	room := state.rooms[req.RoomId]
	if room == nil {
		return nil, nil, errors.New("sync room not found")
	}
	if !contains(room.players, req.ClientId) {
		return nil, nil, errors.New("client is not in sync room")
	}

	frame := &Frame{
		Index: int64(len(room.frames) + 1),
		Inputs: map[string]InputRecord{
			req.ClientId: {
				ClientId: req.ClientId,
				Input:    cloneInput(req.Input),
			},
		},
	}
	room.frames = append(room.frames, frame)

	event := &PushEvent{
		Type:    "frame",
		RoomId:  room.id,
		Players: append([]string(nil), room.players...),
		Frame:   cloneFrame(frame),
	}
	return &SubmitResp{Frame: cloneFrame(frame)}, event, nil
}

// RoomPlayers returns the current member ids of a room.
func RoomPlayers(roomId string) []string {
	state.mu.Lock()
	defer state.mu.Unlock()

	room := state.rooms[roomId]
	if room == nil {
		return nil
	}
	return append([]string(nil), room.players...)
}

func newSyncState() *syncState {
	return &syncState{
		waiting:     make([]string, 0, 16),
		rooms:       make(map[string]*syncRoom),
		clientRooms: make(map[string]string),
	}
}

func (s *syncState) roomByClient(clientId string) *syncRoom {
	roomId := s.clientRooms[clientId]
	if roomId == "" {
		return nil
	}
	return s.rooms[roomId]
}

func (s *syncState) isWaiting(clientId string) bool {
	return contains(s.waiting, clientId)
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func cloneFrames(frames []*Frame) []*Frame {
	cloned := make([]*Frame, 0, len(frames))
	for _, frame := range frames {
		cloned = append(cloned, cloneFrame(frame))
	}
	return cloned
}

func cloneFrame(frame *Frame) *Frame {
	if frame == nil {
		return nil
	}

	inputs := make(map[string]InputRecord, len(frame.Inputs))
	for clientId, record := range frame.Inputs {
		inputs[clientId] = InputRecord{
			ClientId: record.ClientId,
			Input:    cloneInput(record.Input),
		}
	}
	return &Frame{
		Index:  frame.Index,
		Inputs: inputs,
	}
}

func cloneInput(input map[string]any) map[string]any {
	cloned := make(map[string]any, len(input))
	for key, value := range input {
		cloned[key] = value
	}
	return cloned
}
