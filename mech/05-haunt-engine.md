# 作祟引擎与信息权限

> 拆分来源：[core-mechanics.md](../core-mechanics.md)
>
> 说明：本文档聚焦场景卡与作祟映射、作祟启动流程、阵营建模、内容脚本抽象与信息公开规则。

## 1. 场景卡与作祟映射

第三版的作祟由“场景卡 + 触发预兆牌”共同决定，不再使用旧版 `haunt matrix`。

场景卡至少承担三件事：

- 提供本局进入山屋的主题动机
- 将某个触发预兆映射到某个 `hauntId`
- 指定该作祟的奸徒选定规则

因此，作祟定位不是单一查表，而是：

1. 开局先锁定 `scenario`
2. 作祟触发时记录 `triggerOmen`
3. 依据两者组合找到 `hauntId`

## 2. 作祟类型

第三版共有四类作祟：

- `No Traitor`
  - 无奸徒，所有玩家同阵营，只读英雄书
- `One Traitor`
  - 一名奸徒对抗其余英雄
- `Hidden Traitor`
  - 作祟开始时奸徒身份不对全体公开，可互相攻击
- `Free-for-All`
  - 人人为敌，不共享胜利

这意味着系统不能把阵营硬编码成“英雄 vs 奸徒”二元模型。

## 3. 奸徒选定规则

奸徒选定常见规则包括：

- `Haunt Revealer`
- `Left of Haunt Revealer`
- 某属性最高者
- 某属性最低者
- 无奸徒
- 自由混战

平局规则：

- 若按某属性最高或最低选奸徒时出现平局，由距离 `haunt revealer` 最近的轮序玩家成为奸徒
- 若规则写“haunt revealer 左手边玩家是奸徒”，则 `haunt revealer` 自己明确排除在外

## 4. 作祟启动流程

定位到对应作祟后，按以下顺序执行：

1. 英雄方朗读 `introduction`
2. 英雄方执行 `setup`
3. 奸徒方朗读 `introduction`
4. 奸徒方执行 `setup`
5. 双方分开阅读各自剩余内容并制定策略

信息规则：

- `introduction` 与 `setup` 需要先公开朗读
- 之后未使用的作祟细则可以暂时保密
- 一旦某条规则或特殊行动被使用，对方可以要求完整朗读相关文本

这是“使用前可保密、使用后需公开说明”的半公开模型。

## 5. 作祟后的共享基础规则

作祟开始后，探索机制并不会消失。

英雄和奸徒通常仍保留：

- 移动
- 发现新房间
- 使用 Item / Omen
- 交易
- 执行特殊行动

变化点在于：

- 不再进行 `haunt roll`
- 可攻击对象由阵营关系决定
- 剧本会增加额外 special rules 与 special actions

## 6. 作祟后的首回合

规则书明确指出：所有作祟 `setup` 都发生在“作祟后的第一回合之前”。

因此实现上应支持：

- 某些作祟在 setup 中直接指定第一位行动者
- 若作祟未覆写，则默认沿用正常轮序继续

## 7. 作祟内容层抽象

为了支持完整复刻，底层应把作祟视为“数据 + 脚本”。

每个作祟至少应提供：

- `Identification`
  - 场景卡
  - 触发预兆
  - 奸徒判定
- `Introduction`
- `Setup`
- `Special Rules`
- `Special Actions`
- `Objective`
- `Tokens Needed`
- `Important Locations`
- `If You Win`
- `Monster Boxes`

## 8. 可变玩家数

部分作祟使用 `{a/b/c/d}` 形式表示不同玩家数下的参数，对应：

- 3 人局
- 4 人局
- 5 人局
- 6 人局

这意味着作祟脚本必须能按玩家数解析变量，而不是在内容层手写分支逻辑。

## 9. Number Track、Token 与指定房间

### 9.1 Token 与 Number Track

第三版中 token 语义大量由具体作祟定义，因此不应把所有 token 类型硬编码为固定规则对象。

更合理的抽象是：

- 底层只提供 token 容器、位置、朝向、数值槽、显示图标
- 某个作祟再解释该 token 代表什么

`Number Track` 同理，应视为通用状态容器，而不是固定用途的计数器。

### 9.2 查找指定房间

有些作祟依赖某个特定房间。

规则允许：

- 玩家查看板块堆，确认该房间位于哪个区域
- 查看后重新洗混板块堆

因此系统需要允许“受规则授权的未翻板块信息查看”。

## 10. 信息公开与权限模型

数字版至少要区分以下可见性级别：

- `Public`
  - 角色属性
  - 当前持有的 Item / Omen
  - 已翻开的地图
  - 当前 token / Number Track 状态
- `Team Secret`
  - 作祟书尚未用到的细则
  - 某方尚未宣告的计划
- `Hidden Identity`
  - 仅 hidden-traitor 作祟存在
- `Reveal on Use`
  - 某条作祟规则或特殊行动一旦被使用，对方可要求朗读原文

这意味着权限系统至少要支持：

- 按阵营分视图
- 按已使用与未使用动态解锁文本
- 在 hidden-traitor 作祟中维持身份未公开状态

## 11. 与其他边界的关系

- Omen 与 `haunt roll` 的牌堆侧规则见 [04-card-and-action-system.md](./04-card-and-action-system.md)
- 作祟后的攻击、死亡、怪物规则见 [06-combat-death-and-monsters.md](./06-combat-death-and-monsters.md)
