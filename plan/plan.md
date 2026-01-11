# Backend Development Plan for GoFrame + Vben Admin

## Objective

Implement the backend services and database schema to support the user, department, RBAC, and menu/routing features required by the Vben Admin frontend.

## Key Technologies

-   **Backend:** GoFrame
-   **Database:** MySQL (or similar relational database)

## Detailed Plan

### Phase 1: Core User and Authentication Services

**Database Schema Design - User Table (`sys_user`)**:

*   `id` (PK, BIGSERIAL) - Unique identifier for the user.
*   `username` (UNIQUE, VARCHAR) - User's login username.
*   `password` (VARCHAR) - Hashed password.
*   `salt` (VARCHAR) - Salt for password hashing.
*   `real_name` (VARCHAR) - User's real name.
*   `avatar` (VARCHAR) - URL to user's avatar image.
*   `home_path` (VARCHAR) - Default dashboard or user-specific home path.
*   `status` (TINYINT) - User account status (e.g., active, inactive, disabled).
*   `created_at` (DATETIME) - Timestamp for creation.
*   `updated_at` (DATETIME) - Timestamp for last update.
*   `deleted_at` (DATETIME) - Soft deletion timestamp.
*   `roles` (JSON or separate `sys_user_role` pivot table) - Initially, can store as JSON, later move to pivot table for normalization.

**GoFrame Model/DAO Generation for `sys_user`**:

*   Use `gf gen dao` to generate DAO/Model code for the `sys_user` table.

**Authentication API Implementation (`app/service/auth` & `app/controller/auth`)**:

*   **`POST /auth/login`**:
    *   **Input**: `username`, `password`
    *   **Output**: `accessToken`, `refreshToken`, `userInfo` (including `roles`)
    *   **Logic**: Validate credentials, generate JWT tokens, retrieve `userInfo` and associated `roles`.
*   **`POST /auth/refresh`**:
    *   **Input**: `refreshToken`
    *   **Output**: `accessToken`
    *   **Logic**: Validate `refreshToken`, generate new `accessToken`.
*   **`POST /auth/logout`**:
    *   **Logic**: Invalidate tokens (if using server-side token invalidation).
*   **`GET /auth/codes`**:
    *   **Output**: `string[]` (permission codes)
    *   **Logic**: Retrieve `accessCodes` associated with the authenticated user's roles.

**User Information API Implementation (`app/service/user` & `app/controller/user`)**:

*   **`GET /user/info`**:
    *   **Output**: `UserInfo` (current authenticated user's details, including `roles` and `homePath`).

### Phase 2: RBAC (Roles and Permissions)

**Database Schema Design**:

*   **Role Table (`sys_role`)**:
    *   `id` (PK, BIGSERIAL)
    *   `tenant_id` (FK to `sys_tenant.id`) - Tenant scope for the role.
    *   `name` (UNIQUE, VARCHAR) - Role name (e.g., 'admin', 'user', 'super').
    *   `description` (VARCHAR) - Description of the role.
    *   `status` (TINYINT) - Role status.
    *   `created_at` (DATETIME), `updated_at` (DATETIME), `deleted_at` (DATETIME).
    *   Unique constraint on (`tenant_id`, `name`) to allow same role names across tenants.
*   **Permission Table (`sys_permission`)**:
    *   `id` (PK, BIGSERIAL)
    *   `tenant_id` (FK to `sys_tenant.id`) - Tenant scope for the permission.
    *   `code` (UNIQUE, VARCHAR) - Permission code (e.g., 'user:create', 'menu:edit'), corresponds to `accessCodes` in frontend.
    *   `name` (VARCHAR) - Name of the permission.
    *   `description` (VARCHAR) - Description of the permission.
    *   `created_at` (DATETIME), `updated_at` (DATETIME), `deleted_at` (DATETIME).
    *   Unique constraint on (`tenant_id`, `code`) to allow same codes across tenants.
*   **User-Role Pivot Table (`sys_user_role`)**:
    *   `tenant_id` (FK to `sys_tenant.id`)
    *   `user_id` (FK to `sys_user.id`)
    *   `role_id` (FK to `sys_role.id`)
    *   Composite unique key on (`tenant_id`, `user_id`, `role_id`).
*   **Role-Permission Pivot Table (`sys_role_permission`)**:
    *   `tenant_id` (FK to `sys_tenant.id`)
    *   `role_id` (FK to `sys_role.id`)
    *   `permission_id` (FK to `sys_permission.id`)
    *   Composite unique key on (`tenant_id`, `role_id`, `permission_id`).
*   **Casbin Policy Table (`casbin_rule`)**:
    *   `id` (PK, BIGSERIAL)
    *   `ptype` (VARCHAR) - Policy type (`p`, `g`).
    *   `v0`..`v5` (VARCHAR) - Casbin rule fields; store `domain` (tenant) in `v1`.
    *   Indexes on `ptype`, `v0`, `v1`, `v2` for lookup.

**GoFrame Model/DAO Generation for RBAC tables**.

**RBAC Services (`app/service/role`, `app/service/permission`)**:

*   CRUD operations for roles and permissions.
*   Functions to assign roles to users.
*   Functions to assign permissions to roles.
*   Functions to retrieve all roles/permissions for a user.

**RBAC API Implementation (`app/controller/role`, `app/controller/permission`)**:

*   Endpoints for managing roles and permissions.

**Authentication and Authorization Middleware**:

*   Implement GoFrame middleware to validate JWT tokens for protected routes.
*   Implement GoFrame middleware to check user roles/permissions via Casbin (`domain=tenant_id`).

### Phase 3: Menu and Routing Management

**Database Schema Design - Menu Table (`sys_menu`)**:

*   `id` (PK, BIGSERIAL)
*   `parent_id` (FK to `sys_menu.id` for nested menus)
*   `name` (VARCHAR) - Menu item name.
*   `path` (VARCHAR) - Frontend route path.
*   `component` (VARCHAR) - Frontend component path (if applicable).
*   `icon` (VARCHAR) - Icon for the menu item.
*   `order` (INT) - Order for sorting menu items.
*   `type` (TINYINT) - Type of item (e.g., 'menu', 'button', 'iframe').
*   `visible` (TINYINT) - Whether to show in menu (0=false, 1=true).
*   `status` (TINYINT) - Menu item status.
*   `permission_code` (VARCHAR) - Related permission code for this menu item.
*   `meta` (JSON) - JSON field for additional frontend-specific metadata (e.g., `authority`, `menuVisibleWithForbidden`, `ignoreAuth`).
*   `created_at` (DATETIME), `updated_at` (DATETIME), `deleted_at` (DATETIME).

**GoFrame Model/DAO Generation for `sys_menu`**.

**Menu Services (`app/service/menu`)**:

*   CRUD operations for menus.
*   Function to generate a hierarchical menu tree based on user roles/permissions.
*   Function to provide menu data in a format suitable for Vben Admin's dynamic routing.

**Menu API Implementation (`app/controller/menu`)**:

*   **`GET /menu/routes`**:
    *   **Input**: (authentication token)
    *   **Output**: `array` of menu/route objects, filtered by user permissions.

### Phase 4: Department Management

**Database Schema Design - Department Table (`sys_dept`)**:

*   `id` (PK, BIGSERIAL)
*   `parent_id` (FK to `sys_dept.id` for hierarchical departments)
*   `name` (UNIQUE, VARCHAR) - Department name.
*   `order` (INT) - Order for sorting departments.
*   `status` (TINYINT) - Department status.
*   `created_at` (DATETIME), `updated_at` (DATETIME), `deleted_at` (DATETIME).

**GoFrame Model/DAO Generation for `sys_dept`**.

**Department Services (`app/service/dept`)**:

*   CRUD operations for departments.
*   Function to retrieve hierarchical department data.

**Department API Implementation (`app/controller/dept`)**:

*   Endpoints for managing departments.

### General Tasks (Across all phases)

*   **Error Handling:** Implement consistent error handling and response structures across all API endpoints.
*   **Validation:** Add robust input validation for all API requests.
*   **Logging:** Implement structured logging for debugging and monitoring.
*   **Documentation:** Maintain API documentation (e.g., using Swagger/OpenAPI) for easy frontend integration.
*   **Testing:** Write comprehensive unit and integration tests for all services and controllers.
