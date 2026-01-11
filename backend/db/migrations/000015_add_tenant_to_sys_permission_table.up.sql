DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'sys_permission' AND column_name = 'tenant_id'
    ) THEN
        ALTER TABLE "sys_permission"
            ADD COLUMN "tenant_id" BIGINT NOT NULL DEFAULT 1 REFERENCES sys_tenant(id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'uq_sys_permission_tenant_name'
    ) THEN
        ALTER TABLE "sys_permission"
            ADD CONSTRAINT "uq_sys_permission_tenant_name" UNIQUE ("tenant_id", "name");
    END IF;
END $$;
