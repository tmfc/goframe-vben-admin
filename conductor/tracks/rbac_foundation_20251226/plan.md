# Plan: RBAC Data Tables and API Foundation

## Phase 1: Database Schema and DAO Generation

- [x] Task: Design and create `sys_role` table
    - [x] Sub-task: Create migration file for `sys_role` table with `id`, `name`, `description`, `parent_id`, `status`, `created_at`, `updated_at`, `creator_id`, `modifier_id`, `dept_id`.
    - [x] Sub-task: Run migration to create `sys_role` table.
    - [x] Sub-task: Generate DAO for `sys_role`.
- [x] Task: Design and create `sys_permission` table
    - [x] Sub-task: Create migration file for `sys_permission` table with `id`, `name`, `description`, `parent_id`, `status`, `created_at`, `updated_at`, `creator_id`, `modifier_id`, `dept_id`.
    - [x] Sub-task: Run migration to create `sys_permission` table.
    - [x] Sub-task: Generate DAO for `sys_permission`.
- [x] Task: Design and create Role-Permission Pivot Table (`sys_role_permission`)
    - [x] Sub-task: Create migration file for `sys_role_permission` table with `role_id`, `permission_id`, `created_at`, `updated_at`, `scope`.
    - [x] Sub-task: Run migration to create `sys_role_permission` table.
    - [x] Sub-task: Generate DAO for `sys_role_permission`.
- [x] Task: Update `sys_user` table with metadata fields
    - [x] Sub-task: Create migration file to add `created_at`, `updated_at`, `creator_id`, `modifier_id`, `dept_id` to `sys_user`.
    - [x] Sub-task: Run migration to update `sys_user` table.
- [x] Task: Update `sys_menu` table with metadata fields
    - [x] Sub-task: Create migration file to add `created_at`, `updated_at`, `creator_id`, `modifier_id`, `dept_id` to `sys_menu`.
    - [x] Sub-task: Run migration to update `sys_menu` table.
- [x] Task: Implement default `casbin_rule` table structure
    - [x] Sub-task: Create migration file for `casbin_rule` table with `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`.
    - [x] Sub-task: Run migration to create `casbin_rule` table.
    - [x] Sub-task: Generate DAO for `cas_bin_rule`.
- [x] Task: Conductor - User Manual Verification 'Phase 1: Database Schema and DAO Generation' (Protocol in workflow.md)

## Phase 2: Initial Role and Permission Services & APIs

- [x] Task: Implement `sys_role` service (ee644cb)
    - [ ] Sub-task: Write tests for `sys_role` service.
    - [ ] Sub-task: Implement `CreateRole` function.
    - [ ] Sub-task: Implement `GetRole` function.
    - [ ] Sub-task: Implement `UpdateRole` function.
    - [ ] Sub-task: Implement `DeleteRole` function.
- [ ] Task: Implement `sys_permission` service
    - [ ] Sub-task: Write tests for `sys_permission` service.
    - [ ] Sub-task: Implement `CreatePermission` function.
    - [ ] Sub-task: Implement `GetPermission` function.
    - [ ] Sub-task: Implement `UpdatePermission` function.
    - [ ] Sub-task: Implement `DeletePermission` function.
- [ ] Task: Implement `sys_role` API
    - [ ] Sub-task: Write tests for `sys_role` API.
    - [ ] Sub-task: Implement `CreateRole` API endpoint.
    - [ ] Sub-task: Implement `GetRole` API endpoint.
    - [ ] Sub-task: Implement `UpdateRole` API endpoint.
    - [ ] Sub-task: Implement `DeleteRole` API endpoint.
- [ ] Task: Implement `sys_permission` API
    - [ ] Sub-task: Write tests for `sys_permission` API.
    - [ ] Sub-task: Implement `CreatePermission` API endpoint.
    - [ ] Sub-task: Implement `GetPermission` API endpoint.
    - [ ] Sub-task: Implement `UpdatePermission` API endpoint.
    - [ ] Sub-task: Implement `DeletePermission` API endpoint.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Initial Role and Permission Services & APIs' (Protocol in workflow.md)
