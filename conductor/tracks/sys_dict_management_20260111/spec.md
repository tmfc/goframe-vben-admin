# Track Specification: System Dictionary Management (Backend)

## 1. 概述
本 Track 旨在实现系统字典管理的后端功能。字典管理用于维护系统中的固定、可枚举数据（如性别、状态、行业分类等），支持多语言动态扩展和高性能缓存。

## 2. 功能需求

### 2.1 数据库设计与迁移
- 创建 `sys_dict_type` 表（字典类型）。
- 创建 `sys_dict_data` 表（字典数据）。
- **多语言重构**: `sys_dict_data` 使用 `label` (VARCHAR) 作为默认显示，增加 `label_i18n` (JSONB) 存储多语言映射（如 `{"en": "Male", "zh-TW": "男"}`）。
- 执行初始化数据（Seed Data）。
- **表结构参考**: `plan/dict.md` (需根据 JSONB 方案进行微调)。

### 2.2 管理端 API (CRUD)
需实现以下 RESTful 接口，并对接 RBAC 权限系统：
- **字典类型**: 完整 CRUD 接口（Create, Read, Update, Delete, List）。
- **字典数据**: 完整 CRUD 接口。需支持 `label_i18n` 字段的输入与更新。

### 2.3 公共查询 API
面向前端组件调用的只读接口：
- `GET /sys-dict/options/{typeCode}`: 获取选项列表。需根据 `Accept-Language` 自动从 `label_i18n` 或 `label` 提取文本。
- `GET /sys-dict/label/{typeCode}/{value}`: 根据类型和值获取标签。同样支持国际化自动回退。

### 2.4 核心逻辑
- **国际化 (I18n)**: 实现 I18n 工具函数，逻辑为：`label_i18n[lang] ?? label`。
- **缓存机制**: 实现缓存层（Redis 或内存缓存），读取字典数据优先查缓存，变更时自动失效。
- **审计记录**: 在创建和更新操作中自动记录 `creator_id`, `modifier_id`, `created_at`, `updated_at`。

## 3. 非功能需求
- **性能**: 字典读取接口必须经过缓存优化。
- **扩展性**: 通过 JSONB 支持任意数量的语言种类。
- **规范**: 遵循 GoFrame 项目的 Controller-Service-Dao 架构模式。

## 4. 排除范围
- 前端管理界面与 UI 组件开发。
