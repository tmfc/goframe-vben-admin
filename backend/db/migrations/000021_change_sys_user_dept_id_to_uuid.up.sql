ALTER TABLE sys_user
    ALTER COLUMN dept_id TYPE BIGINT USING NULLIF(dept_id::text, '0')::bigint;
