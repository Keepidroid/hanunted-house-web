# 后端架构设计

> 拆分来源：docs/technical-architecture.md
>
> 说明：本文档是主技术架构文档的拆分子文档，便于后续按阶段引用。

## 3. 后端架构

### 3.1 模块划分

建议按如下 Go 包结构拆分：

- `cmd/server`
- `internal/httpapi`
- `internal/ws`
- `internal/room`
- `internal/game`
- `internal/game/commands`
- `internal/game/events`
- `internal/game/view`
- `internal/content`
- `internal/store`
- `internal/authsession`
- `internal/telemetry`

核心模块职责如下。

#### `room`

- 管理房间生命周期
- 维护 `RoomRuntime` 注册表
- 按需加载房间运行时
- 清理空闲房间

#### `game`

- 权威 `GameState`
- 回合流程控制
- 地图、卡牌、属性、战斗、怪物、作祟引擎
- 命令校验与事件应用

#### `view`

- 从权威状态生成玩家视角
- 处理 public / team secret / player secret / hidden identity

#### `content`

- 加载角色、卡牌、房间、作祟配置
- 提供内容查询接口
- 校验内容引用完整性

#### `store`

- 房间、会话、快照、事件日志的数据库读写

### 3.2 房间管理

### 3.2.1 房间基础模型

每个房间包含：

- 唯一 `roomId`
- 短 `roomCode`
- 房主 seat
- 3-6 个座位
- 房间状态：`lobby / running / finished / archived`
- 活动 `gameId`

房间不需要复杂匹配系统，直接以房间码组织即可。

### 3.2.2 创建与加入

推荐：

- 房主创建房间时自动占用 `seat0`
- 加入房间时按空位分配 seat
- 房间链接中只暴露 `roomCode`
- 真正的玩家身份凭据使用后端签发的 `playerToken`

这样分离有两个好处：

- 房间码适合分享
- 玩家令牌适合重连，不应直接暴露为可猜测标识

### 3.2.3 重连

首版建议采用**游客优先**方案，不强依赖账号系统。

实现：

- 玩家首次加入房间时，服务端签发 `playerToken`
- 前端保存在 `localStorage`
- 重连时用该令牌恢复到同一 seat

好处：

- 首版无需登录系统
- 对桌游房间分享最自然
- 足够满足断线重连需求

### 3.2.4 断线处理

推荐策略：

- WS 断开后，seat 不立即释放
- 标记该玩家为 `offline`
- 游戏继续保留其位置和私密信息
- 若在 lobby 阶段掉线，房主可选择踢出
- 若游戏进行中掉线，只允许原 seat 的令牌重连

服务端还应维护：

- 心跳超时检测
- 断线时间戳
- 房间空置超时回收策略

建议：

- `lobby` 空房 30 分钟回收
- `running` 房间保留至少 12-24 小时，防止长局中断

### 3.2.5 RoomRuntime：一房间一个串行执行器

这是后端最关键的设计。

每个房间对应一个 `RoomRuntime`：

- 内部维护权威状态
- 只有一个输入命令邮箱
- 同时只有一个 goroutine 消费该邮箱

优点：

- 同一房间天然无并发写冲突
- 不需要为游戏状态加复杂锁
- 非常适合回合制桌游

推荐结构：

```go
type RoomRuntime struct {
    RoomID        string
    State         *GameState
    Connections   map[PlayerID]*ClientConn
    Mailbox       chan RoomCommand
    Persist       PersistenceWriter
    Projector     ViewProjector
}
```

### 3.3 游戏状态管理

### 3.3.1 建模原则

一局游戏的状态必须满足：

- **完整**：足以恢复一局
- **权威**：客户端不拥有独立真实状态
- **可序列化**：可存快照
- **可过滤**：可导出不同玩家视角
- **可回放**：可由事件推进

### 3.3.2 推荐的 `GameState` 结构

```ts
type GameState = {
  meta: {
    roomId: string
    gameId: string
    contentVersion: string
    stateVersion: number
    phase: "lobby" | "setup" | "preHaunt" | "hauntReveal" | "postHaunt" | "finished"
  }
  seats: SeatState[]
  turn: TurnState
  board: BoardState
  decks: DeckState
  explorers: Record<ExplorerId, ExplorerState>
  corpses: Record<ExplorerId, CorpseState>
  haunt: HauntState | null
  monsters: Record<MonsterInstanceId, MonsterState>
  prompts: PromptState[]
  modifiers: OngoingEffectState[]
  logCursor: number
}
```

其中关键子结构应包括：

#### `TurnState`

- 当前轮到谁
- 当前阶段
- 当前回合可用移动值
- 本回合已交易次数
- 本回合是否已攻击
- 本回合已使用过的特殊行动
- 是否处于等待选择目标/分配伤害/确认等子流程

#### `BoardState`

- 已放置板块列表
- 每块板块的区域、坐标、朝向、占用 footprint
- room node 邻接图
- 未探索门口集合
- 特殊连接集合
- 当前 token 的地图位置

#### `DeckState`

- Event / Item / Omen 当前顺序
- 已持有牌归属
- buried 结果反映在牌序变化上
- 已选择场景卡

#### `ExplorerState`

- 玩家归属
- 角色 ID
- 当前所在空间
- 四项属性格位与显示值
- 持有 Item/Omen
- 是否为奸徒/英雄/未知身份
- 是否死亡

#### `HauntState`

- 是否已触发
- hauntId
- hauntType
- 触发预兆
- haunt revealer
- 阵营分配
- 当前公开规则
- 每方已解锁信息
- 作祟脚本的局部状态

### 3.3.3 Canonical State 与 Player View 分离

不要在服务端维护多份“按玩家定制后的真实状态”。

推荐：

- 服务端只维护一份 `canonical state`
- 发送前通过 `ViewProjector` 投影成 `PlayerView`

好处：

- 减少状态分叉
- 规则逻辑只写一份
- 更容易排查隐藏信息错误

### 3.3.4 操作日志 + 状态快照

推荐采用混合方案：

- 每个已提交命令都会生成 `event log`
- 按策略写入 `snapshot`

推荐快照策略：

- 游戏开始时
- 每回合结束时
- 作祟触发时
- 游戏结束时
- 或每 N 个事件一次

这比“每个命令都存整局快照”更节省，也比“只有日志没有快照”更容易恢复。

### 3.3.5 随机性处理

所有随机结果都必须由服务端决定，并进入事件日志。

例如：

- `DiceRolled`
- `TileDrawn`
- `CardDrawn`
- `HauntTriggered`

不要依赖“重放时重跑随机数”。

正确做法是：

- 日志直接记录随机结果
- 重放时应用这些结果

这样回放与恢复才能绝对一致。

### 3.4 实时通信方案

### 3.4.1 REST 与 WebSocket 的职责划分

推荐分工：

#### REST 负责

- 创建房间
- 加入房间
- 恢复会话
- 获取初始房间/游戏快照
- 健康检查

#### WebSocket 负责

- 房间内所有实时操作
- 房间状态广播
- 在线状态与重连
- 回合推进通知
- 抽牌/掷骰/作祟/死亡等事件推送

原因：

- REST 适合无状态入口操作
- WS 适合长连接房间同步
- 两者职责清晰，调试简单

### 3.4.2 推荐消息协议

#### Client -> Server

```json
{
  "type": "command",
  "commandId": "cmd_123",
  "baseVersion": 42,
  "payload": {
    "kind": "moveExplorer",
    "args": { "path": ["room_a", "room_b"] }
  }
}
```

#### Server -> Client

- `commandAccepted`
- `commandRejected`
- `eventBatch`
- `viewSync`
- `presenceChanged`
- `toast`

### 3.4.3 同步策略：事件流 + 视图同步

推荐采用：

- **事件用于表现和动画**
- **视图同步用于纠偏和恢复**

也就是每次命令处理后：

1. 服务端广播公开事件批次
2. 向每个玩家发送过滤后的 `viewSync`

首版可以直接发送**完整过滤视图**，而不是复杂 diff。

原因：

- 一局游戏状态规模很小
- 玩家数量少
- 规则正确性比节省几十 KB 更重要
- 完整视图更利于重连、排错和版本升级

因此，推荐：

- 加入/重连：发送完整 `PlayerView`
- 每次提交后：发送事件批次 + 精简版完整视图

### 3.4.4 信息隔离

这是本项目最关键的实时通信问题之一。

推荐设计三层可见性：

- `public`
- `team`
- `player`

例如：

- 角色属性、地图、公共日志：`public`
- 某方已读到的作祟规则：`team`
- 某玩家可见的手牌选择或待确认提示：`player`

具体实现上，不建议在前端依赖“收到后自己隐藏”。

必须做到：

- **服务端根本不把不该看的数据发给你**

### 3.5 游戏逻辑引擎

### 3.5.1 命令驱动模型

推荐采用：

- `Command`：玩家意图
- `Validator`：合法性校验
- `Event`：真正发生的事
- `Apply`：把事件应用到状态

流程：

1. 接收命令
2. 校验命令
3. 生成事件
4. 应用事件
5. 持久化事件
6. 投影视图

这个模型非常适合《山屋惊魂》，因为它是：

- 回合制
- 强规则
- 高审计需求
- 经常有随机结算
- 经常有连锁效果

### 3.5.2 回合流程控制

后端应以显式状态机维护游戏阶段。

推荐阶段：

- `Lobby`
- `CharacterSelect`
- `ScenarioSelect`
- `GameSetup`
- `PreHauntTurn`
- `HauntReveal`
- `PostHauntTurn`
- `MonsterTurn`
- `EndGame`

注意：

- `MonsterTurn` 应作为显式阶段，而不是普通副作用
- `Prompt/Choice` 子状态应独立表示

例如：

- 等待玩家分配伤害
- 等待选择交易对象
- 等待执行房间放置
- 等待作祟 setup 完成

### 3.5.3 规则校验

服务端必须校验：

- 当前是不是该玩家回合
- 当前阶段是否允许该动作
- 玩家是否满足资源条件
- 路径是否合法
- 房间放置是否合法
- 攻击目标是否合法
- 信息可见性是否允许该选择

前端高亮只用于 UX，不能代替后端校验。

### 3.5.4 作祟与内容引擎

不建议把 50 个作祟全部硬编码在控制器里。

推荐采用**数据优先 + Go Hook 兜底**模式。

#### 通用规则层负责

- 移动
- 探索
- 抽牌
- 属性变化
- 攻击
- 视线
- obstacle
- 怪物基础规则

#### 内容定义层负责

- 角色数据
- 房间数据
- 卡牌数据
- 作祟元数据
- 作祟目标、特殊动作、token 需求、重要地点

#### Hook 兜底层负责

当某个作祟超出通用 DSL 能力时，提供：

- `hauntId -> Go handler`

这比“一开始造一个通用脚本语言”更现实。

### 3.6 API 设计思路

#### REST 示例

- `POST /api/rooms`
- `POST /api/rooms/{roomCode}/join`
- `POST /api/rooms/{roomCode}/reconnect`
- `GET /api/rooms/{roomCode}`
- `GET /api/rooms/{roomCode}/snapshot`
- `GET /api/healthz`

#### WebSocket 示例

- `GET /ws/rooms/{roomCode}?playerToken=...`

WebSocket 命令类型可包括：

- `setReady`
- `selectCharacter`
- `selectScenario`
- `startGame`
- `moveExplorer`
- `discoverRoom`
- `tradeCards`
- `useSpecialAction`
- `attack`
- `assignDamage`
- `resolvePrompt`
- `endTurn`
