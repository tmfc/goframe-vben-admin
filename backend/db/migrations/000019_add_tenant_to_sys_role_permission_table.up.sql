DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'sys_role_permission' AND column_name = 'tenant_id'
    ) THEN
        ALTER TABLE "sys_role_permission"
            ADD COLUMN "tenant_id" UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';
    END IF;
END $$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'sys_role_permission_pkey'
          AND conrelid = 'sys_role_permission'::regclass
    ) THEN
        ALTER TABLE "sys_role_permission" DROP CONSTRAINT sys_role_permission_pkey;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'pk_sys_role_permission'
          AND conrelid = 'sys_role_permission'::regclass
    ) THEN
        ALTER TABLE "sys_role_permission"
            ADD CONSTRAINT pk_sys_role_permission PRIMARY KEY ("tenant_id", "role_id", "permission_id");
    END IF;
END $$;
