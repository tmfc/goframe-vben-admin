ALTER TABLE sys_user
    ALTER COLUMN dept_id TYPE UUID USING NULLIF(dept_id::text, '0')::uuid;
