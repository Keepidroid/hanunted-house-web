# server

Go 后端（Gin + WebSocket）实现。当前提供：

- 创建房间 `/api/rooms`
- 加入房间 `/api/rooms/:roomCode/join`
- 开始游戏 `/api/rooms/:roomCode/start`
- 测试指令：伤害 `/api/rooms/:roomCode/commands/damage`
- 测试指令：掷骰 `/api/rooms/:roomCode/commands/roll`
- 房间实时快照 `/ws/:roomCode?playerId=...`

当前已按机制文档推进：

- 初始探索者与四维属性轨（Might/Speed/Knowledge/Sanity）
- gain/lose/heal 与 critical 约束
- trait/attack/haunt/other 掷骰分类基础

运行：

```bash
go run ./cmd/server
```
