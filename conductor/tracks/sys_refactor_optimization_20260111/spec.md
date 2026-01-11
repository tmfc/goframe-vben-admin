# 规格说明书 (spec.md) - 系统模块重构与优化 (Refactor & Optimization)

## 1. 概述 (Overview)
本 Track 旨在提升 `web-naive` 前端应用中系统管理模块（主要是菜单和权限）的代码质量、一致性和安全性。通过重构表单实现方式、完善 TypeScript 类型定义、提取通用工具函数以及本地化图标选择器，使代码更符合 Vben Admin 的开发规范。

## 2. 功能需求 (Functional Requirements)

### 2.1 表单重构 (Form Refactoring)
- 将 `src/views/sys/menu/index.vue` 和 `src/views/sys/permission/index.vue` 中的内联表单抽离到各自目录下的 `modules/form.vue`。
- 使用 `useVbenForm` 替代原生的 Naive UI 表单实现，确保与 `user`、`role` 等模块的开发模式一致。

### 2.2 类型安全 (Type Safety)
- 定义 `MenuItem` 和 `PermissionItem` 接口。
- 为 `src/api/sys/menu.ts` 和 `src/api/sys/permission.ts` 中的 API 函数添加完整的参数和返回值类型定义。
- 消除视图文件和 API 文件中不必要的 `any` 使用。

### 2.3 树形工具函数 (Tree Helpers)
- 在 `frontend/apps/web-naive/src/utils` 下创建通用的树形处理工具函数。
- 提取并统一 `listToTree`、`flattenTree`、`collectExpandedKeys` 等逻辑。

### 2.4 图标选择器优化 (Icon Picker Security)
- 重构 `MenuIconPicker`，将常用的图标集本地化/离线化。
- 移除对外部 `api.iconify.design` 的直接网络请求，确保内网环境可用性及合规性。

## 3. 非功能需求 (Non-Functional Requirements)
- **代码一致性**: 遵循 Vben Admin 的最佳实践和项目既有的代码风格。
- **可维护性**: 通过逻辑抽离和类型约束降低后续维护成本。

## 4. 验收标准 (Acceptance Criteria)
- [ ] 菜单和权限模块的功能（增删改查）在重构后依然正常运行。
- [ ] `web-naive` 应用能够成功通过 TypeScript 类型检查 (`pnpm typecheck`)。
- [ ] `MenuIconPicker` 在断网或拦截外部请求的情况下依然能显示和选择本地图标。
- [ ] 树形相关的代码冗余得到消除。

## 5. 超出范围 (Out of Scope)
- 不涉及后端的业务逻辑修改。
- 不涉及系统管理模块之外的其他功能模块。
