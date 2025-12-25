-- Rollback seeded super user
DELETE FROM sys_user
WHERE tenant_id = '00000000-0000-0000-0000-000000000000'
  AND username = 'vben';

-- Optionally rollback default tenant if it is empty
-- DELETE FROM sys_tenant WHERE id = '00000000-0000-0000-0000-000000000000';
