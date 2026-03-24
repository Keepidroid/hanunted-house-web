# 属性轨、伤害与骰子系统

> 拆分来源：[core-mechanics.md](../core-mechanics.md)
>
> 说明：本文档聚焦角色属性轨、伤害分配、死亡阈值与掷骰分类，适合作为角色系统与掷骰引擎的直接参考。

## 1. 属性轨建模

第三版的属性处理不能简化成普通数值相加减。实现时建议把每项属性建成轨道状态，而不是单一整数。

每项属性至少保存：

- `trackIndex`
- `displayValue`
- `startIndex`
- `startValue`

原因：

- 规则按“移动几格”生效
- 轨道上可能存在重复数值
- 治疗的返回目标是起始绿色位置，不是“当前数值加到起始数值”

## 2. 属性变化规则

### 2.1 Gain / Lose

- `Gain X Trait`
  - 属性夹子上移 `X` 格
- `Lose X Trait`
  - 属性夹子下移 `X` 格

边界规则：

- 超过最高值时，停在顶端
- 作祟前不会死亡
- 如果继续下降会越过 critical，作祟前只停在 critical

### 2.2 Heal

- `Heal Trait`
  - 该属性回到起始绿色格位
- 如果当前已不低于起始位置，则不发生变化

这意味着治疗是“回基线”，不是“恢复固定点数”。

## 3. Critical 与死亡

实现中至少要明确三个区间：

- 正常区
- `critical`
- `skull`

规则定义：

- `critical`
  - 属性夹子位于最底部、紧邻骷髅的位置
- 死亡
  - 作祟开始后，若任一属性降到 `skull` 或以下，则探索者死亡
  - 除非该作祟明确覆写死亡条件

开发约束：

- 作祟前不得因为属性降到底而死亡
- 作祟后才进入真实死亡判定

## 4. 伤害类型与分配规则

伤害共有三类：

- `Physical damage`
  - 分配到 `Might` 和或 `Speed`
- `Mental damage`
  - 分配到 `Knowledge` 和或 `Sanity`
- `General damage`
  - 由受伤者分配到任意属性

通用分配原则：

- 伤害按“格数”分配，不按显示数值差
- 可在允许范围内自由拆分
- 作祟前，如果存在尚未到 critical 的合法属性，就不能把伤害压到已 critical 的属性上

这意味着底层接口应支持：

- 限定可分配属性集合
- 校验作祟前的 critical 限制
- 先做合法性校验，再落地格位变化

## 5. 骰子分类

建议在引擎层把掷骰分成四种，而不是统一套一个“掷骰”函数。

### 5.1 Trait Roll

- 按当前属性值掷对应数量的骰子
- 常见于房间、事件、卡牌效果

### 5.2 Attack Roll

- 本质仍是 Trait Roll
- 但语义是战斗比较，通常会与伤害计算联动

### 5.3 Haunt Roll

- 按全场当前持有的预兆牌总数掷骰
- 结果 `5+` 时触发作祟

### 5.4 Other Roll

- 文本直接要求“掷 N 颗骰子”
- 不视为 trait roll
- 不应吃“仅影响 trait roll”的修正

## 6. 实现建议

建议把“属性轨系统”独立成明确模块，并提供以下能力：

- 以轨道格位而不是数值做主状态
- 支持 gain / lose / heal 三类标准操作
- 暴露 critical 与 skull 阈值判定
- 为伤害分配提供合法属性范围校验
- 为骰子系统提供统一的属性读取接口

与回合、战斗、作祟的联动规则分别见：

- [03-turns-map-and-exploration.md](./03-turns-map-and-exploration.md)
- [05-haunt-engine.md](./05-haunt-engine.md)
- [06-combat-death-and-monsters.md](./06-combat-death-and-monsters.md)
