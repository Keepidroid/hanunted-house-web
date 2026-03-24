# 核心机制子文档索引

> 目录目的：将 [core-mechanics.md](../core-mechanics.md) 按后续开发可引用的边界拆分，避免在具体实现阶段反复通读整篇总文档。

## 文档清单

- [01-overview-and-state.md](./01-overview-and-state.md)
  - 游戏阶段总览、核心状态实体、开局设置
  - 适合进入任何机制开发前统一领域认知

- [02-character-traits-and-dice.md](./02-character-traits-and-dice.md)
  - 属性轨、治疗、伤害、死亡阈值、骰子分类
  - 适合角色系统、伤害系统、掷骰引擎实现阶段

- [03-turns-map-and-exploration.md](./03-turns-map-and-exploration.md)
  - 作祟前回合、移动、探索、房间放置、地图连接
  - 适合地图系统、回合系统、探索流程实现阶段

- [04-card-and-action-system.md](./04-card-and-action-system.md)
  - Event / Item / Omen、埋牌、交易、特殊行动
  - 适合牌堆系统、拥有权、行动限制实现阶段

- [05-haunt-engine.md](./05-haunt-engine.md)
  - 场景卡与作祟映射、阵营切换、作祟启动、脚本抽象、信息权限
  - 适合作祟引擎、内容配置、权限过滤实现阶段

- [06-combat-death-and-monsters.md](./06-combat-death-and-monsters.md)
  - 攻击、武器、视线、障碍、死亡、尸体、怪物与 monster turn
  - 适合战斗系统、怪物系统、作祟后通用规则实现阶段

- [07-checklists-and-module-boundaries.md](./07-checklists-and-module-boundaries.md)
  - 易错规则清单、推荐模块拆分、实现结论
  - 适合做任务拆分、评审对照、联调验收阶段

## 推荐阅读顺序

建议默认顺序：

1. [01-overview-and-state.md](./01-overview-and-state.md)
2. [02-character-traits-and-dice.md](./02-character-traits-and-dice.md)
3. [03-turns-map-and-exploration.md](./03-turns-map-and-exploration.md)
4. [04-card-and-action-system.md](./04-card-and-action-system.md)
5. [05-haunt-engine.md](./05-haunt-engine.md)
6. [06-combat-death-and-monsters.md](./06-combat-death-and-monsters.md)
7. [07-checklists-and-module-boundaries.md](./07-checklists-and-module-boundaries.md)

## 按开发阶段引用建议

### 阶段 1：领域建模与基础状态骨架

优先阅读：

- [01-overview-and-state.md](./01-overview-and-state.md)
- [02-character-traits-and-dice.md](./02-character-traits-and-dice.md)
- [03-turns-map-and-exploration.md](./03-turns-map-and-exploration.md)

### 阶段 2：作祟前核心引擎

优先阅读：

- [02-character-traits-and-dice.md](./02-character-traits-and-dice.md)
- [03-turns-map-and-exploration.md](./03-turns-map-and-exploration.md)
- [04-card-and-action-system.md](./04-card-and-action-system.md)

### 阶段 3：作祟引擎与权限系统

优先阅读：

- [05-haunt-engine.md](./05-haunt-engine.md)
- [04-card-and-action-system.md](./04-card-and-action-system.md)
- [07-checklists-and-module-boundaries.md](./07-checklists-and-module-boundaries.md)

### 阶段 4：作祟后通用规则与怪物

优先阅读：

- [05-haunt-engine.md](./05-haunt-engine.md)
- [06-combat-death-and-monsters.md](./06-combat-death-and-monsters.md)
- [07-checklists-and-module-boundaries.md](./07-checklists-and-module-boundaries.md)

### 阶段 5：联调与规则验收

优先阅读：

- [07-checklists-and-module-boundaries.md](./07-checklists-and-module-boundaries.md)
- [03-turns-map-and-exploration.md](./03-turns-map-and-exploration.md)
- [05-haunt-engine.md](./05-haunt-engine.md)
- [06-combat-death-and-monsters.md](./06-combat-death-and-monsters.md)

## 说明

- 原始总文档仍保留在 [core-mechanics.md](../core-mechanics.md)，适合整体回顾与全局校对。
- 子文档以“边界可单独引用”为优先目标，因此内容有少量交叉引用，但不重复展开完整规则。
