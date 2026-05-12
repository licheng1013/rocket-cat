# Frame Sync Example Rules

Applies to:

```text
examples/frame-sync.html
modules/framesync
```

## Goal

The page simulates frame synchronization over the Rocket Cat WebSocket protocol.

It must support this flow:

1. Connect to the WebSocket server.
2. Input an ID.
3. Confirm the ID.
4. Send a check request before any matching action.
5. If a sync room already exists for the ID, reconnect to it and replay frames from the beginning.
6. If no sync room exists for the ID, allow joining or exiting matching.
7. When two clients are matched, create a sync room.
8. After a room starts, clients submit input snapshots and receive broadcast frames.

## Language Rules

- Do not use the word "游戏" on the page.
- Prefer "帧同步", "同步", "匹配", "房间", "帧流", and "客户端".

## Technical Rules

- Use a single HTML file.
- Use Vue from a CDN.
- Do not add npm dependencies.
- Do not add a build step.
- Keep HTML, CSS, and JavaScript inline.
- The page must be directly openable by a browser.

## Default Endpoint

```text
ws://localhost:10100/ws
```

## Commands

| Action | cmd | subCmd |
| --- | ---: | ---: |
| Check existing sync room | 5001 | 1 |
| Join matching | 5001 | 2 |
| Exit matching | 5001 | 3 |
| Submit input snapshot | 5001 | 4 |
| Server push | 5001 | 6 |

## Interaction Rules

- Matching buttons stay disabled until the ID has been confirmed.
- Confirming an ID must call `Check`.
- If `Check` returns `exists: true`, matching actions remain disabled.
- If `Check` returns `exists: true`, the page replays returned frames from the beginning.
- If `Check` returns `exists: false`, the page can join or exit matching.
- Once a room is started, the page submits input snapshots periodically.
- Incoming frame pushes update the latest frame index and frame stream.

## Backend Rules

- `cmd.go` contains command and route definitions only.
- `model.go` contains request and response structs only.
- `handler.go` binds params, calls service, returns or pushes results.
- `service.go` owns matching state, room state, and frame history.
- Service code must not operate `websocket.Conn` directly.
- Route registration must stay explicit.

## Maintenance Rules

When frame sync commands change, update:

- `modules/framesync/cmd.go`
- `examples/frame-sync.html`
- `examples/ws-client.html`
- `docs/ws-client-rules.md`
- this document
