# 数据模型概要

> 拆分来源：docs/technical-architecture.md
>
> 说明：本文档是主技术架构文档的拆分子文档，便于后续按阶段引用。

## 6. 数据模型概要

## 6.1 持久化实体

### `rooms`

```ts
type RoomRow = {
  id: string
  code: string
  status: "lobby" | "running" | "finished" | "archived"
  hostSeatId: string
  currentGameId: string | null
  createdAt: string
  updatedAt: string
}
```

### `room_seats`

```ts
type RoomSeatRow = {
  id: string
  roomId: string
  seatIndex: number
  nickname: string
  playerTokenHash: string
  online: boolean
  joinedAt: string
  lastSeenAt: string
}
```

### `games`

```ts
type GameRow = {
  id: string
  roomId: string
  contentVersion: string
  phase: string
  stateVersion: number
  startedAt: string | null
  endedAt: string | null
}
```

### `game_snapshots`

```ts
type GameSnapshotRow = {
  id: string
  gameId: string
  stateVersion: number
  snapshotJson: unknown
  createdAt: string
}
```

### `game_events`

```ts
type GameEventRow = {
  id: string
  gameId: string
  seq: number
  actorSeatId: string | null
  visibility: "public" | "team" | "player"
  eventType: string
  payloadJson: unknown
  createdAt: string
}
```

## 6.2 运行时实体

### `RoomRuntime`

```ts
type RoomRuntime = {
  roomId: string
  gameState: GameState
  mailbox: RoomCommand[]
  onlineConnections: Record<string, ConnectionRef>
  dirtySinceSnapshot: boolean
}
```

### `BoardPlacement`

```ts
type BoardPlacement = {
  placementId: string
  tileId: string
  floor: "basement" | "ground" | "upper"
  anchor: { x: number; y: number }
  orientation: 0 | 90 | 180 | 270
  roomNodeIds: string[]
}
```

### `ExplorerState`

```ts
type ExplorerState = {
  explorerId: string
  seatId: string
  characterId: string
  roomNodeId: string
  traits: {
    might: TraitState
    speed: TraitState
    knowledge: TraitState
    sanity: TraitState
  }
  items: string[]
  omens: string[]
  side: "hero" | "traitor" | "neutral" | "hidden"
  dead: boolean
}
```

### `HauntState`

```ts
type HauntState = {
  started: boolean
  hauntId: string
  hauntType: "noTraitor" | "oneTraitor" | "hiddenTraitor" | "freeForAll"
  triggerOmenId: string
  hauntRevealerSeatId: string
  sides: Record<string, SideAssignment>
  runtimeFlags: Record<string, unknown>
  revealedTextKeys: string[]
}
```

### `MonsterState`

```ts
type MonsterState = {
  instanceId: string
  monsterDefId: string
  roomNodeId: string
  stunned: boolean
  alive: boolean
  runtimeFlags: Record<string, unknown>
}
```

## 6.3 静态内容模型

静态内容不建议进 PostgreSQL 主业务表。

推荐结构：

- `characters.json`
- `room_tiles.json`
- `cards_event.json`
- `cards_item.json`
- `cards_omen.json`
- `scenarios.json`
- `haunts/*.json`

这些内容应有稳定 ID，并由前后端共享类型定义。

## 6.4 持久化 vs 内存态

### 必须持久化

- 房间元数据
- 座位与玩家令牌
- 当前游戏元数据
- 操作日志
- 状态快照
- 内容版本号

### 只保存在内存

- 当前房间连接对象
- 房间命令队列
- 临时心跳状态
- 视图缓存
- 短期 debounce 写盘状态

### 可从持久化恢复，不必单独存

- 当前权威游戏状态
- 当前玩家视角
- 当前可见日志
