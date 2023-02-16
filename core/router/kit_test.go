package router

import (
	"testing"
)

func TestKit(t *testing.T) {
	cmd := 1
	subCmd := 1
	merge := CmdKit.GetMerge(int64(cmd), int64(subCmd))
	if CmdKit.GetCmd(merge) != int64(cmd) {
		t.Error("Cmd Error")
	}
	if CmdKit.GetSubCmd(merge) != int64(subCmd) {
		t.Error("SubCmd Error")
	}
}
