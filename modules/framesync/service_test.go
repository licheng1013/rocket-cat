package framesync

import "testing"

func TestCheckServiceMissingRoom(t *testing.T) {
	state = newSyncState()

	resp, err := CheckService(&CheckReq{ClientId: "10001"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Exists {
		t.Fatal("expected missing room")
	}
}

func TestJoinMatchServiceCreatesRoom(t *testing.T) {
	state = newSyncState()

	first, event, err := JoinMatchService(&JoinMatchReq{ClientId: "10001"})
	if err != nil {
		t.Fatal(err)
	}
	if first.Matched || event != nil {
		t.Fatalf("expected first client waiting, got resp=%+v event=%+v", first, event)
	}

	second, event, err := JoinMatchService(&JoinMatchReq{ClientId: "10002"})
	if err != nil {
		t.Fatal(err)
	}
	if !second.Matched || second.RoomId == "" {
		t.Fatalf("expected matched room, got %+v", second)
	}
	if event == nil || event.Type != "started" {
		t.Fatalf("expected started event, got %+v", event)
	}
}

func TestSubmitServiceCreatesFrame(t *testing.T) {
	state = newSyncState()

	_, _, err := JoinMatchService(&JoinMatchReq{ClientId: "10001"})
	if err != nil {
		t.Fatal(err)
	}
	joined, _, err := JoinMatchService(&JoinMatchReq{ClientId: "10002"})
	if err != nil {
		t.Fatal(err)
	}

	resp, event, err := SubmitService(&SubmitReq{
		ClientId: "10001",
		RoomId:   joined.RoomId,
		Input:    map[string]any{"seq": float64(1)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Frame.Index != 1 {
		t.Fatalf("expected frame 1, got %d", resp.Frame.Index)
	}
	if event == nil || event.Type != "frame" {
		t.Fatalf("expected frame event, got %+v", event)
	}
}
