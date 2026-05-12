package ws

import "testing"

func TestRouterDispatch(t *testing.T) {
	router := NewRouter()
	called := false

	router.Register(1001, 1, func(ctx *Context) {
		called = true
	})
	router.Dispatch(&Context{
		Packet: &Packet{Cmd: 1001, SubCmd: 1},
	})

	if !called {
		t.Fatal("handler was not called")
	}
}

func TestRouterMiddlewareOrder(t *testing.T) {
	router := NewRouter()
	order := make([]int, 0, 3)

	router.Use(func(next Handler) Handler {
		return func(ctx *Context) {
			order = append(order, 1)
			next(ctx)
		}
	})
	router.Use(func(next Handler) Handler {
		return func(ctx *Context) {
			order = append(order, 2)
			next(ctx)
		}
	})
	router.Register(1001, 1, func(ctx *Context) {
		order = append(order, 3)
	})

	router.Dispatch(&Context{
		Packet: &Packet{Cmd: 1001, SubCmd: 1},
	})

	want := []int{1, 2, 3}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("unexpected middleware order: %+v", order)
		}
	}
}

func TestRouterDuplicateRoutePanic(t *testing.T) {
	router := NewRouter()
	router.Register(1001, 1, func(ctx *Context) {})

	defer func() {
		if recover() == nil {
			t.Fatal("expected duplicate route panic")
		}
	}()

	router.Register(1001, 1, func(ctx *Context) {})
}
