# 易错规则清单与模块边界建议

> 拆分来源：[core-mechanics.md](../core-mechanics.md)
>
> 说明：本文档用于联调、评审与任务拆分阶段，对照最容易出错的规则，并给出推荐模块边界。

## 1. 最容易出错的规则清单

以下规则建议作为单独验收 checklist：

- 角色没有通用的“独特主动能力”系统，角色差异主要来自属性轨、起始值与卡牌或剧本效果
- 属性变化按“移动夹子几格”结算，不按显示数字直接加减
- 作祟前不会死亡，属性最低只能压到 critical
- Event 结算后埋回牌堆底部，不进入普通弃牌堆
- 抽到 Omen 后才进行 `haunt roll`
- 抽到最后一张 Omen 且尚未作祟时，会自动触发作祟
- 发现新房间后，本回合必须立即结束
- 新房间摆放时必须保证该区域仍可继续扩张，必要时需要最小重排
- 默认没有“攻击后偷窃”
- 作祟后仍可继续探索和发现新房间，只是不再进行 `haunt roll`
- 视线要求直线、同区域、不中途转向
- obstacle 的额外代价发生在离开当前空间时，不是进入目标空间时
- 敌对角色与怪物会构成 obstacle，被击晕怪物在翻面期间不再正常阻碍移动
- 交易每回合只有一次，且不能转出本回合已用于特殊行动或用于攻击的牌
- 只能使用你在本回合开始时就持有的 Item / Omen 特殊行动
- 作祟书内容不是绝对私密，规则一旦被使用，对方可要求公开相关全文
- 怪物需要独立的 `monster turn` 概念，且很多作祟会在该阶段前后挂触发器

## 2. 推荐模块拆分

若后续进入实现阶段，建议至少拆为以下模块：

- `character-system`
  - 属性轨、伤害、治疗、死亡阈值
- `map-system`
  - 区域、板块放置、邻接、视线、obstacle
- `turn-system`
  - 回合顺序、阶段切换、行动上限、回合结束
- `card-system`
  - 三类牌堆、抽牌、埋牌、拥有权、交易、特殊行动
- `combat-system`
  - 攻击、武器、伤害分配、目标合法性
- `haunt-engine`
  - 场景卡映射、作祟启动、阵营切换、脚本执行、权限过滤
- `monster-system`
  - monster turn、移动、stun / kill、monster box
- `content-data`
  - 角色、房间、卡牌、50 个作祟的数据定义

## 3. 推荐开发顺序

如果以后要按阶段开发，建议按依赖关系推进：

1. `character-system` + `map-system` + `turn-system`
2. `card-system`
3. `haunt-engine`
4. `combat-system` + `monster-system`
5. `content-data`

原因很直接：

- 属性轨与地图增长是强约束底层
- 牌堆与探索流程依赖前两者稳定
- 作祟引擎建立在前置规则之上
- 战斗与怪物主要发生在作祟后
- 内容数据应压在可运行骨架之后落地

## 4. 最终实现结论

第三版真正的实现难点不在普通回合流程，而在以下三点：

- 属性轨与地图增长是强约束底层规则，必须先做对
- 作祟系统是内容驱动的规则覆写层，必须数据化、脚本化
- 信息公开是半公开模型，必须从权限设计一开始就考虑

如果这三块底座建对，后续房间、卡牌、怪物与作祟内容才能稳定扩展。

## 5. 交叉引用

- 总状态入口见 [01-overview-and-state.md](./01-overview-and-state.md)
- 作祟前基础流程见 [03-turns-map-and-exploration.md](./03-turns-map-and-exploration.md)
- 作祟引擎与权限见 [05-haunt-engine.md](./05-haunt-engine.md)
- 作祟后通用规则见 [06-combat-death-and-monsters.md](./06-combat-death-and-monsters.md)
