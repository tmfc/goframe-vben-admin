# Implementation Plan: System Dictionary Management (Full-stack Enhancement)

## Phase 1: 后端接口与逻辑增强 (Backend APIs & Logic)
本阶段重点完善后端字典管理的能力，包括导入导出和原子化状态更新。

- [ ] Task: TDD - 编写字典数据导出接口单元测试并实现逻辑（支持 Excel/JSON）
- [ ] Task: TDD - 编写字典数据导入接口单元测试并实现逻辑（含事务处理与重复校验）
- [ ] Task: TDD - 编写字典/项状态切换（Status）的原子化更新接口单元测试并实现
- [ ] Task: 优化 `label_i18n` 处理逻辑，确保在各种操作（含导入导出）下多语言数据的一致性
- [ ] Task: 配置新增接口的 Casbin 权限规则
- [ ] Task: Conductor - User Manual Verification 'Phase 1: 后端接口与逻辑增强' (Protocol in workflow.md)

## Phase 2: 前端 API 与类型定义 (Frontend API & Types)
为 `web-naive` 定义与后端交互的接口及 TypeScript 类型。

- [ ] Task: 在 `frontend/apps/web-naive/src/api/sys/dict.ts` 中创建 API 方法（含导入导出、状态切换）
- [ ] Task: 定义字典类型及数据的 TypeScript 接口，特别加强对 `label_i18n` 的类型支持
- [ ] Task: Conductor - User Manual Verification 'Phase 2: 前端 API 与类型定义' (Protocol in workflow.md)

## Phase 3: 字典管理主页面开发 (Dictionary Management Main Page)
实现字典类型的列表展示及右侧管理面板的基础布局。

- [ ] Task: 实现字典类型管理主页面 `frontend/apps/web-naive/src/views/sys/dict/type/index.vue`
- [ ] Task: 实现基于 Naive UI 的字典类型搜索、分页列表及状态切换 Switch
- [ ] Task: 实现字典数据管理侧边栏（Drawer/Sidebar）容器及基础列表
- [ ] Task: Conductor - User Manual Verification 'Phase 3: 字典管理主页面开发' (Protocol in workflow.md)

## Phase 4: 多语言编辑与高级功能实现 (i18n & Advanced Features)
实现动态多语言表单及导入导出 UI。

- [ ] Task: 实现字典数据编辑表单中的动态多语言配置组件（根据已有数据自动带出语种）
- [ ] Task: 实现字典数据的增删改功能及权限按钮控制
- [ ] Task: 实现前端导入导出交互（上传文件、下载文件流）
- [ ] Task: 最终联调：验证全流程 CRUD、状态切换、导入导出及多语言回显
- [ ] Task: Conductor - User Manual Verification 'Phase 4: 多语言编辑与高级功能实现' (Protocol in workflow.md)
