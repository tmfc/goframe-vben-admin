-- Seed default tenant and super user

-- Ensure default tenant exists
INSERT INTO sys_tenant (id, name)
VALUES ('00000000-0000-0000-0000-000000000000', 'Default Tenant')
ON CONFLICT (id) DO NOTHING;

-- Insert super user vben with bcrypt hash of "123456"
-- roles stored as JSONB array ["super"], status=1 (enabled)
-- Note: "super" corresponds to consts.RoleSuper
INSERT INTO sys_user (
    id, tenant_id, username, password, real_name, roles, home_path, status
) VALUES (
    '11111111-1111-1111-1111-111111111111',
    '00000000-0000-0000-0000-000000000000',
    'vben',
    '$2a$10$Nhe5dVZhZ3lcVzWH2bvOT.rS95QSk1PCBqzy7FSewXU7GiKiK3mlW',
    'Vben Super Admin',
    '["super"]',
    '/dashboard',
    1
)
ON CONFLICT (tenant_id, username) DO NOTHING;
