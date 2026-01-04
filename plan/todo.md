# 多 Agent 并行开发计划

## Agent 1: 用户-角色关联功能
- [x] 创建 `sys_user_role` 表迁移文件（如不存在）
- [x] 在 `sys_role` 服务中添加用户角色分配功能
  - [x] `AssignRoleToUser(ctx, userID, roleID)` 方法
  - [x] `RemoveRoleFromUser(ctx, userID, roleID)` 方法
  - [x] `GetUserRoles(ctx, userID)` 方法
- [x] 在 `sys_role` 服务中添加批量操作
  - [x] `AssignRolesToUser(ctx, userID, roleIDs)` 方法
  - [x] `GetUsersByRole(ctx, roleID)` 方法
- [x] 为用户角色关联功能编写单元测试

## Agent 2: 角色-权限关联功能
- [ ] 创建 `sys_role_permission` 表迁移文件（如不存在）
- [ ] 在 `sys_role` 服务中添加角色权限分配功能
  - [ ] `AssignPermissionToRole(ctx, roleID, permissionID)` 方法
  - [ ] `RemovePermissionFromRole(ctx, roleID, permissionID)` 方法
  - [ ] `GetRolePermissions(ctx, roleID)` 方法
- [ ] 在 `sys_role` 服务中添加批量操作
  - [ ] `AssignPermissionsToRole(ctx, roleID, permissionIDs)` 方法
  - [ ] `GetPermissionsByUser(ctx, userID)` 方法（通过角色间接获取）
- [ ] 为角色权限关联功能编写单元测试

## Agent 3: Casbin 策略同步
- [ ] 在 `sys_role` 服务中添加 Casbin 策略同步方法
  - [ ] `SyncRoleToCasbin(ctx, roleID)` - 将角色权限同步到 Casbin
  - [ ] `RemoveRoleFromCasbin(ctx, roleID)` - 从 Casbin 删除角色策略
  - [ ] `SyncAllRolesToCasbin(ctx)` - 同步所有角色到 Casbin
- [ ] 在角色创建/更新/删除时自动触发 Casbin 同步
- [ ] 为 Casbin 同步功能编写单元测试

## Agent 4: API Controller 增强
- [ ] 增强 `sys_role` controller
  - [ ] `POST /role/:id/assign-users` - 分配用户到角色
  - [ ] `POST /role/:id/remove-users` - 从角色移除用户
  - [ ] `GET /role/:id/users` - 获取角色下的用户列表
  - [ ] `POST /role/:id/assign-permissions` - 分配权限到角色
  - [ ] `POST /role/:id/remove-permissions` - 从角色移除权限
  - [ ] `GET /role/:id/permissions` - 获取角色的权限列表
- [ ] 增强 `sys_permission` controller
  - [ ] `GET /permission/by-user/:userId` - 获取用户的所有权限
- [ ] 为新增的 API 端点编写单元测试

## Agent 5: 中间件优化和测试
- [ ] 优化 Casbin 认证中间件
  - [ ] 添加缓存机制减少数据库查询
  - [ ] 改进错误处理和日志记录
- [ ] 创建中间件测试用例
- [ ] 测试不同角色和权限组合的授权场景
- [ ] 测试跨租户隔离是否正常工作

## Agent 6: 集成测试和文档
- [ ] 编写端到端集成测试
  - [ ] 测试完整的 RBAC 流程（创建角色、分配权限、分配用户、验证权限）
  - [ ] 测试多租户场景下的权限隔离
- [ ] 更新 API 文档（如果需要）
- [ ] 性能测试（大量角色/权限场景）
- [ ] 修复测试中发现的问题

---

## 依赖关系说明
- **Agent 1 和 Agent 2** 可以完全并行执行，各自处理不同的关联表
- **Agent 3** 依赖 Agent 2 的完成（需要角色-权限关联功能）
- **Agent 4** 依赖 Agent 1 和 Agent 2 的完成
- **Agent 5** 可以与 Agent 3 并行执行
- **Agent 6** 依赖所有其他 Agent 的完成

## 执行顺序建议
1. **第一批（并行）**: Agent 1, Agent 2, Agent 5
2. **第二批（并行）**: Agent 3（等待 Agent 2 完成）
3. **第三批（并行）**: Agent 4（等待 Agent 1 和 Agent 2 完成）
4. **第四批**: Agent 6（等待所有其他完成）