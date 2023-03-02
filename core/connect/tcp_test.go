package connect

import (
	"fmt"
	"testing"
)

func TestMyProtocol(t *testing.T) {
	m := MyProtocol{}
	m.SetData([]byte("HelloWorld"))
	encode := Encode(&m)
	result := Decode(encode)
	fmt.Println(string(result.Data))
}
