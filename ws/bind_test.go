package ws

import (
	"encoding/json"
	"testing"
)

type bindReq struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestBind(t *testing.T) {
	data, err := json.Marshal(bindReq{Name: "rocket", Age: 1})
	if err != nil {
		t.Fatal(err)
	}

	req, err := Bind[bindReq](&Context{
		Packet: &Packet{Data: data},
	})
	if err != nil {
		t.Fatal(err)
	}

	if req.Name != "rocket" || req.Age != 1 {
		t.Fatalf("unexpected req: %+v", req)
	}
}

func TestBindEmptyPacket(t *testing.T) {
	_, err := Bind[bindReq](nil)
	if err != ErrEmptyPacket {
		t.Fatalf("expected ErrEmptyPacket, got %v", err)
	}
}
