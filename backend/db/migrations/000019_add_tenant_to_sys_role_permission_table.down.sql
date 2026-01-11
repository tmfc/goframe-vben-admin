DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'sys_role_permission' AND column_name = 'tenant_id'
    ) THEN
        ALTER TABLE "sys_role_permission"
            DROP COLUMN "tenant_id";
    END IF;

    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'uq_sys_role_permission_tenant_role_permission'
          AND conrelid = 'sys_role_permission'::regclass
    ) THEN
        ALTER TABLE "sys_role_permission" DROP CONSTRAINT uq_sys_role_permission_tenant_role_permission;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'uq_sys_role_permission_role_permission'
          AND conrelid = 'sys_role_permission'::regclass
    ) THEN
        ALTER TABLE "sys_role_permission"
            ADD CONSTRAINT "uq_sys_role_permission_role_permission" UNIQUE ("role_id", "permission_id");
    END IF;
END $$;
