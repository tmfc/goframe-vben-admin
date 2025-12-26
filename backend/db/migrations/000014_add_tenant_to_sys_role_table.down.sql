ALTER TABLE "sys_role"
    DROP CONSTRAINT IF EXISTS "uq_sys_role_tenant_name";

ALTER TABLE "sys_role"
    DROP COLUMN IF EXISTS "tenant_id";
