CREATE TABLE sys_menu (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES sys_tenant(id),
    parent_id BIGINT,
    name VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    component VARCHAR(255),
    icon VARCHAR(255),
    "order" INT NOT NULL DEFAULT 0,
    type VARCHAR(32) NOT NULL,
    visible SMALLINT NOT NULL DEFAULT 1,
    status SMALLINT NOT NULL DEFAULT 1,
    permission_code VARCHAR(255),
    meta JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_sys_menu_tenant ON sys_menu (tenant_id);
CREATE INDEX idx_sys_menu_parent ON sys_menu (parent_id);
CREATE INDEX idx_sys_menu_path ON sys_menu (path);
CREATE INDEX idx_sys_menu_status ON sys_menu (status);
