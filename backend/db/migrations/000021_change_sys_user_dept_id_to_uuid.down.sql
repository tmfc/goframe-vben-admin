ALTER TABLE sys_user
    ALTER COLUMN dept_id TYPE BIGINT USING dept_id::bigint;
