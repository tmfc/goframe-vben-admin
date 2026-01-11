# 实施计划 (plan.md) - 系统模块重构与优化

## Phase 1: 类型定义与工具提取 (Foundations)
本阶段旨在建立坚实的基础，通过完善类型系统和提取公共逻辑，为后续的 UI 重构做准备。

- [~] Task: 1.1 定义 Menu 和 Permission 的 TypeScript 接口
    - 在 `src/api/sys/menu.ts` 和 `src/api/sys/permission.ts` 中定义实体接口（如 `MenuItem`, `PermissionItem`）和请求/响应类型。
    - 替换 API 文件中的 `any` 类型。
- [ ] Task: 1.2 创建树形结构工具库
    - 在 `frontend/apps/web-naive/src/utils/tree.ts` 中实现 `listToTree`, `flattenTree`, `collectExpandedKeys` 等函数。
    - 编写单元测试验证工具函数的正确性。
- [ ] Task: 1.3 替换现有代码中的树形逻辑
    - 更新 `dept`, `role`, `user` 等模块，使用新的工具函数替换重复的内联逻辑。
    - 运行测试确保替换后功能未受影响。
- [ ] Task: Conductor - User Manual Verification 'Foundations' (Protocol in workflow.md)

## Phase 2: 图标选择器优化 (Icon Picker)
本阶段专注于安全性和离线支持，确保图标选择器不再依赖外部服务。

- [ ] Task: 2.1 引入本地图标集
    - 确定项目中需要的常用图标集（如 `ant-design` 等），并配置本地离线使用方式（可能利用 `@iconify/json` 或 Vben 的图标机制）。
- [ ] Task: 2.2 重构 MenuIconPicker 组件
    - 修改 `src/views/sys/menu/components/MenuIconPicker.vue`。
    - 移除 `fetch` 请求，改为从本地数据源加载图标。
- [ ] Task: Conductor - User Manual Verification 'Icon Picker' (Protocol in workflow.md)

## Phase 3: 菜单与权限模块重构 (Refactoring)
本阶段是核心工作，将菜单和权限模块迁移到标准的 Vben Form 架构。

- [ ] Task: 3.1 重构菜单表单 (Menu Form)
    - 创建 `src/views/sys/menu/modules/form.vue`。
    - 使用 `useVbenForm` 定义表单 Schema。
    - 迁移现有的权限代码自动生成逻辑 (`generatePermissionCode`)。
    - 更新 `src/views/sys/menu/index.vue` 使用新的弹窗组件。
- [ ] Task: 3.2 重构权限表单 (Permission Form)
    - 创建 `src/views/sys/permission/modules/form.vue`。
    - 使用 `useVbenForm` 定义表单 Schema。
    - 更新 `src/views/sys/permission/index.vue` 使用新的弹窗组件。
- [ ] Task: 3.3 完善测试用例
    - 确保 `menu` 和 `permission` 模块有基本的单元测试覆盖，特别是新的表单逻辑。
- [ ] Task: Conductor - User Manual Verification 'Refactoring' (Protocol in workflow.md)
