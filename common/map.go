package common

import "sync"

func MapToList[K int64, V any](m map[K]V) ([]K, []V) {
	kList := make([]K, 0)
	vList := make([]V, 0)
	for k, v := range m {
		kList = append(kList, k)
		vList = append(vList, v)
	}
	return kList, vList
}

func SyncMapToList(m *sync.Map) ([]any, []any) {
	kList := make([]any, 0)
	vList := make([]any, 0)
	m.Range(func(key, value any) bool {
		kList = append(kList, key)
		vList = append(vList, value)
		return true
	})
	return kList, vList
}
