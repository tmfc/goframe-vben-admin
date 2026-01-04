# Backend Development Plan - TODO List

## Phase 1: Core User and Authentication Services
- [x] Database Schema Design - User Table (`sys_user`)
- [x] GoFrame Model/DAO Generation for `sys_user`
- [x] Authentication API Implementation (`/auth/login`, `/auth/refresh`, `/auth/logout`, `/auth/codes`)
- [x] User Information API Implementation (`/user/info`)
- [x] Tenant support implementation

## Phase 2: RBAC (Roles and Permissions)
- [x] Database Schema Design (`sys_role`, `sys_permission`, `sys_user_role`, `sys_role_permission`, `casbin_rule`)
- [x] GoFrame Model/DAO Generation for RBAC tables
- [x] RBAC Services (role CRUD, permission CRUD, assign roles to users, assign permissions to roles)
- [x] RBAC API Implementation (role endpoints, permission endpoints)
- [x] Authentication and Authorization Middleware with Casbin
- [x] User-Role association functionality
- [x] Role-Permission association functionality
- [x] Casbin policy synchronization

## Phase 3: Menu and Routing Management
- [x] Database Schema Design - Menu Table (`sys_menu`)
- [x] GoFrame Model/DAO Generation for `sys_menu`
- [x] Menu Services (CRUD operations)
- [x] Hierarchical menu tree generation logic
- [x] Dynamic routing for frontend integration
- [x] Menu API Implementation (`/menu/all`, `/sys-menu/*`)

## Phase 4: Department Management
- [x] Database Schema Design - Department Table (`sys_dept`)
- [x] GoFrame Model/DAO Generation for `sys_dept`
- [x] Department Services (CRUD operations)
- [x] Hierarchical department data logic
- [x] Department API Implementation (`/sys-dept/*`)

## General Tasks
- [x] Error Handling
- [x] Input Validation
- [x] Logging
- [x] API Documentation (Swagger/OpenAPI) - GoFrame auto-generated, configured in config.yaml
- [x] Comprehensive unit and integration tests - 13 test files covering services, controllers, middleware, and E2E scenarios

## Optional Enhancements
- [x] Add unit tests for Department service
- [x] Add unit tests for Menu controller
- [x] Add integration tests for Department API endpoints
- [x] Add integration tests for Menu API endpoints
