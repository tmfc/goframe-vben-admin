CREATE TABLE sys_user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES sys_tenant(id),
    username VARCHAR(64) NOT NULL,
    password VARCHAR(255) NOT NULL,
    salt VARCHAR(64) NOT NULL,
    real_name VARCHAR(128),
    avatar VARCHAR(255),
    home_path VARCHAR(255),
    status SMALLINT NOT NULL DEFAULT 1,
    roles JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE (tenant_id, username)
);
