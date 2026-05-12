# Rocket Cat

Rocket Cat is a small Go WebSocket action framework built around explicit routes:

```text
route = (cmd << 16) | subCmd
```

The project favors simple module boundaries, explicit registration, unified request binding, and unified responses.

## Layout

```text
/ws
    bind.go
    context.go
    packet.go
    response.go
    router.go
    server.go
    session.go

/modules
    /user
        cmd.go
        model.go
        handler.go
        service.go
    /chat
        cmd.go
        model.go
        handler.go
        service.go
    /room
        cmd.go
        model.go
        handler.go
        service.go

/main.go
```

## Run

```bash
go run .
```

The WebSocket endpoint is:

```text
ws://localhost:10100/ws
```

## Packet

Requests use JSON:

```json
{
  "cmd": 1001,
  "subCmd": 1,
  "data": {
    "account": "demo",
    "password": "123456"
  }
}
```

Responses use:

```json
{
  "cmd": 1001,
  "subCmd": 1,
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

## Module Rules

Each module must contain:

```text
cmd.go      command and route definitions only
model.go    request and response structs only
handler.go  bind request, call service, return response
service.go  core business logic only
```

Handlers should use `ws.Bind[T](ctx)` instead of calling `json.Unmarshal` directly.

Modules register routes explicitly:

```go
func Register(router *ws.Router) {
	router.Register(Cmd, Login, LoginHandler)
}
```

## Example Handler

```go
func LoginHandler(ctx *ws.Context) {
	req, err := ws.Bind[LoginReq](ctx)
	if err != nil {
		return
	}

	resp, err := LoginService(req)
	if err != nil {
		ws.Fail(ctx, 1, err.Error())
		return
	}

	ws.OK(ctx, resp)
}
```

## Tests

```bash
go test ./...
```
