# 消息中心实现计划 (plan.md)

## 阶段 1: 数据库设计与后端基础 (Phase 1: Database & Backend Foundation) [checkpoint: 5a9504a]
- [x] 任务: 创建数据库迁移文件 (SQL Migration) [52bb0dd]
    - [x] 创建 `sys_message` 表存储消息主体
    - [x] 创建 `sys_message_read` 表记录用户已读状态
- [x] 任务: 实现 DAO 和 Model 层 (GoFrame dao/model) [4861eb2]
    - [x] 使用 `gf gen dao` 生成基础代码
    - [x] 定义消息传输对象 (DTO) 和业务常量 (消息类型、跳转类型)
- [x] 任务: 核心业务逻辑实现 (Logic 层) [67e6ae1]
    - [x] 编写发送消息的基础接口 (支持异步推送到 Redis)
    - [x] 编写消息列表查询、标记已读的业务逻辑
- [x] 任务: 单元测试 [67e6ae1]
    - [x] 针对消息发送、读取、已读状态更新编写单元测试
- [x] Task: Conductor - User Manual Verification 'Phase 1: Database & Backend Foundation' (Protocol in workflow.md) [5a9504a]

## 阶段 2: 异步机制与 Redis 集成 (Phase 2: Async & Redis Integration)
- [x] 任务: 配置 Redis 集成 [a26098c]
    - [x] 在 `internal/config` 中确保 Redis 配置正确
- [ ] 任务: 实现消息生产者与消费者
    - [ ] 业务模块作为生产者，将消息推入 Redis 队列
    - [ ] 实现后台消费协程，监听队列并准备分发到 WebSocket
- [ ] 任务: 健壮性测试
    - [ ] 编写测试用例验证消息在 Redis 中的堆积与正确消费
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Async & Redis Integration' (Protocol in workflow.md)

## 阶段 3: WebSocket 服务端实现 (Phase 3: WebSocket Server)
- [ ] 任务: 实现 WebSocket 控制器 (Controller)
    - [ ] 处理客户端连接、鉴权与心跳
- [ ] 任务: 实现连接管理器 (Hub/Manager)
    - [ ] 维护用户 ID 与 WebSocket 连接的映射关系
    - [ ] 实现点对点 (P2P) 和广播 (Broadcast) 功能
- [ ] 任务: 消费 Redis 消息并推送
    - [ ] 将 Redis 消费者获取的消息通过 WebSocket 发送给指定用户
- [ ] 任务: 集成测试
    - [ ] 使用测试脚本模拟 WebSocket 连接并验证消息实时接收
- [ ] Task: Conductor - User Manual Verification 'Phase 3: WebSocket Server' (Protocol in workflow.md)

## 阶段 4: 消息中心 API 开发 (Phase 4: Message Center API)
- [ ] 任务: 编写 API 接口文档 (api/v1)
    - [ ] 定义消息列表分页、详情、设为已读、全部已读等接口
- [ ] 任务: 实现 API 控制器层
    - [ ] 调用 Logic 层完成业务闭环
- [ ] 任务: 权限配置 (Casbin)
    - [ ] 为消息中心相关 API 配置 RBAC 权限
- [ ] 任务: 自动化测试 (E2E)
    - [ ] 验证全流程 API 调用
- [ ] Task: Conductor - User Manual Verification 'Phase 4: Message Center API' (Protocol in workflow.md)

## 阶段 5: 前端集成 - Vben Admin (Phase 5: Frontend Integration)
- [ ] 任务: 实现前端 WebSocket 客户端封装
    - [ ] 维护连接状态、重连逻辑
- [ ] 任务: 实现消息通知组件 (Header Bell)
    - [ ] 实时更新未读总数，展示最近几条通知
- [ ] 任务: 实现消息中心管理页面
    - [ ] 包含分类列表、已读/未读筛选、一键已读功能
- [ ] 任务: 跳转逻辑实现
    - [ ] 根据后端返回的 `Path` 和 `Params` 执行路由跳转或打开弹窗
- [ ] Task: Conductor - User Manual Verification 'Phase 5: Frontend Integration' (Protocol in workflow.md)
