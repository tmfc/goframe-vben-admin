# Implementation Plan: System Dictionary Management (Backend)

## Phase 1: 数据库与实体生成 [checkpoint: 28d0baa]
- [x] Task: 创建数据库迁移文件（支持 JSONB 字段 `label_i18n`） [bb467ba]
- [x] Task: 执行迁移并验证表结构 [bb467ba]
- [x] Task: 编写并执行初始化 Seed 数据脚本 [bb467ba]
- [x] Task: 使用 `gf gen dao` 生成实体、DAO 和模型文件 [bb467ba]
- [x] Task: Conductor - User Manual Verification 'Phase 1: 数据库与实体生成' (Protocol in workflow.md) [bb467ba]

## Phase 2: 核心业务逻辑实现 (TDD) [checkpoint: 28d0baa]
- [x] Task: 编写 `SysDictType` Service 单元测试并实现 CRUD 逻辑 [bb467ba]
- [x] Task: 编写 `SysDictData` Service 单元测试并实现 CRUD 逻辑 [bb467ba]
- [x] Task: 实现 I18n 标签提取工具函数 (`label_i18n[lang] ?? label`) [bb467ba]
- [x] Task: 实现字典数据缓存层（Redis/内存），确保 CRUD 时同步清理 [bb467ba]
- [x] Task: 使用 `gf gen service` 生成服务接口定义 [bb467ba]
- [x] Task: Conductor - User Manual Verification 'Phase 2: 核心业务逻辑实现 (TDD)' (Protocol in workflow.md) [bb467ba]

## Phase 3: API 与控制器实现 (TDD) [checkpoint: 28d0baa]
- [x] Task: 在 `backend/api/` 定义请求与响应结构体（含多语言字段处理） [bb467ba]
- [x] Task: 编写 Controller 单元测试并实现管理端 CRUD 接口 [bb467ba]
- [x] Task: 编写公共查询 API (`/options`, `/label`) 并实现逻辑 [bb467ba]
- [x] Task: 在 `cmd/` 中注册路由并配置中间件 [bb467ba]
- [x] Task: Conductor - User Manual Verification 'Phase 3: API 与控制器实现 (TDD)' (Protocol in workflow.md) [bb467ba]

## Phase 4: 权限集成与最终验证 [checkpoint: 28d0baa]
- [x] Task: 配置 RBAC 权限（Casbin），限制管理端 API 访问 [bb467ba]
- [x] Task: 验证审计字段（`creator_id` 等）在创建/更新时是否自动填充 [388d0a0]
- [x] Task: 进行端到端（E2E）集成测试，验证缓存、I18n 与数据库的一致性 [388d0a0]
- [ ] Task: Conductor - User Manual Verification 'Phase 4: 权限集成与最终验证' (Protocol in workflow.md)
