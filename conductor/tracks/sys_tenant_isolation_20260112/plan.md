# Implementation Plan: 租户隔离配套工作及登录逻辑改造

本计划旨在实现多租户隔离的基础设施补全及登录认证逻辑的升级。

## Phase 1: 数据库与模型适配 (Database & Model Adaptation)
补全业务表的 `tenant_id` 字段并完善租户基础信息。

- [ ] Task: 生成补全 `tenant_id` 的数据库迁移文件
- [ ] Task: 在 `sys_tenant` 中增加 `code` 字段并建立唯一索引
- [ ] Task: 更新 Go 后端业务模型 (internal/model/entity) 以包含 `tenant_id` 字段
- [ ] Task: Conductor - User Manual Verification 'Phase 1: 数据库与模型适配' (Protocol in workflow.md)

## Phase 2: 配置管理与基础支持 (Config & Core Support)
增加多租户开关并实现基础工具类。

- [ ] Task: 在 `config.toml` 中增加 `app.multiTenant` 配置项
- [ ] Task: TDD - 编写解析登录名的工具函数（需支持有无后缀及模式开关两种情况）
- [ ] Task: 实现登录名解析逻辑（若开关关闭或无后缀，则返回默认租户 ID 1）
- [ ] Task: Conductor - User Manual Verification 'Phase 2: 配置管理与基础支持' (Protocol in workflow.md)

## Phase 3: 登录逻辑与租户隔离增强 (Login & Tenant Isolation)
支持跨租户同名用户登录及多租户过滤逻辑优化。

- [ ] Task: TDD - 编写多模式下的登录认证单元测试
- [ ] Task: 改造登录 Service：处理单租户（直接查）与多租户（按 code + name 查）的分支逻辑
- [ ] Task: 优化多租户中间件：当开关关闭时，强制注入 `tenant_id = 1` 过滤条件
- [ ] Task: Conductor - User Manual Verification 'Phase 3: 登录逻辑与租户隔离增强' (Protocol in workflow.md)

## Phase 4: 超管平台管理与身份切换 (Admin Management & Identity Switching)
实现平台级功能访问及超管切换租户 Token 逻辑。

- [ ] Task: 预设 `system` 租户作为平台管理租户
- [ ] Task: TDD - 编写超管切换租户身份（签发新 Token）的单元测试
- [ ] Task: 实现切换租户接口：验证权限并重新签发携带目标 `tenant_id` 的 JWT
- [ ] Task: 调整菜单权限逻辑：允许 `system` 租户访问平台级管理功能
- [ ] Task: Conductor - User Manual Verification 'Phase 4: 超管平台管理与身份切换' (Protocol in workflow.md)

## Phase 5: 前端适配与验收 (Frontend Adaptation & Final Verification)
调整登录界面及管理交互。

- [ ] Task: 修改前端登录页，根据配置调整提示语并引导用户使用新格式
- [ ] Task: 实现前端超管切换租户的 UI 交互
- [ ] Task: 进行单/多租户模式切换的全量验收测试
- [ ] Task: Conductor - User Manual Verification 'Phase 5: 前端适配与验收' (Protocol in workflow.md)
