ALTER TABLE "sys_permission"
    DROP CONSTRAINT IF EXISTS "uq_sys_permission_tenant_name";

ALTER TABLE "sys_permission"
    DROP COLUMN IF EXISTS "tenant_id";
