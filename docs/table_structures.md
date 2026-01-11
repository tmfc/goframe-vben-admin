Here are the proposed table structures for `sys_tenant` and `sys_user` in a tabular format, incorporating multi-tenancy:

### `sys_tenant` Table Structure

| Column Name | Data Type              | Constraints/Description                                      |
| :---------- | :--------------------- | :----------------------------------------------------------- |
| `id`        | BIGSERIAL PRIMARY KEY  | Auto-incrementing tenant identifier.                         |
| `name`      | VARCHAR(255) NOT NULL  | Name of the tenant.                                          |
| `status`    | SMALLINT NOT NULL      | Tenant status (e.g., 1 for active, 0 for inactive). Default: 1 |
| `created_at`| TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP | Timestamp for creation.                                      |
| `updated_at`| TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP | Timestamp for last update.                                   |
| `deleted_at`| TIMESTAMP WITH TIME ZONE| Soft deletion timestamp (for logical deletion).              |

### `sys_user` Table Structure

| Column Name   | Data Type              | Constraints/Description                                      |
| :------------ | :--------------------- | :----------------------------------------------------------- |
| `id`          | BIGSERIAL PRIMARY KEY  | Auto-incrementing user identifier.                           |
| `tenant_id`   | BIGINT NOT NULL        | Foreign key to `sys_tenant.id`, identifies the tenant this user belongs to. |
| `username`    | VARCHAR(64) NOT NULL   | User's login username. **Unique per tenant.**                |
| `password`    | VARCHAR(255) NOT NULL  | Hashed password.                                             |
| `salt`        | VARCHAR(64) NOT NULL   | Salt for password hashing.                                   |
| `real_name`   | VARCHAR(128)           | User's real name.                                            |
| `avatar`      | VARCHAR(255)           | URL to user's avatar image.                                  |
| `home_path`   | VARCHAR(255)           | Default dashboard or user-specific home path.                |
| `status`      | SMALLINT NOT NULL      | User account status (e.g., 1 for active, 0 for inactive, 2 for disabled). Default: 1 |
| `roles`       | JSONB                  | User's roles (JSON array of strings, e.g., `["admin", "editor"]`). |
| `created_at`  | TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP | Timestamp for creation.                                      |
| `updated_at`  | TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP | Timestamp for last update.                                   |
| `deleted_at`  | TIMESTAMP WITH TIME ZONE| Soft deletion timestamp (for logical deletion).              |

**Unique Constraint for `sys_user`**:

A unique constraint will be placed on `(tenant_id, username)` to ensure that usernames are unique only within a specific tenant.
Additionally, a foreign key constraint will link `sys_user.tenant_id` to `sys_tenant.id`.
