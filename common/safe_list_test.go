package common

import "testing"

func TestSafeList(t *testing.T) {
	// 测试
	safeList := &SafeList{}
	safeList.Add(1)
	safeList.Add(3)
	t.Log(safeList.GetList())
}
