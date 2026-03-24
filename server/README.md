# server

Go 后端（Gin + WebSocket）实现。当前提供：

- 创建房间 `/api/rooms`
- 加入房间 `/api/rooms/:roomCode/join`
- 开始游戏 `/api/rooms/:roomCode/start`
- 房间实时快照 `/ws/:roomCode?playerId=...`

运行：

```bash
go run ./cmd/server
```
