DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'sys_role_permission' AND column_name = 'tenant_id'
    ) THEN
        ALTER TABLE "sys_role_permission"
            ADD COLUMN "tenant_id" BIGINT NOT NULL DEFAULT 1 REFERENCES sys_tenant(id);
    END IF;
END $$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'uq_sys_role_permission_role_permission'
          AND conrelid = 'sys_role_permission'::regclass
    ) THEN
        ALTER TABLE "sys_role_permission" DROP CONSTRAINT uq_sys_role_permission_role_permission;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'uq_sys_role_permission_tenant_role_permission'
          AND conrelid = 'sys_role_permission'::regclass
    ) THEN
        ALTER TABLE "sys_role_permission"
            ADD CONSTRAINT uq_sys_role_permission_tenant_role_permission UNIQUE ("tenant_id", "role_id", "permission_id");
    END IF;
END $$;
