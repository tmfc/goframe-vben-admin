# 实施计划：API 权限绑定与自动化校验 (Casbin)

## 阶段 1：基础设施与模型更新 (Infrastructure & Model Update)
- [ ] **Task: 数据库变更**
    - [ ] 子任务：创建 SQL 迁移文件，在 `sys_permission` 表中增加 `apis` (JSONB) 字段。
    - [ ] 子任务：执行迁移。
- [ ] **Task: 模型与 DAO 更新**
    - [ ] 子任务：运行 `gf gen dao` 更新后端模型和 DAO 代码。
- [ ] **Task: Conductor - User Manual Verification 'Infrastructure & Model Update' (Protocol in workflow.md)**

## 阶段 2：核心授权与同步逻辑 (Core Authorization Logic)
- [ ] **Task: Casbin 策略升级**
    - [ ] 子任务：更新 Casbin 模型配置（如果需要），支持 `(sub, obj, act)` 匹配，其中 `obj` 是路径模式，`act` 是 HTTP 方法。
    - [ ] 子任务：实现支持路由模式匹配（如 `:id`）的 `Enforce` 扩展逻辑。
- [ ] **Task: 权限同步服务**
    - [ ] 子任务：编写测试用例：验证将 `sys_permission.apis` 同步至 Casbin Policy 的正确性。
    - [ ] 子任务：实现同步逻辑（在权限创建/更新/删除时触发同步）。
- [ ] **Task: Conductor - User Manual Verification 'Core Authorization Logic' (Protocol in workflow.md)**

## 阶段 3：中间件实现与集成 (Middleware Integration)
- [ ] **Task: 后端权限中间件**
    - [ ] 子任务：编写失败测试：尝试访问受限 API，预期返回 403。
    - [ ] 子任务：实现全局中间件 `MiddlewareApiAuth`，提取请求信息并调用 Casbin 校验。
    - [ ] 子任务：在路由组中注册该中间件。
- [ ] **Task: Conductor - User Manual Verification 'Middleware Integration' (Protocol in workflow.md)**

## 阶段 4：前端管理界面 (Frontend Management UI)
- [ ] **Task: 权限配置界面更新 (web-naive)**
    - [ ] 子任务：在 `web-naive` 的权限/菜单编辑表单中增加 API 绑定功能（支持动态增加 Path 和 Method 的列表）。
    - [ ] 子任务：确保 JSON 数据能够正确序列化并提交至后端。
- [ ] **Task: Conductor - User Manual Verification 'Frontend Management UI' (Protocol in workflow.md)**

## 阶段 5：全链路验证与交付 (Full System Verification)
- [ ] **Task: 集成测试与端到端验证**
    - [ ] 子任务：编写全链路集成测试：创建权限 -> 绑定 API -> 分配角色 -> 测试不同角色的访问权限。
- [ ] **Task: Conductor - User Manual Verification 'Full System Verification' (Protocol in workflow.md)**
