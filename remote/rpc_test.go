package remote

import (
	"log"
	"testing"
)

func TestCallBody(t *testing.T) {
	b := &CallBody{Id: 30, Data: []byte("Hello")}
	data, _ := b.ToMarshal()

	c := &CallBody{}
	_ = c.ToUnmarshal(data)

	log.Println(c.Id, string(c.Data))
	if c.Id != 30 {
		t.Errorf("c.Id=%d", c.Id)
	}
}
