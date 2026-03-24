# 前端架构设计

> 拆分来源：docs/technical-architecture.md
>
> 说明：本文档是主技术架构文档的拆分子文档，便于后续按阶段引用。

## 4. 前端架构

### 4.1 页面结构与路由设计

推荐保持路由简单。

- `/`
  - 首页，输入昵称、创建房间、加入房间
- `/room/:roomCode`
  - 房间主页面
  - 根据房间状态渲染 Lobby / CharacterSelect / Game / EndGame
- `/rejoin`
  - 可选，仅用于读取本地玩家令牌并跳转回房间

不建议把 lobby、game、post-game 拆成太多独立路由，因为它们共享同一个房间上下文和连接。

### 4.2 页面层级

推荐页面骨架：

- `AppShell`
- `RoomPage`
- `LobbyView`
- `GameView`
- `EndGameView`

其中 `GameView` 再拆为：

- `BoardViewport`
- `FloorBoard`
- `PlayerSidebar`
- `HandPanel`
- `ActionPanel`
- `GameLogPanel`
- `ModalHost`
- `ConnectionStatusBar`

### 4.3 核心组件拆分

#### 地图区域

- `BoardViewport`
  - 负责缩放、拖拽、重置视角
- `FloorBoard`
  - 渲染单个楼层
- `TilePlacement`
  - 渲染房间板块
- `DoorOverlay`
  - 渲染可用门口和特殊连接
- `TokenLayer`
  - 渲染 token
- `ExplorerLayer`
  - 渲染探索者和怪物

#### 玩家信息区

- `PlayerList`
- `CharacterPanel`
- `TraitTrack`
- `InventoryPanel`
- `TeamStatusPanel`

#### 回合与动作区

- `TurnBanner`
- `ActionBar`
- `PromptPanel`
- `DiceRollModal`
- `CardModal`
- `TradeModal`
- `AttackResolutionModal`
- `HauntRevealModal`

#### 其他

- `GameLogPanel`
- `ConnectionBadge`
- `ToastHost`
- `RulesExcerptPanel`

### 4.4 前端状态管理方案

推荐把前端状态分成 4 类。

#### 1. 会话状态

用 Zustand 保存：

- 当前昵称
- 当前 `playerToken`
- 当前房间码
- 当前连接状态

#### 2. 服务端同步状态

用 Zustand 保存当前 `PlayerView`。

这是前端真正的“游戏数据源”：

- 房间状态
- 游戏阶段
- 玩家可见地图
- 玩家可见手牌
- 当前 prompt
- 当前日志

#### 3. REST 生命周期状态

用 TanStack Query 处理：

- 创建房间请求
- 加入房间请求
- 恢复会话请求
- 初始快照加载

#### 4. 局部 UI 状态

用 Zustand 或组件局部 state 保存：

- 当前选中的房间
- 当前是否打开卡牌弹窗
- 当前地图缩放值
- 当前底部抽屉是否展开

原则：

- **前端不缓存“独立推演出来的游戏真相”**
- 一切以服务端推送的 `PlayerView` 为准

### 4.5 地图渲染方案

这是前端里最重要的结构设计。

### 4.5.1 推荐方案：DOM 绝对定位 + 坐标驱动渲染

推荐把地图建成一个逻辑坐标平面：

- 每个已放置 tile 都有：
  - `floor`
  - `anchorX`
  - `anchorY`
  - `orientation`
  - `footprint`

前端只负责：

- 按服务端给出的坐标与朝向渲染
- 不负责自己计算“是否可放置”

推荐采用 DOM 绝对定位，而不是 Canvas。

原因：

- 房间数量不会大到需要 Canvas 性能
- DOM 更容易做点击、hover、tooltip、可访问性
- 手牌、人物、token、可选动作高亮都更适合 DOM 叠层
- CSS transform 做缩放和平移更简单

### 4.5.2 为什么不推荐 Canvas 作为首版主方案

Canvas 的优点是性能，但本项目更缺的不是性能，而是：

- 命中测试
- 状态叠层
- 可访问性
- 调试便利性

在只有几十块房间板块、少量角色 token 的情况下，Canvas 收益很小。

### 4.5.3 三层地图展示策略

地图有 Basement / Ground Floor / Upper Floor 三层。

推荐 UI 策略：

- 桌面端：支持同时展示 3 层，或主视图 + 其余楼层缩略图
- 移动端：默认单层显示，通过 floor tabs 切换

但底层渲染模型应统一：

- 每层都是一个 `FloorBoard`
- 每层共享同一种坐标体系

这样桌面与移动端只是布局差异，不影响数据结构。

### 4.5.4 多空间起始板块

由于起始 Ground Floor 板块是三空间大板块，地图数据不能假设“一块 tile = 一个 room”。

推荐：

- 逻辑上按 `RoomNode` 建图
- 渲染上按 `TilePlacement` 建板块
- 由 `TilePlacement` 指向一个或多个 `RoomNode`

这样既能正确处理移动与邻接，也能正确显示大板块。

### 4.6 移动端适配策略

移动端的重点不是“缩小桌面布局”，而是重新组织交互。

推荐：

- 地图占据主视口
- 玩家信息区和手牌区放到底部抽屉
- 当前回合动作放在底部固定 Action Bar
- 弹窗改为全屏 bottom sheet 或全屏 overlay
- 增大点击热区
- 支持双指缩放、单指拖动、双击聚焦
- 使用 CSS `env(safe-area-inset-*)` 处理安全区

PWA 策略：

- 缓存静态资源和内容资源
- 不提供离线对战
- 断线时给出明确提示和自动重连
