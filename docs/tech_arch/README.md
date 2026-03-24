# 技术架构子文档索引

> 目录目的：将 [technical-architecture.md](../technical-architecture.md) 按后续开发可引用的边界拆分，避免在具体实施阶段频繁通读整篇总文档。

## 文档清单

- [01-overview.md](./01-overview.md)
  - 架构结论、系统分层、职责边界、整体数据流
  - 适合在进入任何开发阶段前先统一设计认知

- [02-technology-selection.md](./02-technology-selection.md)
  - 技术栈清单、选型理由、明确不引入的技术
  - 适合做基础设施落地、脚手架初始化、技术评审

- [03-backend-architecture.md](./03-backend-architecture.md)
  - 房间管理、`RoomRuntime`、状态管理、WebSocket、规则引擎、API 边界
  - 适合后端框架搭建、房间系统、引擎实现、联机同步阶段

- [04-frontend-architecture.md](./04-frontend-architecture.md)
  - 路由、页面骨架、组件拆分、前端状态、地图渲染、移动端适配
  - 适合前端界面搭建、地图视图、交互层实现阶段

- [05-core-challenges.md](./05-core-challenges.md)
  - 动态地图、作祟分叉、信息不对称、作祟配置化、恢复与动画一致性
  - 适合在实现高风险模块前集中对照

- [06-data-models.md](./06-data-models.md)
  - 持久化实体、运行时实体、静态内容模型、持久化与内存态划分
  - 适合数据库建模、内容结构设计、状态序列化阶段

- [07-delivery-and-deployment.md](./07-delivery-and-deployment.md)
  - 里程碑、MVP、优先级、`docker compose`、最终实施建议
  - 适合做开发排期、任务拆分、部署落地阶段

## 推荐阅读顺序

建议默认顺序：

1. [01-overview.md](./01-overview.md)
2. [02-technology-selection.md](./02-technology-selection.md)
3. [03-backend-architecture.md](./03-backend-architecture.md)
4. [04-frontend-architecture.md](./04-frontend-architecture.md)
5. [05-core-challenges.md](./05-core-challenges.md)
6. [06-data-models.md](./06-data-models.md)
7. [07-delivery-and-deployment.md](./07-delivery-and-deployment.md)

## 按开发阶段引用建议

### 阶段 1：房间、连接、基础骨架

优先阅读：

- [01-overview.md](./01-overview.md)
- [02-technology-selection.md](./02-technology-selection.md)
- [03-backend-architecture.md](./03-backend-architecture.md)
- [06-data-models.md](./06-data-models.md)
- [07-delivery-and-deployment.md](./07-delivery-and-deployment.md)

### 阶段 2：作祟前核心引擎

优先阅读：

- [01-overview.md](./01-overview.md)
- [03-backend-architecture.md](./03-backend-architecture.md)
- [05-core-challenges.md](./05-core-challenges.md)
- [06-data-models.md](./06-data-models.md)

### 阶段 3：前端地图与游戏主界面

优先阅读：

- [01-overview.md](./01-overview.md)
- [04-frontend-architecture.md](./04-frontend-architecture.md)
- [05-core-challenges.md](./05-core-challenges.md)

### 阶段 4：作祟、隐藏信息、怪物与内容扩展

优先阅读：

- [03-backend-architecture.md](./03-backend-architecture.md)
- [05-core-challenges.md](./05-core-challenges.md)
- [06-data-models.md](./06-data-models.md)

### 阶段 5：联调、排错、恢复与部署

优先阅读：

- [03-backend-architecture.md](./03-backend-architecture.md)
- [06-data-models.md](./06-data-models.md)
- [07-delivery-and-deployment.md](./07-delivery-and-deployment.md)

## 说明

- 原始总文档仍保留在 [technical-architecture.md](../technical-architecture.md)，适合整体验证与全局回顾。
- 后续如果继续扩展架构文档，建议优先维护子文档，再视需要回写总文档。  
