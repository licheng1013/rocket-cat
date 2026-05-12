package ws

import "testing"

func TestRoute(t *testing.T) {
	route := Merge(1001, 2)

	if route != Route(65601538) {
		t.Fatalf("unexpected route: %d", route)
	}
	if Cmd(route) != 1001 {
		t.Fatalf("unexpected cmd: %d", Cmd(route))
	}
	if Sub(route) != 2 {
		t.Fatalf("unexpected subCmd: %d", Sub(route))
	}
}
