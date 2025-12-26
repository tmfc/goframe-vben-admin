DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'sys_role' AND column_name = 'tenant_id'
    ) THEN
        ALTER TABLE "sys_role"
            ADD COLUMN "tenant_id" UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000' REFERENCES sys_tenant(id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'uq_sys_role_tenant_name'
    ) THEN
        ALTER TABLE "sys_role"
            ADD CONSTRAINT "uq_sys_role_tenant_name" UNIQUE ("tenant_id", "name");
    END IF;
END $$;
