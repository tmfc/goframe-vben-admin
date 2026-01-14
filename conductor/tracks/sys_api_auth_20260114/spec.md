# 规范：API 权限绑定与自动化校验 (Casbin)

## 1. 概述 (Overview)
本 Track 的目标是增强系统现有的 RBAC（基于角色的访问控制）权限系统，实现“权限标识”与“后端 API 接口”的动态绑定。通过在全局中间件中自动拦截请求并进行 Casbin 权限校验，确保只有拥有特定权限的用户（或角色）才能访问受限的 API 接口。

## 2. 功能要求 (Functional Requirements)

### 2.1 模型增强 (Model Enhancement)
- 在现有的 `permission`（或菜单/权限）模型中增加 `apis` 字段。
- 该字段类型为 JSON，用于存储该权限关联的 API 列表。
- 存储格式示例：
  ```json
  [
    {"path": "/api/v1/user/:id", "method": "GET"},
    {"path": "/api/v1/user/:id", "method": "PUT"}
  ]
  ```

### 2.2 后台管理界面 (Admin UI)
- 在权限管理/菜单管理界面中，增加一个配置项，允许管理员为每个权限项添加、编辑或删除绑定的 API。
- 需要提供 API 路径（支持参数，如 `:id`）和 HTTP 方法（GET, POST, PUT, DELETE, ALL）的选择。

### 2.3 权限同步逻辑 (Permission Sync)
- 当权限与 API 的绑定关系发生变化时，系统应自动更新 Casbin 的策略库（Policy）。
- Casbin 策略格式建议：`p, role_id, path, method`。

### 2.4 全局权限中间件 (Global Authorization Middleware)
- 实现一个全局中间件，用于拦截所有受保护的 API 请求。
- **匹配逻辑：**
    1. 获取当前请求的 `Path` 和 `Method`。
    2. 使用 GoFrame 的路由匹配引擎（或兼容的逻辑），将当前 `Path` 与 Casbin 策略中的路由模式（如 `/user/:id`）进行匹配。
    3. 调用 Casbin 的 `Enforce` 方法校验当前登录用户的角色是否拥有访问该 `(path, method)` 的权限。
- **响应处理：**
    - 校验通过：继续执行后续业务逻辑。
    - 校验失败：返回 `403 Forbidden` 错误。

## 3. 非功能要求 (Non-Functional Requirements)
- **性能：** 由于中间件拦截所有请求，API 路径匹配逻辑应进行优化（如使用缓存或高效的路由树）。
- **兼容性：** 现有的菜单展示逻辑不受影响。

## 4. 验收标准 (Acceptance Criteria)
- [ ] 可以在后台界面为某个菜单项绑定具体的 API 接口。
- [ ] 绑定完成后，未分配该角色的用户访问该接口应返回 403。
- [ ] 分配了对应角色的用户可以正常访问该接口。
- [ ] 支持动态路径匹配（如 `/v1/sys_dict/data/1` 能够匹配配置的 `/v1/sys_dict/data/:id`）。

## 5. 超出范围 (Out of Scope)
- 接口级别的细粒度数据权限过滤（本阶段仅关注 API 访问控制）。
