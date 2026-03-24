# 核心机制总览与状态边界

> 拆分来源：[core-mechanics.md](../core-mechanics.md)
>
> 说明：本文档聚焦整体阶段结构、核心状态实体与开局流程，适合作为所有机制实现的共同入口。

## 1. 游戏阶段总览

《山屋惊魂》第三版可以拆成两个主阶段：

1. `Pre-Haunt`
   - 所有玩家合作探索房屋
   - 拼出地图、翻房间、抽牌、强化角色
2. `Post-Haunt`
   - 根据场景卡和触发预兆牌确定具体作祟
   - 阵营、目标、特殊规则与额外单位开始生效

从实现角度，这不是“固定地图上的战棋”，而是四层组合：

- 会持续生长的三层房屋图结构
- 一组带公开属性轨的探索者单位
- 由房间、卡牌、作祟共同改写的规则引擎
- 在作祟后强力覆写行为的内容系统

## 2. 核心状态实体

### 2.1 玩家、探索者与阵营

- `Player`
  - 真实玩家
- `Explorer`
  - 桌面上的基础可行动单位
- `Side / Team`
  - 作祟后生效的阵营层
  - 常见值包括 `heroes`、`traitor`、`no_traitor`、`hidden_traitor`、`free_for_all`

实现重点：

- 作祟后英雄和奸徒仍然都是 `Explorer`
- 阵营变化属于规则层与权限层变化，不是单位类型变化

### 2.2 角色属性轨

每个角色有四项属性：

- `Might`
- `Speed`
- `Knowledge`
- `Sanity`

建议角色状态至少包含：

- `traitTrackIndex`
- `traitValue`
- `startIndex`
- `startValue`

原因：

- 规则按“移动几格”结算，不是简单按数值加减
- 角色板存在重复数字，移动一格不一定改变显示值
- 治疗返回的是起始绿色位置，本质是重置到初始格位

### 2.3 房屋、板块与空间

房屋由三个区域组成：

- `Basement`
- `Ground Floor`
- `Upper Floor`

建议把地图拆成两个层次：

- `PhysicalTile`
  - 物理板块
- `SpaceNode`
  - 可站立、可计算邻接、可进入战斗与视线判定的逻辑空间

这样建模的原因是：

- 大多数房间板块只有一个空间
- Ground Floor 起始板块实际上包含 `Entrance Hall`、`Hallway`、`Ground Floor Staircase` 三个空间
- `Upper Landing` 与 `Basement Landing` 也是独立起始板块

### 2.4 牌堆与拥有权

核心牌堆包括：

- `Event Deck`
- `Item Deck`
- `Omen Deck`
- `Scenario Deck`

牌的核心状态建议区分：

- 在哪个牌堆中
- 被哪个 `Explorer` 持有
- 是否被埋回牌堆底部

第三版的重要约束：

- Event 结算后不是进入弃牌堆，而是 `bury` 到牌堆底部
- `bury` 是通用术语，后续也会用于板块堆处理

### 2.5 作祟内容模块

作祟不应视为若干散落条件，而应作为独立内容模块建模，至少包括：

- `hauntId`
- `scenario`
- `triggerOmen`
- `traitorSelectionRule`
- `hauntType`
- `heroIntro`
- `heroSetup`
- `heroSpecialRules`
- `heroSpecialActions`
- `heroObjective`
- `traitorIntro`
- `traitorSetup`
- `traitorSpecialRules`
- `traitorSpecialActions`
- `traitorObjective`
- `monsterBoxes`
- `tokensNeeded`
- `importantLocations`
- `playerCountVariants`

## 3. 开局设置流程

建议系统支持以下顺序：

1. 玩家选择角色
2. 发放角色板、模型、底座
3. 四个属性夹子放到绿色起始位置
4. 准备 8 颗骰子
5. 分别洗混 Event / Item / Omen
6. 摆放怪物与奸徒参考卡
7. 每位玩家获得玩家参考卡
8. 摆好三个起始板块
9. 其余房间板块洗混形成暗置板块堆
10. 所有探索者放在 `Entrance Hall`
11. 从 5 张场景卡中选 1 张作为本局场景
12. Number Track 与各类 token 放旁边待命
13. 由“生日最接近下一个真实生日”的玩家先手，之后顺时针

实现注意：

- 第三版不再使用旧版 `haunt matrix`
- 数字版若不读取真实生日，需要额外定义一个可操作的等价先手方案

## 4. 实现边界建议

后续开发中，建议把本文档当成“总状态边界说明”而不是具体规则手册：

- 角色轨细节见 [02-character-traits-and-dice.md](./02-character-traits-and-dice.md)
- 地图与探索流程见 [03-turns-map-and-exploration.md](./03-turns-map-and-exploration.md)
- 牌堆与行动限制见 [04-card-and-action-system.md](./04-card-and-action-system.md)
- 作祟映射与权限见 [05-haunt-engine.md](./05-haunt-engine.md)
- 战斗、死亡与怪物见 [06-combat-death-and-monsters.md](./06-combat-death-and-monsters.md)
