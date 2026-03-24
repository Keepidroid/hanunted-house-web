# 核心技术难点分析

> 拆分来源：docs/technical-architecture.md
>
> 说明：本文档是主技术架构文档的拆分子文档，便于后续按阶段引用。

## 5. 核心技术难点分析

### 5.1 房间板块的动态地图生成

难点：

- 三层区域独立增长
- 有门连接与特殊连接
- 有三空间起始板块
- 有“必须保留可继续扩张空间”的放置限制
- 极端情况下需要最小重排

解决思路：

- 服务端维护 `BoardState + RoomNode adjacency graph`
- 内容层为每种 tile 定义：
  - footprint
  - door topology
  - region
  - special links
- 写一个专门的 `PlacementValidator` 和 `PlacementSolver`
- 客户端只消费结果，不自己做合法性判断

### 5.2 作祟触发后的游戏分叉

难点：

- 同一局游戏会从合作阶段切换到对抗/隐匿/自由混战
- 阵营、目标、提示、规则、可见信息都变
- 还要支持 setup 前后、monster turn、特殊死亡等状态分支

解决思路：

- 显式相位状态机
- 作祟内容抽象为 `HauntDefinition + HauntRuntimeState`
- 通用引擎负责阶段切换
- 具体作祟只覆写必要规则

### 5.3 信息不对称的实现

难点：

- 作祟后双方能看到的内容不同
- Hidden Traitor 还要求身份延迟公开
- 规则文本是“使用前可保密，使用后可要求公开”

解决思路：

- 服务端只维护 canonical state
- 用 `ViewProjector` 生成不同 `PlayerView`
- 事件也带 visibility scope
- 永远不把不该看的字段发到客户端

### 5.4 50 个作祟的可配置化

难点：

- 每个作祟都可能引入不同 token 语义、特殊动作、胜利条件、怪物规则
- 纯硬编码会让后期极其痛苦
- 纯 DSL 又很容易做成半套脚本语言，复杂度失控

解决思路：

- 通用规则数据化
- 作祟内容尽量用声明式配置描述
- 提供少量 Go Hook 作为兜底扩展点
- 用自动化测试覆盖每个作祟的 setup、目标、关键触发器

### 5.5 断线重连与服务重启恢复

难点：

- 游戏存在隐藏信息，不能简单把全量状态广播回来
- 玩家刷新页面后必须回到自己的正确视角
- 服务器重启后不能丢局

解决思路：

- 玩家 seat 绑定 `playerToken`
- PostgreSQL 保存 `snapshot + event log`
- 恢复时重建 canonical state，再重新投影玩家视图

### 5.6 动画表现与权威状态一致性

难点：

- 前端想做更流畅的抽牌、掷骰、攻击动画
- 但游戏不能因为动画导致客户端与服务端状态不一致

解决思路：

- 前端只做“表现层动画”
- 动画数据来自服务端事件
- 真正状态更新以 `viewSync` 为准
- 不做强 optimistic update
