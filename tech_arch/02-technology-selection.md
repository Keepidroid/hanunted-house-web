# 技术选型与理由

> 拆分来源：docs/technical-architecture.md
>
> 说明：本文档是主技术架构文档的拆分子文档，便于后续按阶段引用。

## 2. 技术选型与理由

## 2.1 总体技术栈

### 前端

- `React 19 + TypeScript`
- `Vite`
- `React Router`
- `Zustand`
- `TanStack Query`
- `Tailwind CSS`
- `Radix UI`
- `react-zoom-pan-pinch`
- `vite-plugin-pwa`
- `Vitest + React Testing Library`
- `Playwright`

### 后端

- `Go 1.26+`
- `gin` — 高性能 HTTP 路由框架
- `huma` — REST API 框架，通过 `humagin` 适配器包装 gin，自动生成 OpenAPI 3.x 文档
- `gws` (`github.com/lxzan/gws`) — 高性能事件驱动 WebSocket 库
- `pgx/v5`
- `sqlc` — 从 SQL 生成类型安全的 Go 数据库访问代码
- `golang-migrate`
- `slog`
- `testify`

### 数据与内容

- `PostgreSQL 16`
- 静态内容文件：`JSON/YAML`
- `JSON Schema` 或启动时自校验

### 部署

- `docker compose`
- `Caddy`
- `PostgreSQL`

### 可观测性

- `Prometheus metrics`（首版建议保留接口，监控可后补）
- 结构化日志 `slog`

## 2.2 关键选型说明

### React + TypeScript

这是题设约束，同时也符合项目特点：

- UI 复杂，状态多，组件层次深
- 房间、角色、地图、卡牌、弹窗、动画都适合组件化拆分
- TypeScript 对复杂消息协议、视图状态和内容类型非常重要

### Vite

理由：

- 启动和热更新快
- 对 React/TS/PWA 支持成熟
- 首版开发效率高于更重的 SSR 方案

不推荐一开始上 Next.js，原因是本项目核心不是 SEO 或 SSR，而是实时房间交互。

### Zustand

推荐用于前端实时状态，而不是 Redux。

理由：

- 游戏客户端状态是“长生命周期单房间状态 + 少量 UI 局部状态”
- Zustand 足够轻，适合按切片组织
- 比 Redux ceremony 少很多
- 与 WebSocket 推送结合自然

### TanStack Query

推荐只用于 REST 生命周期数据，而不是整个游戏实时态。

适合处理：

- 创建房间
- 加入房间
- 恢复会话
- 获取初始快照
- 非实时配置读取

不建议用 React Query 直接承载整局实时游戏状态，因为房间内核心状态应由 WebSocket 驱动。

### Tailwind CSS + Radix UI

推荐组合：

- Tailwind 负责高效布局与响应式
- Radix 负责可访问性强的基础交互组件，如 Dialog、Popover、Tabs、Toast

理由：

- 桌游 UI 自定义度高，不适合被重型组件库绑住
- 但弹窗、菜单、提示等基础交互又不值得手写一遍

### react-zoom-pan-pinch

用于地图区域缩放与拖拽。

理由：

- 山屋地图会动态增长
- PC 需要平滑拖拽缩放
- 移动端需要双指缩放和视口重置
- 自己实现会把精力浪费在手势细节上

### Go + gin + huma

采用 `gin` 作为 HTTP 路由框架，搭配 `huma` 自动生成 OpenAPI 3.x 文档。

理由：

- `gin` 高性能、生态丰富、社区活跃
- `huma` 通过 `humagin` 适配器包装 gin，可自动暴露 `/openapi.json`、`/openapi.yaml` 等标准 OpenAPI 端点
- Handler 通过类型化的 Input/Output 结构体定义请求/响应，编译期即保证参数一致性
- 内置参数校验、错误格式化、内容协商，减少样板代码
- 对 WebSocket 端点可直接使用 gin 原生路由注册，不受 huma 约束

### gws

采用 `github.com/lxzan/gws` 作为 WebSocket 库。

理由：

- 事件驱动模型（实现 `gws.Event` 接口），代码组织更清晰
- 内置并发安全写、Ping/Pong 心跳管理
- 高性能、低内存占用
- 支持 `ParallelEnabled` 并行消息处理
- 替代已归档的 `gorilla/websocket` 和 `coder/websocket`

### PostgreSQL + pgx + sqlc

推荐组合：

- PostgreSQL 作为唯一数据库
- `pgx` 驱动
- `sqlc` 生成类型安全的数据访问层

理由：

- 房间、玩家座位、会话、日志、快照都适合关系模型 + JSONB 混合建模
- `sqlc` 比 ORM 更适合规则复杂、结构明确的系统
- 对事件日志和快照 JSONB 支持好

不推荐首版使用 GORM，原因是：

- 模型边界明确，不需要 ORM 的动态便利
- 事件表、快照表、查询约束更适合手写 SQL
- 类型安全和可控 SQL 对可维护性更重要

### 静态内容文件而不是数据库内容配置

房间、卡牌、角色、作祟属于**版本化静态内容**，不属于玩家运营数据。

推荐：

- 把内容放在仓库中的 `content/` 目录
- 以 JSON/YAML 描述
- 构建时或服务启动时校验

理由：

- 内容跟代码一起版本管理
- 更利于 review、diff、测试与回滚
- 不需要为了静态内容增加后台 CMS 复杂度

### Caddy

推荐用作前端静态文件与反向代理入口。

理由：

- 配置比 Nginx 更简单
- docker compose 场景下易用
- 能统一代理 `/api` 和 `/ws`

### 明确不引入的技术

首版不建议引入：

- Redis
- Kafka / RabbitMQ
- 微服务拆分
- GraphQL
- 分布式锁

原因：

- 这些都不解决本项目真正的首要问题
- 3-6 人房间制桌游，用单进程房间 actor 足够
- 复杂基础设施会显著拖慢规则实现和联调速度
