CREATE TABLE sys_dict_data (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES sys_tenant(id),
    dict_type_id BIGINT NOT NULL REFERENCES sys_dict_type(id) ON DELETE CASCADE,
    label VARCHAR(100) NOT NULL,
    label_i18n JSONB,
    value VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    color VARCHAR(20),
    icon VARCHAR(100),
    css_class VARCHAR(100),
    status SMALLINT DEFAULT 1,
    sort INT DEFAULT 0,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    creator_id BIGINT,
    modifier_id BIGINT,
    dept_id BIGINT
);

CREATE INDEX idx_sys_dict_data_tenant ON sys_dict_data (tenant_id);
CREATE INDEX idx_sys_dict_data_type_id ON sys_dict_data (dict_type_id);
CREATE INDEX idx_sys_dict_data_value ON sys_dict_data (dict_type_id, value);
CREATE INDEX idx_sys_dict_data_status ON sys_dict_data (status);
