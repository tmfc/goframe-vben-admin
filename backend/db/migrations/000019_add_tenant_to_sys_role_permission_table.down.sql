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
        WHERE conname = 'pk_sys_role_permission'
          AND conrelid = 'sys_role_permission'::regclass
    ) THEN
        ALTER TABLE "sys_role_permission" DROP CONSTRAINT pk_sys_role_permission;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'sys_role_permission_pkey'
          AND conrelid = 'sys_role_permission'::regclass
    ) THEN
        ALTER TABLE "sys_role_permission"
            ADD PRIMARY KEY ("role_id", "permission_id");
    END IF;
END $$;
