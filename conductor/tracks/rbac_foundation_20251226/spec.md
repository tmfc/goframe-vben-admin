# Specification: RBAC Data Tables and API Foundation

## Overview
This track focuses on establishing the foundational elements for Role-Based Access Control (RBAC) within the SaaS management platform. This includes designing and creating the necessary data tables (`sys_role`, `sys_permission`, a role-permission pivot table, and `casbin_rule`), generating Data Access Objects (DAOs) for these entities, and initiating the implementation of core role and permission services and APIs. This work will lay the groundwork for future features such as menu and department management, which will rely on these access control mechanisms.

## Functional Requirements

### Data Model Definition

#### 1. `sys_role` Table
-   **Purpose:** Stores information about different user roles.
-   **Fields:**
    -   `id`: Primary Key, unique identifier for the role.
    -   `name`: Unique name of the role (e.g., "Administrator", "Editor").
    -   `description`: Brief description of the role's purpose.
    -   `parent_id`: (Nullable) Foreign key referencing `sys_role.id` to support hierarchical roles.
    -   `status`: Current status of the role (e.g., active, inactive).
    -   `created_at`: Timestamp for record creation.
    -   `updated_at`: Timestamp for last record update.
    -   `creator_id`: ID of the user who created the role.
    -   `modifier_id`: ID of the user who last modified the role.
    -   `dept_id`: Department ID associated with the role.

#### 2. `sys_permission` Table
-   **Purpose:** Stores definitions of application permissions.
-   **Fields:**
    -   `id`: Primary Key, unique identifier for the permission.
    -   `name`: Unique name of the permission (e.g., "user:create", "report:view").
    -   `description`: Brief description of what the permission grants access to.
    -   `parent_id`: (Nullable) Foreign key referencing `sys_permission.id` to support hierarchical permissions.
    -   `status`: Current status of the permission (e.g., active, inactive).
    -   `created_at`: Timestamp for record creation.
    -   `updated_at`: Timestamp for last record update.
    -   `creator_id`: ID of the user who created the permission.
    -   `modifier_id`: ID of the user who last modified the permission.
    -   `dept_id`: Department ID associated with the permission.

#### 3. Role-Permission Pivot Table (e.g., `sys_role_permission`)
-   **Purpose:** Establishes a many-to-many relationship between `sys_role` and `sys_permission`.
-   **Fields:**
    -   `role_id`: Foreign key referencing `sys_role.id`.
    -   `permission_id`: Foreign key referencing `sys_permission.id`.
    -   `created_at`: Timestamp for record creation.
    -   `updated_at`: Timestamp for last record update.
    -   `scope`: Defines the scope of the permission grant (e.g., 'all', 'department', 'self').

#### 4. `casbin_rule` Table
-   **Purpose:** Stores Casbin authorization rules.
-   **Fields:** Adhere to the default Casbin table structure:
    -   `ptype`
    -   `v0`
    -   `v1`
    -   `v2`
    -   `v3`
    -   `v4`
    -   `v5`

#### 5. Additional Metadata Fields for Existing Tables
-   The `sys_user` and `sys_menu` tables must be updated to include the following metadata fields:
    -   `created_at`: Timestamp for record creation.
    -   `updated_at`: Timestamp for last record update.
    -   `creator_id`: ID of the user who created the record.
    -   `modifier_id`: ID of the user who last modified the record.
    -   `dept_id`: Department ID associated with the record.

### Data Access Layer (DAO Generation)
-   Generate DAOs for `sys_role`, `sys_permission`, the role-permission pivot table, and `casbin_rule` using GoFrame's `gf gen dao` command.

### Services and APIs
-   Implement initial services and APIs for basic CRUD operations on `sys_role` and `sys_permission`.
-   Integrate with Casbin for policy enforcement.

## Non-Functional Requirements
-   **Performance:** API responses for role/permission management should be efficient.
-   **Security:** All API endpoints must enforce proper authentication and authorization checks.
-   **Maintainability:** Code should be well-structured, commented, and follow GoFrame conventions.

## Acceptance Criteria
-   All specified data tables are created with the correct schema, including all custom metadata fields.
-   DAOs are successfully generated for the new tables.
-   Basic CRUD APIs for roles and permissions are functional and protected by authentication.
-   Role-permission assignments through the pivot table are correctly managed.
-   Casbin policies can be stored and retrieved via the `casbin_rule` table.
-   Existing `sys_user` and `sys_menu` tables are updated with the new metadata fields.

## Out of Scope
-   Complex role/permission inheritance logic beyond simple `parent_id` relationships.
-   User interface development for RBAC management.
-   Integration of RBAC with other modules (e.g., full menu authorization, data filtering by department) beyond basic API implementation.
