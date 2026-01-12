CREATE TABLE sys_dict_type (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES sys_tenant(id),
    type_code VARCHAR(50) NOT NULL,
    type_name VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    is_system BOOLEAN DEFAULT FALSE,
    status SMALLINT DEFAULT 1,
    sort INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    creator_id BIGINT,
    modifier_id BIGINT,
    dept_id BIGINT
);

CREATE UNIQUE INDEX idx_sys_dict_type_tenant_code ON sys_dict_type (tenant_id, type_code);
CREATE INDEX idx_sys_dict_type_tenant ON sys_dict_type (tenant_id);
CREATE INDEX idx_sys_dict_type_status ON sys_dict_type (status);
