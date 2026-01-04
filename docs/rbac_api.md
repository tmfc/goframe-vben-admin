# RBAC API Notes

## Role Management Extensions

- `POST /sys-role/{id}/assign-users`
  - Body: `{ "userIds": ["..."], "createdBy": 1 }`
- `POST /sys-role/{id}/remove-users`
  - Body: `{ "userIds": ["..."] }`
- `GET /sys-role/{id}/users`
- `POST /sys-role/{id}/assign-permissions`
  - Body: `{ "permissionIds": [1, 2, 3] }`
- `POST /sys-role/{id}/remove-permissions`
  - Body: `{ "permissionIds": [1, 2, 3] }`
- `GET /sys-role/{id}/permissions`

## Permission Lookup

- `GET /sys-permission/by-user/{userId}`

## Notes

- Casbin policies are synchronized per role with domain set to `tenant_id`.
- Role permissions use `sys_permission.name` as the Casbin object match.
